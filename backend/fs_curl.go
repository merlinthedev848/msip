package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// xmlEscape properly escapes characters for FreeSWITCH XML configuration
func xmlEscape(s string) string {
	var buf bytes.Buffer
	xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

// HandleFSCurl receives mod_xml_curl requests from FreeSWITCH
func HandleFSCurl(c *gin.Context) {
	section := c.PostForm("section")

	switch section {
	case "directory":
		handleDirectory(c)
	case "dialplan":
		handleDialplan(c)
	case "configuration":
		handleConfiguration(c)
	default:
		c.String(http.StatusOK, notFoundXML())
	}
}

func handleDirectory(c *gin.Context) {
	user := strings.TrimSpace(c.PostForm("user"))
	domain := strings.TrimSpace(c.PostForm("domain"))

	if user == "" {
		c.String(http.StatusOK, notFoundXML())
		return
	}

	var ext Extension
	query := DB.Where("extension_number = ? AND is_active = ?", user, true)

	if domain != "" {
		var tenant Tenant
		if err := DB.Where("domain = ?", domain).First(&tenant).Error; err == nil {
			query = query.Where("tenant_id = ?", tenant.ID)
		}
	}

	if err := query.First(&ext).Error; err != nil {
		c.String(http.StatusOK, notFoundXML())
		return
	}

	var tenant Tenant
	if err := DB.First(&tenant, "id = ?", ext.TenantID).Error; err != nil {
		tenant.Domain = domain
		tenant.Name = "MSIP PBX"
	}

	effectiveDomain := domain
	if effectiveDomain == "" {
		effectiveDomain = tenant.Domain
	}

	xmlResponse := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="directory">
    <domain name="%s">
      <params>
        <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"/>
      </params>
      <groups>
        <group name="default">
          <users>
            <user id="%s">
              <params>
                <param name="password" value="%s"/>
                <param name="vm-password" value="%s"/>
              </params>
              <variables>
                <variable name="toll_allow" value="domestic,international,local"/>
                <variable name="accountcode" value="%s"/>
                <variable name="user_context" value="default"/>
                <variable name="effective_caller_id_name" value="Extension %s"/>
                <variable name="effective_caller_id_number" value="%s"/>
                <variable name="outbound_caller_id_name" value="%s"/>
                <variable name="outbound_caller_id_number" value="%s"/>
              </variables>
            </user>
          </users>
        </group>
      </groups>
    </domain>
  </section>
</document>`,
		xmlEscape(effectiveDomain),
		xmlEscape(ext.ExtensionNumber),
		xmlEscape(ext.PasswordHash),
		xmlEscape(ext.ExtensionNumber),
		xmlEscape(ext.TenantID.String()),
		xmlEscape(ext.ExtensionNumber),
		xmlEscape(ext.ExtensionNumber),
		xmlEscape(tenant.Name),
		xmlEscape(ext.ExtensionNumber),
	)

	c.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(xmlResponse))
}

func handleDialplan(c *gin.Context) {
	domain := strings.TrimSpace(c.PostForm("domain"))
	context := strings.TrimSpace(c.PostForm("context"))
	if context == "" {
		context = "default"
	}

	var tenant Tenant
	if domain != "" {
		_ = DB.Where("domain = ?", domain).First(&tenant).Error
	}

	tenantFilter := func(tx *gorm.DB) *gorm.DB {
		if tenant.ID != uuid.Nil {
			return tx.Where("tenant_id = ?", tenant.ID)
		}
		return tx
	}

	var ivrs []IVRMenu
	tenantFilter(DB.Session(&gorm.Session{})).Preload("Choices").Find(&ivrs)

	var inboundRoutes []InboundRoute
	tenantFilter(DB.Session(&gorm.Session{})).Find(&inboundRoutes)

	var outboundRoutes []OutboundRoute
	tenantFilter(DB.Session(&gorm.Session{})).Order("priority asc").Find(&outboundRoutes)

	var trunks []Trunk
	tenantFilter(DB.Session(&gorm.Session{})).Find(&trunks)
	trunkMap := make(map[string]Trunk)
	for _, t := range trunks {
		trunkMap[t.ID.String()] = t
	}

	var videoRooms []VideoRoom
	tenantFilter(DB.Session(&gorm.Session{})).Where("is_active = ?", true).Find(&videoRooms)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="dialplan" description="Dynamic FreeSWITCH / FusionPBX Dialplan">
    <context name="%s">
`, xmlEscape(context)))

	// If context is "public" (inbound trunk carrier calls) or if context is "default":
	// Inbound DID Routes
	for _, ir := range inboundRoutes {
		sb.WriteString(fmt.Sprintf(`
      <!-- Inbound DID Route: %s -->
      <extension name="inbound_%s">
        <condition field="destination_number" expression="^%s$">
          <action application="set" data="domain_name=${domain_name}"/>
`, xmlEscape(ir.DIDNumber), xmlEscape(ir.DIDNumber), xmlEscape(ir.DIDNumber)))

		switch strings.ToUpper(ir.DestinationType) {
		case "IVR":
			sb.WriteString(fmt.Sprintf(`          <action application="answer"/>
          <action application="sleep" data="1000"/>
          <action application="ivr" data="%s"/>
`, xmlEscape(ir.DestinationTarget)))
		case "VOICEMAIL":
			sb.WriteString(fmt.Sprintf(`          <action application="answer"/>
          <action application="voicemail" data="default ${domain_name} %s"/>
`, xmlEscape(ir.DestinationTarget)))
		default: // EXTENSION
			sb.WriteString(fmt.Sprintf(`          <action application="bridge" data="user/%s@${domain_name}"/>
`, xmlEscape(ir.DestinationTarget)))
		}

		sb.WriteString(`        </condition>
      </extension>
`)
	}

	// Internal extensions & feature codes (executed for default / internal contexts)
	if context != "public" {
		sb.WriteString(`
      <!-- Echo Test Route (FusionPBX 9999 standard) -->
      <extension name="echo_test">
        <condition field="destination_number" expression="^9999$">
          <action application="answer"/>
          <action application="echo"/>
        </condition>
      </extension>

      <!-- Check Personal Voicemail (*97 standard) -->
      <extension name="check_personal_voicemail">
        <condition field="destination_number" expression="^\*97$">
          <action application="answer"/>
          <action application="sleep" data="1000"/>
          <action application="voicemail" data="check default ${domain_name} ${caller_id_number}"/>
        </condition>
      </extension>

      <!-- Check Any Voicemail Box (*98 standard) -->
      <extension name="check_general_voicemail">
        <condition field="destination_number" expression="^\*98$">
          <action application="answer"/>
          <action application="sleep" data="1000"/>
          <action application="voicemail" data="check default ${domain_name}"/>
        </condition>
      </extension>

      <!-- AI Dictation / Transcription Service (*88) -->
      <extension name="Dictation_AI">
        <condition field="destination_number" expression="^\*88$">
          <action application="answer"/>
          <action application="sleep" data="1000"/>
          <action application="playback" data="tone_stream://%(1000,0,800)"/>
          <action application="record" data="/usr/local/freeswitch/recordings/${uuid}_dictation.wav 600 200"/>
          <action application="hangup"/>
        </condition>
      </extension>

      <!-- Conversational AI Receptionist (*99 WebSocket stream) -->
      <extension name="Conversational_AI_Receptionist">
        <condition field="destination_number" expression="^\*99$">
          <action application="answer"/>
          <action application="uuid_audio_fork" data="${uuid} ws://pbx_backend:8081/api/ai/stream mono 8k"/>
          <action application="park"/>
        </condition>
      </extension>

      <!-- Default Video MCU Conference (9000) -->
      <extension name="Video_Room_9000">
        <condition field="destination_number" expression="^9000$">
          <action application="answer"/>
          <action application="conference" data="video_mcu_9000@video_layout_1++flags{video_muxing}"/>
        </condition>
      </extension>
`)

		// Dynamic Video Rooms
		for _, vr := range videoRooms {
			sb.WriteString(fmt.Sprintf(`
      <!-- Dynamic Video Room: %s -->
      <extension name="video_room_%s">
        <condition field="destination_number" expression="^%s$">
          <action application="answer"/>
          <action application="conference" data="room_%s@video_layout_1++flags{video_muxing}"/>
        </condition>
      </extension>
`, xmlEscape(vr.Name), xmlEscape(vr.RoomCode), xmlEscape(vr.RoomCode), xmlEscape(vr.RoomCode)))
		}

		// Dynamic IVR Menus & Ring Groups
		for _, ivr := range ivrs {
			if ivr.Greeting == "ring_group" || strings.HasPrefix(ivr.Name, "Ring Group:") {
				var targets []string
				for _, choice := range ivr.Choices {
					if choice.Destination != "" {
						targets = append(targets, fmt.Sprintf("user/%s@${domain_name}", choice.Destination))
					}
				}
				if len(targets) > 0 {
					bridgeTargets := strings.Join(targets, ",")
					sb.WriteString(fmt.Sprintf(`
      <!-- Ring Group Route: %s -->
      <extension name="ring_group_%s">
        <condition field="destination_number" expression="^%s$">
          <action application="set" data="hangup_after_bridge=true"/>
          <action application="set" data="call_timeout=30"/>
          <action application="bridge" data="%s"/>
        </condition>
      </extension>
`, xmlEscape(ivr.Name), xmlEscape(ivr.Extension), xmlEscape(ivr.Extension), xmlEscape(bridgeTargets)))
					continue
				}
			}

			sb.WriteString(fmt.Sprintf(`
      <!-- IVR Route: %s -->
      <extension name="ivr_%s">
        <condition field="destination_number" expression="^%s$">
          <action application="answer"/>
          <action application="sleep" data="1000"/>
          <action application="ivr" data="%s"/>
        </condition>
      </extension>
`, xmlEscape(ivr.Name), xmlEscape(ivr.Name), xmlEscape(ivr.Extension), xmlEscape(ivr.Name)))
		}

		// Outbound Routes (Dialpatterns prioritized asc)
		for _, obr := range outboundRoutes {
			trunk, ok := trunkMap[obr.TrunkID.String()]
			if ok && obr.DialPattern != "" {
				bridgeData := fmt.Sprintf("sofia/gateway/%s/$1", trunk.ProviderName)
				if trunk.AuthMethod == "IP_ACL" {
					bridgeData = fmt.Sprintf("sofia/external/$1@%s", trunk.SIPServer)
				}
				sb.WriteString(fmt.Sprintf(`
      <!-- Outbound Route: %s via %s -->
      <extension name="outbound_%s">
        <condition field="destination_number" expression="%s">
          <action application="set" data="effective_caller_id_number=${outbound_caller_id_number}"/>
          <action application="bridge" data="%s"/>
        </condition>
      </extension>
`, xmlEscape(obr.DialPattern), xmlEscape(trunk.ProviderName), xmlEscape(obr.ID.String()[:8]), xmlEscape(obr.DialPattern), xmlEscape(bridgeData)))
			}
		}

		// Local Extension Dialing (2 to 5 digits) with Voicemail Fallback
		sb.WriteString(`
      <!-- Local Extension Routing with Voicemail Fallback -->
      <extension name="local_extension">
        <condition field="destination_number" expression="^(\d{2,5})$">
          <action application="set" data="dialed_extension=$1"/>
          <action application="set" data="call_timeout=30"/>
          <action application="set" data="hangup_after_bridge=true"/>
          <action application="bridge" data="user/${dialed_extension}@${domain_name}"/>
          <action application="answer"/>
          <action application="sleep" data="1000"/>
          <action application="voicemail" data="default ${domain_name} ${dialed_extension}"/>
        </condition>
      </extension>
`)
	}

	sb.WriteString(`
    </context>
  </section>
</document>`)

	c.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(sb.String()))
}

