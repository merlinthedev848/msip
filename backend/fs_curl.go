package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleFSCurl receives mod_xml_curl requests from FreeSWITCH
func HandleFSCurl(c *gin.Context) {
	// FreeSWITCH sends data as standard POST form values
	section := c.PostForm("section")
	
	switch section {
	case "directory":
		handleDirectory(c)
	case "dialplan":
		handleDialplan(c)
	case "configuration":
		handleConfiguration(c)
	default:
		// Return empty NOT FOUND document for unhandled sections
		c.String(http.StatusOK, notFoundXML())
	}
}

func handleDirectory(c *gin.Context) {
	user := c.PostForm("user")
	domain := c.PostForm("domain")

	if user == "" {
		c.String(http.StatusOK, notFoundXML())
		return
	}

	var ext Extension
	result := DB.Where("extension_number = ?", user).First(&ext)

	if result.Error != nil || !ext.IsActive {
		c.String(http.StatusOK, notFoundXML())
		return
	}

	xmlTemplate := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
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
              </params>
              <variables>
                <variable name="toll_allow" value="domestic,international,local"/>
                <variable name="accountcode" value="%s"/>
                <variable name="user_context" value="default"/>
                <variable name="effective_caller_id_name" value="Extension %s"/>
                <variable name="effective_caller_id_number" value="%s"/>
                <variable name="outbound_caller_id_name" value="MSIP PBX"/>
                <variable name="outbound_caller_id_number" value="0000000000"/>
              </variables>
            </user>
          </users>
        </group>
      </groups>
    </domain>
  </section>
</document>`

	xmlResponse := fmt.Sprintf(xmlTemplate, domain, ext.ExtensionNumber, ext.PasswordHash, ext.TenantID.String(), ext.ExtensionNumber, ext.ExtensionNumber)
	c.Data(http.StatusOK, "text/xml", []byte(xmlResponse))
}

func handleDialplan(c *gin.Context) {
	xml := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="dialplan" description="Dynamic Go Backend Dialplan">
    <context name="default">
      
      <!-- Echo Test Route -->
      <extension name="echo_test">
        <condition field="destination_number" expression="^9999$">
          <action application="answer"/>
          <action application="echo"/>
        </condition>
      </extension>`

	// Inject IVR Routes
	var ivrs []IVRMenu
	DB.Find(&ivrs)
	for _, ivr := range ivrs {
		xml += fmt.Sprintf(`
      <!-- IVR Route: %s -->
      <extension name="ivr_%s">
        <condition field="destination_number" expression="^%s$">
          <action application="answer"/>
          <action application="sleep" data="1000"/>
          <action application="ivr" data="%s"/>
        </condition>
      </extension>`, ivr.Name, ivr.Name, ivr.Extension, ivr.Name)
	}

	xml += `
      <!-- AI Dictation / Transcription Service -->
      <extension name="ai_dictation">
        <condition field="destination_number" expression="^\*88$">
          <action application="answer"/>
          <!-- Play a 1-second beep (800Hz) to signal start of recording -->
          <action application="playback" data="tone_stream://%(1000,0,800)"/>
          <action application="set" data="playback_terminators=#"/>
          <!-- Record the audio (max 60 seconds) to the shared volume -->
          <action application="record" data="/usr/local/freeswitch/recordings/${uuid}_dictation.wav 60"/>
          <!-- Play a double beep to signal completion -->
          <action application="playback" data="tone_stream://%(200,50,800);%(200,0,800)"/>
          <action application="hangup"/>
        </condition>
      </extension>

      <!-- Local Extension Routing with Voicemail -->
      <extension name="local_extension">
        <condition field="destination_number" expression="^(\d{2,5})$">
          <!-- Bridge to registered SIP user with 30s timeout -->
          <action application="set" data="dialed_extension=$1"/>
          <action application="set" data="call_timeout=30"/>
          <action application="set" data="hangup_after_bridge=true"/>
          <action application="bridge" data="user/${dialed_extension}@${domain_name}"/>
          
          <!-- If unanswered/offline, fallback to voicemail -->
          <action application="answer"/>
          <action application="sleep" data="1000"/>
          <action application="voicemail" data="default ${domain_name} ${dialed_extension}"/>
        </condition>
      </extension>

    </context>
  </section>
</document>`

	c.Data(http.StatusOK, "text/xml", []byte(xml))
}

func handleConfiguration(c *gin.Context) {
	keyName := c.PostForm("key_value") // FreeSWITCH passes key_value=ivr.conf
	
	if keyName == "ivr.conf" {
		var ivrs []IVRMenu
		DB.Preload("Choices").Find(&ivrs)
		
		xml := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="configuration" description="IVR Menus">
    <configuration name="ivr.conf" description="IVR menus">
      <menus>`
      
		for _, ivr := range ivrs {
			xml += fmt.Sprintf(`
        <menu name="%s"
            greet-long="%s"
            greet-short="%s"
            invalid-sound="ivr/ivr-that_was_an_invalid_entry.wav"
            exit-sound="voicemail/vm-goodbye.wav"
            timeout="%d"
            max-failures="%d">`, ivr.Name, ivr.Greeting, ivr.Greeting, ivr.Timeout, ivr.MaxFailures)
            
			for _, choice := range ivr.Choices {
				xml += fmt.Sprintf(`
          <entry action="%s" digits="%s" param="%s"/>`, choice.Action, choice.Digits, choice.Destination)
			}
			xml += `
        </menu>`
		}
		
		xml += `
      </menus>
    </configuration>
  </section>
</document>`

		c.Data(http.StatusOK, "text/xml", []byte(xml))
		return
	}
	
	c.String(http.StatusOK, notFoundXML())
}

func notFoundXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<document type="freeswitch/xml">
  <section name="result">
    <result status="not found" />
  </section>
</document>`
}