func handleConfiguration(c *gin.Context) {
	keyName := c.PostForm("key_value")

	switch keyName {
	case "ivr.conf":
		handleIVRConfig(c)
	case "acl.conf":
		handleACLConfig(c)
	case "callcenter.conf":
		handleCallCenterConfig(c)
	case "conference.conf":
		handleConferenceConfig(c)
	case "voicemail.conf":
		handleVoicemailConfig(c)
	case "sofia.conf":
		handleSofiaConfig(c)
	default:
		c.String(http.StatusOK, notFoundXML())
	}
}

func handleIVRConfig(c *gin.Context) {
	var ivrs []IVRMenu
	DB.Preload("Choices").Find(&ivrs)

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="configuration" description="IVR Menus">
    <configuration name="ivr.conf" description="IVR menus">
      <menus>
`)

	for _, ivr := range ivrs {
		timeout := ivr.Timeout
		if timeout <= 0 {
			timeout = 10000
		}
		maxFailures := ivr.MaxFailures
		if maxFailures <= 0 {
			maxFailures = 3
		}

		sb.WriteString(fmt.Sprintf(`        <menu name="%s"
            greet-long="%s"
            greet-short="%s"
            invalid-sound="ivr/ivr-that_was_an_invalid_entry.wav"
            exit-sound="voicemail/vm-goodbye.wav"
            timeout="%d"
            max-failures="%d">
`,
			xmlEscape(ivr.Name),
			xmlEscape(ivr.Greeting),
			xmlEscape(ivr.Greeting),
			timeout,
			maxFailures,
		))

		for _, choice := range ivr.Choices {
			sb.WriteString(fmt.Sprintf(`          <entry action="%s" digits="%s" param="%s"/>
`,
				xmlEscape(choice.Action),
				xmlEscape(choice.Digits),
				xmlEscape(choice.Destination),
			))
		}

		sb.WriteString(`        </menu>
`)
	}

	sb.WriteString(`      </menus>
    </configuration>
  </section>
</document>`)

	c.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(sb.String()))
}

func handleACLConfig(c *gin.Context) {
	var firewallRules []FirewallRule
	DB.Find(&firewallRules)

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="configuration" description="Network Lists">
    <configuration name="acl.conf" description="Network Lists">
      <network-lists>
        <list name="domains" default="deny">
          <node type="allow" domain="$${domain}"/>
          <node type="allow" cidr="127.0.0.1/32"/>
          <node type="allow" cidr="10.0.0.0/8"/>
          <node type="allow" cidr="172.16.0.0/12"/>
          <node type="allow" cidr="192.168.0.0/16"/>
`)

	for _, rule := range firewallRules {
		action := "allow"
		if strings.ToLower(rule.Action) == "block" || strings.ToLower(rule.Action) == "drop" {
			action = "deny"
		}
		ip := strings.TrimSpace(rule.IPAddress)
		if !strings.Contains(ip, "/") {
			ip = ip + "/32"
		}
		sb.WriteString(fmt.Sprintf(`          <node type="%s" cidr="%s"/>
`, action, xmlEscape(ip)))
	}

	sb.WriteString(`        </list>
        <list name="lan" default="allow">
          <node type="allow" cidr="127.0.0.1/32"/>
          <node type="allow" cidr="10.0.0.0/8"/>
          <node type="allow" cidr="172.16.0.0/12"/>
          <node type="allow" cidr="192.168.0.0/16"/>
        </list>
      </network-lists>
    </configuration>
  </section>
</document>`)

	c.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(sb.String()))
}

func handleSofiaConfig(c *gin.Context) {
	var trunks []Trunk
	DB.Find(&trunks)

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="configuration" description="Sofia Gateways">
    <configuration name="sofia.conf" description="Sofia Endpoints">
      <profiles>
        <profile name="external">
          <gateways>
`)

	for _, trunk := range trunks {
		reg := "false"
		if trunk.Register {
			reg = "true"
		}
		sb.WriteString(fmt.Sprintf(`            <gateway name="%s">
              <param name="realm" value="%s"/>
              <param name="proxy" value="%s"/>
              <param name="username" value="%s"/>
              <param name="password" value="%s"/>
              <param name="register" value="%s"/>
              <param name="caller-id-in-from" value="false"/>
            </gateway>
`,
			xmlEscape(trunk.ProviderName),
			xmlEscape(trunk.SIPServer),
			xmlEscape(trunk.SIPServer),
			xmlEscape(trunk.Username),
			xmlEscape(trunk.Password),
			reg,
		))
	}

	sb.WriteString(`          </gateways>
        </profile>
      </profiles>
    </configuration>
  </section>
</document>`)

	c.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(sb.String()))
}

func handleCallCenterConfig(c *gin.Context) {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="configuration" description="Call Center Queues">
    <configuration name="callcenter.conf" description="CallCenter">
      <settings>
        <param name="odbc-dsn" value=""/>
      </settings>
      <queues>
        <queue name="support_tier_1@default">
          <param name="strategy" value="longest-idle-agent"/>
          <param name="moh-sound" value="local_stream://default"/>
          <param name="time-base-score" value="system"/>
          <param name="max-wait-time" value="300"/>
          <param name="max-wait-time-with-no-agent" value="30"/>
          <param name="tier-rules-apply" value="false"/>
        </queue>
        <queue name="sales_global@default">
          <param name="strategy" value="round-robin"/>
          <param name="moh-sound" value="local_stream://default"/>
          <param name="time-base-score" value="system"/>
          <param name="max-wait-time" value="300"/>
        </queue>
        <queue name="billing@default">
          <param name="strategy" value="top-down"/>
          <param name="moh-sound" value="local_stream://default"/>
          <param name="time-base-score" value="system"/>
          <param name="max-wait-time" value="180"/>
        </queue>
      </queues>
    </configuration>
  </section>
</document>`)

	c.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(sb.String()))
}

func handleConferenceConfig(c *gin.Context) {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="configuration" description="Conference Profiles">
    <configuration name="conference.conf" description="Audio/Video Conference">
      <caller-controls>
        <group name="default">
          <control action="mute" digits="0"/>
          <control action="deaf mute" digits="*"/>
          <control action="energy up" digits="9"/>
          <control action="energy equ" digits="8"/>
          <control action="energy dn" digits="7"/>
          <control action="vol talk up" digits="3"/>
          <control action="vol talk dn" digits="1"/>
          <control action="hangup" digits="#"/>
        </group>
      </caller-controls>
      <profiles>
        <profile name="video_layout_1">
          <param name="rate" value="16000"/>
          <param name="interval" value="20"/>
          <param name="energy-level" value="300"/>
          <param name="caller-controls" value="default"/>
          <param name="video-mode" value="mux"/>
          <param name="video-layout-name" value="1x1"/>
          <param name="video-canvas-size" value="1280x720"/>
        </profile>
      </profiles>
    </configuration>
  </section>
</document>`)

	c.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(sb.String()))
}

func handleVoicemailConfig(c *gin.Context) {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="configuration" description="Voicemail Configuration">
    <configuration name="voicemail.conf" description="Voicemail">
      <settings>
        <param name="password-complexity" value="false"/>
      </settings>
      <profiles>
        <profile name="default">
          <param name="file-ext" value="wav"/>
          <param name="digit-timeout" value="10000"/>
          <param name="max-login-attempts" value="3"/>
          <param name="min-record-len" value="3"/>
          <param name="max-record-len" value="300"/>
          <param name="tone-spec" value="%(1000, 0, 640)"/>
          <param name="record-silence-hits" value="2"/>
        </profile>
      </profiles>
    </configuration>
  </section>
</document>`)

	c.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(sb.String()))
}

func notFoundXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="result">
    <result status="not found" />
  </section>
</document>`
}
