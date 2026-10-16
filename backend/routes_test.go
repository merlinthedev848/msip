package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var testRouter *gin.Engine
var adminToken string
var viewerToken string
var tenantA uuid.UUID
var tenantB uuid.UUID

func setupTestEnvironment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_DRIVER", "sqlite")
	os.Setenv("DATABASE_URL", "file:test_msip.db?cache=shared&mode=memory")
	os.Setenv("JWT_SECRET", "test-secret-key")
	jwtSecret = []byte("test-secret-key")

	InitDB()

	// Clear out tables for pristine test state
	DB.Exec("DELETE FROM extensions")
	DB.Exec("DELETE FROM trunks")
	DB.Exec("DELETE FROM inbound_routes")
	DB.Exec("DELETE FROM outbound_routes")
	DB.Exec("DELETE FROM ivr_choices")
	DB.Exec("DELETE FROM ivr_menus")
	DB.Exec("DELETE FROM cdrs")
	DB.Exec("DELETE FROM voicemails")
	DB.Exec("DELETE FROM video_rooms")
	DB.Exec("DELETE FROM ip_cameras")
	DB.Exec("DELETE FROM team_chats")
	DB.Exec("DELETE FROM sms_messages")
	DB.Exec("DELETE FROM api_credentials")
	DB.Exec("DELETE FROM webhooks")
	DB.Exec("DELETE FROM throttle_rules")
	DB.Exec("DELETE FROM firewall_rules")
	DB.Exec("DELETE FROM fraud_rules")
	DB.Exec("DELETE FROM users")
	DB.Exec("DELETE FROM tenants")
	DB.Exec("DELETE FROM ai_assistants")
	DB.Exec("DELETE FROM ai_bookings")
	DB.Exec("DELETE FROM ai_call_sessions")
	DB.Exec("DELETE FROM ai_transcription_settings")
	DB.Exec("DELETE FROM call_queues")
	DB.Exec("DELETE FROM queue_members")
	DB.Exec("DELETE FROM call_distribution_groups")
	DB.Exec("DELETE FROM group_members")
	DB.Exec("DELETE FROM call_diverts")
	DB.Exec("DELETE FROM devices")
	DB.Exec("DELETE FROM telephone_numbers")
	DB.Exec("DELETE FROM sites")
	DB.Exec("DELETE FROM contacts")
	DB.Exec("DELETE FROM button_profiles")
	DB.Exec("DELETE FROM recording_rules")
	DB.Exec("DELETE FROM audio_prompts")
	DB.Exec("DELETE FROM feature_codes")
	DB.Exec("DELETE FROM paging_groups")
	DB.Exec("DELETE FROM pickup_groups")
	DB.Exec("DELETE FROM auto_clips")
	DB.Exec("DELETE FROM outbound_restrictions")
	DB.Exec("DELETE FROM holidays")
	DB.Exec("DELETE FROM time_conditions")

	// Seed Tenant A
	tA := Tenant{ID: uuid.New(), Name: "Tenant Alpha", Domain: "alpha.pbx.local"}
	DB.Create(&tA)
	tenantA = tA.ID

	// Seed Tenant B
	tB := Tenant{ID: uuid.New(), Name: "Tenant Beta", Domain: "beta.pbx.local"}
	DB.Create(&tB)
	tenantB = tB.ID

	// Seed Superadmin user in Tenant A
	hashAdmin, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	uAdmin := User{
		ID:           uuid.New(),
		TenantID:     tenantA,
		Email:        "admin@alpha.local",
		PasswordHash: string(hashAdmin),
		Role:         "superadmin",
	}
	DB.Create(&uAdmin)

	// Seed Viewer user in Tenant A
	hashViewer, _ := bcrypt.GenerateFromPassword([]byte("viewer123"), bcrypt.DefaultCost)
	uViewer := User{
		ID:           uuid.New(),
		TenantID:     tenantA,
		Email:        "viewer@alpha.local",
		PasswordHash: string(hashViewer),
		Role:         "viewer",
	}
	DB.Create(&uViewer)

	// Seed an Extension in Tenant A
	extA := Extension{
		Base:             Base{ID: uuid.New()},
		TenantRef:        TenantRef{TenantID: tenantA},
		ExtensionNumber:  "1001",
		PasswordHash:     "secretpass",
		IsActive:         true,
		ProvisioningCode: "10011001",
	}
	DB.Create(&extA)

	// Seed an Extension in Tenant B
	extB := Extension{
		Base:             Base{ID: uuid.New()},
		TenantRef:        TenantRef{TenantID: tenantB},
		ExtensionNumber:  "2001",
		PasswordHash:     "secretpass2",
		IsActive:         true,
		ProvisioningCode: "20012001",
	}
	DB.Create(&extB)

	// Generate JWT for admin
	tokA := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       uAdmin.ID.String(),
		"tenant_id": tenantA.String(),
		"role":      "superadmin",
		"exp":       time.Now().Add(time.Hour).Unix(),
	})
	adminToken, _ = tokA.SignedString(jwtSecret)

	// Generate JWT for viewer
	tokV := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       uViewer.ID.String(),
		"tenant_id": tenantA.String(),
		"role":      "viewer",
		"exp":       time.Now().Add(time.Hour).Unix(),
	})
	viewerToken, _ = tokV.SignedString(jwtSecret)

	testRouter = SetupRouter()
}

func performRequest(r *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf *bytes.Buffer
	if body != nil {
		switch b := body.(type) {
		case []byte:
			buf = bytes.NewBuffer(b)
		case string:
			buf = bytes.NewBufferString(b)
		default:
			data, _ := json.Marshal(b)
			buf = bytes.NewBuffer(data)
		}
	} else {
		buf = bytes.NewBuffer(nil)
	}
	req, _ := http.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performFormRequest(r *gin.Engine, method, path string, form url.Values) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHealthEndpoint(t *testing.T) {
	setupTestEnvironment(t)
	w := performRequest(testRouter, "GET", "/health", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "healthy") {
		t.Fatalf("expected healthy status, got %s", w.Body.String())
	}
}

func TestAuthLogin(t *testing.T) {
	setupTestEnvironment(t)

	// Valid login via /api/auth/login
	loginBody := `{"email": "admin@alpha.local", "password": "admin123"}`
	w := performRequest(testRouter, "POST", "/api/auth/login", []byte(loginBody), "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["token"] == nil || res["token"] == "" {
		t.Fatalf("expected token in response")
	}

	// Valid login via /api/v1/auth/login (direct frontend / nginx path)
	wV1 := performRequest(testRouter, "POST", "/api/v1/auth/login", []byte(loginBody), "")
	if wV1.Code != http.StatusOK {
		t.Fatalf("expected 200 for /api/v1/auth/login, got %d: %s", wV1.Code, wV1.Body.String())
	}

	// Valid login via /auth/login
	wShort := performRequest(testRouter, "POST", "/auth/login", []byte(loginBody), "")
	if wShort.Code != http.StatusOK {
		t.Fatalf("expected 200 for /auth/login, got %d: %s", wShort.Code, wShort.Body.String())
	}

	// Invalid password
	badBody := `{"email": "admin@alpha.local", "password": "wrong"}`
	wBad := performRequest(testRouter, "POST", "/api/auth/login", []byte(badBody), "")
	if wBad.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong password, got %d", wBad.Code)
	}

	// Non-existent email
	noUser := `{"email": "nobody@nowhere.com", "password": "password"}`
	wNo := performRequest(testRouter, "POST", "/api/auth/login", []byte(noUser), "")
	if wNo.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for non-existent user, got %d", wNo.Code)
	}
}

func TestAutoProvisioning(t *testing.T) {
	setupTestEnvironment(t)

	w := performRequest(testRouter, "GET", "/api/provisioning/10011001", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid code, got %d: %s", w.Code, w.Body.String())
	}
	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["sip_username"] != "1001" || res["sip_domain"] != "alpha.pbx.local" {
		t.Fatalf("unexpected provisioning payload: %+v", res)
	}

	wBad := performRequest(testRouter, "GET", "/api/provisioning/unknown", nil, "")
	if wBad.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown code, got %d", wBad.Code)
	}
}

func TestFreeSWITCHCurlDirectory(t *testing.T) {
	setupTestEnvironment(t)

	form := url.Values{
		"section": {"directory"},
		"user":    {"1001"},
		"domain":  {"alpha.pbx.local"},
	}
	w := performFormRequest(testRouter, "POST", "/fs", form)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	xml := w.Body.String()
	if !strings.Contains(xml, `user id="1001"`) || !strings.Contains(xml, `alpha.pbx.local`) {
		t.Fatalf("expected XML to contain user and domain, got: %s", xml)
	}

	// Non-existent user
	formMissing := url.Values{
		"section": {"directory"},
		"user":    {"9999"},
		"domain":  {"alpha.pbx.local"},
	}
	wMissing := performFormRequest(testRouter, "POST", "/fs", formMissing)
	if !strings.Contains(wMissing.Body.String(), `<result status="not found" />`) {
		t.Fatalf("expected not found XML, got: %s", wMissing.Body.String())
	}
}

func TestFreeSWITCHCurlDialplan(t *testing.T) {
	setupTestEnvironment(t)

	// Create IVR and Trunk to test dialplan injection
	ivr := IVRMenu{
		Base:      Base{ID: uuid.New()},
		TenantRef: TenantRef{TenantID: tenantA},
		Name:      "Main_Menu",
		Extension: "7000",
		Greeting:  "welcome.wav",
	}
	DB.Create(&ivr)

	trunk := Trunk{
		Base:         Base{ID: uuid.New()},
		TenantRef:    TenantRef{TenantID: tenantA},
		ProviderName: "Twilio_Primary",
		SIPServer:    "sip.twilio.com",
		AuthMethod:   "IP_ACL",
	}
	DB.Create(&trunk)

	obr := OutboundRoute{
		Base:        Base{ID: uuid.New()},
		TenantRef:   TenantRef{TenantID: tenantA},
		DialPattern: "^0[1-9]\\d+$",
		TrunkID:     trunk.ID,
		Priority:    1,
	}
	DB.Create(&obr)

	form := url.Values{
		"section": {"dialplan"},
		"domain":  {"alpha.pbx.local"},
	}
	w := performFormRequest(testRouter, "POST", "/fs", form)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	xml := w.Body.String()
	if !strings.Contains(xml, "echo_test") || !strings.Contains(xml, "ivr_Main_Menu") || !strings.Contains(xml, "outbound_") {
		t.Fatalf("expected dialplan with echo, IVR, and outbound route, got: %s", xml)
	}

	// Verify FusionPBX feature codes *97 and *98 are in dialplan
	if !strings.Contains(xml, "*97") || !strings.Contains(xml, "*98") {
		t.Fatalf("expected dialplan to include *97 and *98 voicemail feature codes, got: %s", xml)
	}
}

func TestFreeSWITCHConfigurations(t *testing.T) {
	setupTestEnvironment(t)

	// Add a firewall rule and trunk for config tests
	fw := FirewallRule{
		Base:      Base{ID: uuid.New()},
		TenantRef: TenantRef{TenantID: tenantA},
		IPAddress: "203.0.113.5",
		Action:    "allow",
	}
	DB.Create(&fw)

	trunk := Trunk{
		Base:         Base{ID: uuid.New()},
		TenantRef:    TenantRef{TenantID: tenantA},
		ProviderName: "Telnyx_Gateway",
		SIPServer:    "sip.telnyx.com",
		AuthMethod:   "USERPASS",
		Username:     "user123",
		Password:     "pass123",
		Register:     true,
	}
	DB.Create(&trunk)

	configs := []struct {
		name     string
		expected string
	}{
		{"ivr.conf", "ivr.conf"},
		{"acl.conf", "acl.conf"},
		{"acl.conf", "203.0.113.5/32"},
		{"sofia.conf", "Telnyx_Gateway"},
		{"callcenter.conf", "support_tier_1@default"},
		{"conference.conf", "video_layout_1"},
		{"voicemail.conf", "voicemail.conf"},
	}

	for _, cfg := range configs {
		form := url.Values{
			"section":   {"configuration"},
			"key_value": {cfg.name},
		}
		w := performFormRequest(testRouter, "POST", "/fs", form)
		if w.Code != http.StatusOK {
			t.Fatalf("config %s expected 200, got %d", cfg.name, w.Code)
		}
		if !strings.Contains(w.Body.String(), cfg.expected) {
			t.Fatalf("config %s expected to contain %q, got: %s", cfg.name, cfg.expected, w.Body.String())
		}
	}
}

func TestFreeSWITCHCDRIngestion(t *testing.T) {
	setupTestEnvironment(t)

	cdrPayload := `{
		"variables": {
			"uuid": "test-call-uuid-1234",
			"caller_id_name": "John Doe",
			"caller_id_number": "1001",
			"destination_number": "1002",
			"context": "default",
			"duration": 45,
			"billsec": 40,
			"hangup_cause": "NORMAL_CLEARING",
			"domain_name": "alpha.pbx.local"
		}
	}`

	w := performRequest(testRouter, "POST", "/api/cdr", []byte(cdrPayload), "")
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for CDR, got %d: %s", w.Code, w.Body.String())
	}

	var saved CDR
	if err := DB.Where("uuid = ?", "test-call-uuid-1234").First(&saved).Error; err != nil {
		t.Fatalf("expected CDR to be stored in DB: %v", err)
	}
	if saved.Duration != 45 || saved.Billsec != 40 {
		t.Fatalf("CDR duration mismatch: %+v", saved)
	}

	// Verify GET /api/cdr returns the CDR
	wGet := performRequest(testRouter, "GET", "/api/cdr", nil, adminToken)
	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 GET /api/cdr, got %d", wGet.Code)
	}
	var getRes map[string][]CDR
	json.Unmarshal(wGet.Body.Bytes(), &getRes)
	if len(getRes["cdrs"]) == 0 {
		t.Fatalf("expected at least 1 CDR in list")
	}

	// Verify PUT /api/cdr/:id updates the record
	updateBody := `{"Transcription": "Updated transcription text"}`
	wPut := performRequest(testRouter, "PUT", "/api/cdr/"+saved.UUID, []byte(updateBody), adminToken)
	if wPut.Code != http.StatusOK {
		t.Fatalf("expected 200 PUT /api/cdr/:id, got %d", wPut.Code)
	}

	// Verify DELETE /api/cdr/:id removes the record
	wDel := performRequest(testRouter, "DELETE", "/api/cdr/"+saved.UUID, nil, adminToken)
	if wDel.Code != http.StatusOK {
		t.Fatalf("expected 200 DELETE /api/cdr/:id, got %d", wDel.Code)
	}
}

func TestExtensionsCRUD(t *testing.T) {
	setupTestEnvironment(t)

	// List extensions
	wList := performRequest(testRouter, "GET", "/api/extensions", nil, adminToken)
	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", wList.Code)
	}

	// Create extension
	newExt := `{"ExtensionNumber": "1002", "PasswordHash": "pass1002", "IsActive": true}`
	wCreate := performRequest(testRouter, "POST", "/api/extensions", []byte(newExt), adminToken)
	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", wCreate.Code, wCreate.Body.String())
	}
	var created Extension
	var createResp struct {
		Extension Extension `json:"extension"`
	}
	_ = json.Unmarshal(wCreate.Body.Bytes(), &createResp)
	if createResp.Extension.ID != uuid.Nil {
		created = createResp.Extension
	} else {
		_ = json.Unmarshal(wCreate.Body.Bytes(), &created)
	}
	if created.PasswordHash != "" {
		t.Fatalf("expected password hash to be sanitized from response")
	}

	// Duplicate extension number should fail with 409 Conflict
	wDup := performRequest(testRouter, "POST", "/api/extensions", []byte(newExt), adminToken)
	if wDup.Code != http.StatusConflict {
		t.Fatalf("expected 409 conflict for duplicate extension, got %d", wDup.Code)
	}

	// Update extension
	updateExt := fmt.Sprintf(`{"ExtensionNumber": "1002", "IsActive": false}`)
	wUpdate := performRequest(testRouter, "PUT", "/api/extensions/"+created.ID.String(), []byte(updateExt), adminToken)
	if wUpdate.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", wUpdate.Code, wUpdate.Body.String())
	}

	// Delete extension
	wDel := performRequest(testRouter, "DELETE", "/api/extensions/"+created.ID.String(), nil, adminToken)
	if wDel.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", wDel.Code)
	}

	// Delete non-existent
	wDel404 := performRequest(testRouter, "DELETE", "/api/extensions/"+created.ID.String(), nil, adminToken)
	if wDel404.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", wDel404.Code)
	}
}

func TestRoleAuthorizationReadOnlyViewer(t *testing.T) {
	setupTestEnvironment(t)

	// Viewer can read extensions
	wGet := performRequest(testRouter, "GET", "/api/extensions", nil, viewerToken)
	if wGet.Code != http.StatusOK {
		t.Fatalf("viewer expected 200 GET, got %d", wGet.Code)
	}

	// Viewer cannot create extensions
	newExt := `{"ExtensionNumber": "1099", "PasswordHash": "pass1099"}`
	wPost := performRequest(testRouter, "POST", "/api/extensions", []byte(newExt), viewerToken)
	if wPost.Code != http.StatusForbidden {
		t.Fatalf("viewer expected 403 Forbidden for POST, got %d: %s", wPost.Code, wPost.Body.String())
	}
}

func TestTenantIsolationAndIDOREradication(t *testing.T) {
	setupTestEnvironment(t)

	// Create user in Tenant B
	uB := User{
		ID:           uuid.New(),
		TenantID:     tenantB,
		Email:        "admin@beta.local",
		PasswordHash: "pass",
		Role:         "admin",
	}
	DB.Create(&uB)

	tokB := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       uB.ID.String(),
		"tenant_id": tenantB.String(),
		"role":      "admin",
		"exp":       time.Now().Add(time.Hour).Unix(),
	})
	tenantBToken, _ := tokB.SignedString(jwtSecret)

	// Tenant B lists extensions - must NOT see extension 1001 (which belongs to Tenant A)
	wListB := performRequest(testRouter, "GET", "/api/extensions", nil, tenantBToken)
	var listRes map[string][]Extension
	json.Unmarshal(wListB.Body.Bytes(), &listRes)
	for _, e := range listRes["extensions"] {
		if e.ExtensionNumber == "1001" {
			t.Fatalf("IDOR vulnerability! Tenant B can see Tenant A's extension 1001!")
		}
	}

	// Tenant B attempts to delete Tenant A's extension (IDOR attack)
	var extA Extension
	DB.Where("extension_number = ?", "1001").First(&extA)
	wAttack := performRequest(testRouter, "DELETE", "/api/extensions/"+extA.ID.String(), nil, tenantBToken)
	if wAttack.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found on IDOR cross-tenant delete, got %d", wAttack.Code)
	}

	// Verify Tenant A extension was NOT deleted
	var checkExt Extension
	if err := DB.Where("id = ?", extA.ID).First(&checkExt).Error; err != nil {
		t.Fatalf("Extension 1001 was erroneously deleted by Tenant B!")
	}
}

func TestDashboardStats(t *testing.T) {
	setupTestEnvironment(t)

	w := performRequest(testRouter, "GET", "/api/dashboard/stats", nil, adminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var stats map[string]int64
	json.Unmarshal(w.Body.Bytes(), &stats)
	if stats["extensions"] < 1 {
		t.Fatalf("expected at least 1 extension, got %d", stats["extensions"])
	}
}

func TestCommitsEndpoint(t *testing.T) {
	setupTestEnvironment(t)

	w := performRequest(testRouter, "GET", "/api/commits", nil, adminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "commits") {
		t.Fatalf("expected commits list in body, got: %s", w.Body.String())
	}
}

func TestRingGroupsAndQueuesCRUD(t *testing.T) {
	setupTestEnvironment(t)

	// 1. Ring Groups Test
	groupPayload := `{
		"Name": "Technical Support Group",
		"InternalNumber": "601",
		"SiteName": "Main Corporate HQ",
		"RingCadence": "Default",
		"CallerNamePrefix": "Tech:",
		"NoAnswerTimeout": 15,
		"NoAnswerDestination": "Voicemail Box",
		"Members": [
			{"Member": "Alice Smith", "Extension": "1001", "DelaySec": 0},
			{"Member": "Bob Jones", "Extension": "1002", "DelaySec": 5}
		]
	}`
	wGrpCreate := performRequest(testRouter, "POST", "/api/ring-groups", []byte(groupPayload), adminToken)
	if wGrpCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 for ring group create, got %d: %s", wGrpCreate.Code, wGrpCreate.Body.String())
	}

	var createdGrp CallDistributionGroup
	json.Unmarshal(wGrpCreate.Body.Bytes(), &createdGrp)
	if createdGrp.Name != "Technical Support Group" || len(createdGrp.Members) != 2 {
		t.Fatalf("unexpected group content: %+v", createdGrp)
	}

	wGrpList := performRequest(testRouter, "GET", "/api/ring-groups", nil, adminToken)
	if wGrpList.Code != http.StatusOK || !strings.Contains(wGrpList.Body.String(), "Technical Support Group") {
		t.Fatalf("expected group in list, got: %s", wGrpList.Body.String())
	}

	// 2. Call Queues Test
	queuePayload := `{
		"Name": "Tier 2 Helpdesk",
		"InternalNumber": "801",
		"Extension": "8002",
		"MaxCapacity": 20,
		"RingStrategy": "Round Robin",
		"MemberRingTimeout": 25,
		"MemberRetryDelay": 10,
		"MemberWrapUpDelay": 5,
		"PlayMusicWhileWaiting": true,
		"Music": "Default Telecom Harmony",
		"CallerNamePrefix": "Tier2:",
		"NoAnswerTimeout": 30,
		"NoAnswerDestination": "Voicemail Box",
		"Members": [
			{"Member": "Alice Smith", "Extension": "1001", "Penalty": 0, "Disabled": false},
			{"Member": "Bob Jones", "Extension": "1002", "Penalty": 1, "Disabled": true}
		]
	}`
	wQCreate := performRequest(testRouter, "POST", "/api/queues", []byte(queuePayload), adminToken)
	if wQCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 for queue create, got %d: %s", wQCreate.Code, wQCreate.Body.String())
	}

	var createdQ CallQueue
	json.Unmarshal(wQCreate.Body.Bytes(), &createdQ)
	if createdQ.Name != "Tier 2 Helpdesk" || len(createdQ.Members) != 2 {
		t.Fatalf("unexpected queue content: %+v", createdQ)
	}

	wQList := performRequest(testRouter, "GET", "/api/queues", nil, adminToken)
	if wQList.Code != http.StatusOK || !strings.Contains(wQList.Body.String(), "Tier 2 Helpdesk") {
		t.Fatalf("expected queue in list, got: %s", wQList.Body.String())
	}

	// Clean up
	wDelGrp := performRequest(testRouter, "DELETE", "/api/ring-groups/"+createdGrp.ID.String(), nil, adminToken)
	if wDelGrp.Code != http.StatusOK {
		t.Fatalf("expected 200 deleting ring group, got %d", wDelGrp.Code)
	}
	wDelQ := performRequest(testRouter, "DELETE", "/api/queues/"+createdQ.ID.String(), nil, adminToken)
	if wDelQ.Code != http.StatusOK {
		t.Fatalf("expected 200 deleting queue, got %d", wDelQ.Code)
	}
}

func TestAutoCLIPLifecycleAndBackgroundWorker(t *testing.T) {
	setupTestEnvironment(t)

	// Create AutoCLIP entry
	clipPayload := `{
		"caller_number": "07987654321",
		"original_callee_ext": "1001"
	}`
	wCreate := performRequest(testRouter, "POST", "/api/autoclip", []byte(clipPayload), adminToken)
	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 for autoclip create, got %d: %s", wCreate.Code, wCreate.Body.String())
	}

	// Verify GET /api/autoclip returns records
	wGet := performRequest(testRouter, "GET", "/api/autoclip", nil, adminToken)
	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 for GET /api/autoclip, got %d", wGet.Code)
	}
	if !strings.Contains(wGet.Body.String(), "07987654321") {
		t.Fatalf("expected autoclip record in response: %s", wGet.Body.String())
	}

	// Test background maintenance cycle execution
	if err := RunMaintenanceCycle(); err != nil {
		t.Fatalf("expected maintenance cycle to succeed: %v", err)
	}
}

func TestDiagnosticsPing(t *testing.T) {
	setupTestEnvironment(t)

	pingBody := `{"host": "127.0.0.1"}`
	w := performRequest(testRouter, "POST", "/api/diagnostics/ping", []byte(pingBody), adminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for diagnostics ping, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "127.0.0.1") {
		t.Fatalf("expected host in response: %s", w.Body.String())
	}
}

func TestSIPSettingsCRUDAndValidation(t *testing.T) {
	setupTestEnvironment(t)

	// 1. GET /api/sip-settings (should auto-provision defaults)
	wGet := performRequest(testRouter, "GET", "/api/sip-settings", nil, adminToken)
	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 for GET /api/sip-settings, got %d: %s", wGet.Code, wGet.Body.String())
	}
	var settings SIPSettings
	if err := json.Unmarshal(wGet.Body.Bytes(), &settings); err != nil {
		t.Fatalf("failed to unmarshal SIP settings: %v", err)
	}
	if settings.SIPPort != 5060 || settings.RTPPortStart != 16384 || settings.RTPPortEnd != 32768 {
		t.Fatalf("unexpected default settings values: %+v", settings)
	}

	// 2. PUT /api/sip-settings with validation error (RTP start >= RTP end)
	badPayload := `{
		"user_agent": "Test-PBX",
		"sip_port": 5060,
		"sip_tls_port": 5061,
		"external_sip_port": 5080,
		"wss_port": 7443,
		"rtp_port_start": 30000,
		"rtp_port_end": 20000,
		"jitter_buffer_initial_ms": 60,
		"jitter_buffer_max_ms": 200,
		"min_se_sec": 120,
		"session_expires_sec": 1800,
		"register_min_sec": 60,
		"register_max_sec": 3600
	}`
	wBad := performRequest(testRouter, "PUT", "/api/sip-settings", []byte(badPayload), adminToken)
	if wBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for inverted RTP range, got %d: %s", wBad.Code, wBad.Body.String())
	}

	// 3. PUT /api/sip-settings with valid updates
	validPayload := `{
		"user_agent": "Custom-Carrier-PBX/3.0",
		"sip_port": 5060,
		"sip_tls_port": 5061,
		"external_sip_port": 5080,
		"wss_port": 7443,
		"rtp_port_start": 20000,
		"rtp_port_end": 40000,
		"rtp_timeout_sec": 45,
		"jitter_buffer_mode": "adaptive",
		"jitter_buffer_initial_ms": 80,
		"jitter_buffer_max_ms": 250,
		"dtmf_mode": "rfc2833",
		"audio_codecs": "OPUS,G722,PCMA",
		"enable_cng": true,
		"nat_mode": "manual",
		"external_ipv4": "198.51.100.25",
		"min_se_sec": 120,
		"session_expires_sec": 1800,
		"register_min_sec": 60,
		"register_max_sec": 3600,
		"register_default_sec": 1800,
		"dscp_audio": 46,
		"dscp_video": 34,
		"dscp_signaling": 26
	}`
	wPut := performRequest(testRouter, "PUT", "/api/sip-settings", []byte(validPayload), adminToken)
	if wPut.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid PUT /api/sip-settings, got %d: %s", wPut.Code, wPut.Body.String())
	}
	var updated SIPSettings
	json.Unmarshal(wPut.Body.Bytes(), &updated)
	if updated.UserAgent != "Custom-Carrier-PBX/3.0" || updated.RTPPortStart != 20000 || updated.RTPPortEnd != 40000 {
		t.Fatalf("updated values did not persist: %+v", updated)
	}

	// 4. POST /api/sip-settings/apply
	wApply := performRequest(testRouter, "POST", "/api/sip-settings/apply", nil, adminToken)
	if wApply.Code != http.StatusOK {
		t.Fatalf("expected 200 for apply, got %d: %s", wApply.Code, wApply.Body.String())
	}
	if !strings.Contains(wApply.Body.String(), "applied") {
		t.Fatalf("expected applied confirmation, got: %s", wApply.Body.String())
	}

	// 5. POST /api/sip-settings/reset
	wReset := performRequest(testRouter, "POST", "/api/sip-settings/reset", nil, adminToken)
	if wReset.Code != http.StatusOK {
		t.Fatalf("expected 200 for reset, got %d: %s", wReset.Code, wReset.Body.String())
	}
	var resetResp struct {
		Status   string      `json:"status"`
		Settings SIPSettings `json:"settings"`
	}
	json.Unmarshal(wReset.Body.Bytes(), &resetResp)
	if resetResp.Settings.RTPPortStart != 16384 || resetResp.Settings.RTPPortEnd != 32768 {
		t.Fatalf("expected default RTP ports after reset, got start=%d end=%d", resetResp.Settings.RTPPortStart, resetResp.Settings.RTPPortEnd)
	}
}

func TestFreeSWITCHSwitchConfGeneration(t *testing.T) {
	setupTestEnvironment(t)

	form := "section=configuration&key_value=switch.conf"
	req, _ := http.NewRequest("POST", "/fs", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from /fs switch.conf, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "rtp-start-port") || !strings.Contains(body, "rtp-end-port") {
		t.Fatalf("expected switch.conf with RTP ports, got: %s", body)
	}
}

func TestAIAssistantsCRUD(t *testing.T) {
	setupTestEnvironment(t)

	// 1. Create Assistant
	payload := map[string]interface{}{
		"name":             "Jessica - AI Receptionist",
		"extension_number": "891",
		"role_description": "Virtual Front Desk",
		"system_prompt":    "You are an AI receptionist. Help callers book appointments.",
		"greeting_message": "Hello, how can I direct your call?",
		"llm_provider":     "openai",
		"llm_model":        "gpt-4o",
		"voice_id":         "nova",
		"enable_booking":   true,
		"enable_transfer":  true,
		"transfer_target":  "1001",
		"is_active":        true,
	}
	wCreate := performRequest(testRouter, "POST", "/api/ai/assistants", payload, adminToken)
	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 for creating assistant, got %d: %s", wCreate.Code, wCreate.Body.String())
	}
	var created AIAssistant
	if err := json.Unmarshal(wCreate.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode created assistant: %v", err)
	}
	if created.ID == uuid.Nil || (created.ExtensionNumber != "891" && created.Extension != "891") {
		t.Fatalf("unexpected created assistant data: %+v", created)
	}

	// 2. List Assistants
	wList := performRequest(testRouter, "GET", "/api/ai/assistants", nil, adminToken)
	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 for listing assistants, got %d", wList.Code)
	}
	var listResp struct {
		Assistants []AIAssistant `json:"assistants"`
	}
	json.Unmarshal(wList.Body.Bytes(), &listResp)
	if len(listResp.Assistants) == 0 {
		t.Fatalf("expected at least 1 assistant, got %d", len(listResp.Assistants))
	}

	// 3. Update Assistant
	payload["greeting_message"] = "Welcome to MSIP Enterprise. How may I help?"
	wUpdate := performRequest(testRouter, "PUT", fmt.Sprintf("/api/ai/assistants/%s", created.ID), payload, adminToken)
	if wUpdate.Code != http.StatusOK {
		t.Fatalf("expected 200 for updating assistant, got %d", wUpdate.Code)
	}

	// 4. Delete Assistant
	wDel := performRequest(testRouter, "DELETE", fmt.Sprintf("/api/ai/assistants/%s", created.ID), nil, adminToken)
	if wDel.Code != http.StatusOK {
		t.Fatalf("expected 200 for deleting assistant, got %d", wDel.Code)
	}
}

func TestAIBookingsCRUD(t *testing.T) {
	setupTestEnvironment(t)

	// 1. Create Booking
	bookingData := map[string]interface{}{
		"caller_name":      "Alice Walker",
		"caller_number":    "+15559876543",
		"service_type":     "Consultation",
		"booking_date":     "2026-10-25",
		"booking_time":     "15:30",
		"duration_minutes": 30,
		"status":           "confirmed",
		"notes":            "Follow-up discussion",
	}
	wCreate := performRequest(testRouter, "POST", "/api/ai/bookings", bookingData, adminToken)
	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 for creating booking, got %d: %s", wCreate.Code, wCreate.Body.String())
	}
	var created AIBooking
	json.Unmarshal(wCreate.Body.Bytes(), &created)
	if created.ID == uuid.Nil || created.CallerName != "Alice Walker" {
		t.Fatalf("unexpected booking created: %+v", created)
	}

	// 2. List Bookings
	wList := performRequest(testRouter, "GET", "/api/ai/bookings", nil, adminToken)
	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 for listing bookings, got %d", wList.Code)
	}
	var listResp struct {
		Bookings []AIBooking `json:"bookings"`
	}
	json.Unmarshal(wList.Body.Bytes(), &listResp)
	if len(listResp.Bookings) == 0 {
		t.Fatalf("expected bookings list to contain at least 1 item")
	}

	// 3. Update Status
	bookingData["status"] = "completed"
	wUpdate := performRequest(testRouter, "PUT", fmt.Sprintf("/api/ai/bookings/%s", created.ID), bookingData, adminToken)
	if wUpdate.Code != http.StatusOK {
		t.Fatalf("expected 200 for updating booking, got %d", wUpdate.Code)
	}

	// 4. Delete Booking
	wDel := performRequest(testRouter, "DELETE", fmt.Sprintf("/api/ai/bookings/%s", created.ID), nil, adminToken)
	if wDel.Code != http.StatusOK {
		t.Fatalf("expected 200 for deleting booking, got %d", wDel.Code)
	}
}

func TestAISimulatorAndTools(t *testing.T) {
	setupTestEnvironment(t)

	// Seed assistant
	asst := AIAssistant{
		Base:            Base{ID: uuid.New()},
		TenantRef:       TenantRef{TenantID: tenantA},
		Name:            "Sarah Receptionist",
		Extension:       "890",
		ExtensionNumber: "890",
		EnableBooking:   true,
		EnableTransfer:  true,
		TransferTarget:  "1001",
		IsActive:        true,
	}
	DB.Create(&asst)

	// Test appointment booking tool invocation
	simPayload := map[string]interface{}{
		"assistant_id":  asst.ID.String(),
		"caller_number": "+15553334444",
		"messages": []map[string]string{
			{"role": "user", "content": "I would like to book an appointment tomorrow at 14:00 for Dental Service"},
		},
	}
	wSim := performRequest(testRouter, "POST", "/api/ai/simulate", simPayload, adminToken)
	if wSim.Code != http.StatusOK {
		t.Fatalf("expected 200 from AI simulator, got %d: %s", wSim.Code, wSim.Body.String())
	}
	var simResp struct {
		Reply     string       `json:"reply"`
		ToolCalls []AIToolCall `json:"tool_calls"`
	}
	json.Unmarshal(wSim.Body.Bytes(), &simResp)
	if simResp.Reply == "" {
		t.Fatalf("expected AI simulator reply, got empty string")
	}
	if len(simResp.ToolCalls) == 0 || simResp.ToolCalls[0].Name != "book_appointment" {
		t.Fatalf("expected tool call 'book_appointment', got: %+v", simResp.ToolCalls)
	}

	// Verify booking was persisted to DB
	var booking AIBooking
	if err := DB.Where("caller_number = ?", "+15553334444").First(&booking).Error; err != nil {
		t.Fatalf("expected appointment record to be created in DB: %v", err)
	}

	// Test call transfer tool invocation
	transferPayload := map[string]interface{}{
		"assistant_id":  asst.ID.String(),
		"caller_number": "+15553334444",
		"messages": []map[string]string{
			{"role": "user", "content": "Please transfer me to sales extension 1001"},
		},
	}
	wTransfer := performRequest(testRouter, "POST", "/api/ai/simulate", transferPayload, adminToken)
	if wTransfer.Code != http.StatusOK {
		t.Fatalf("expected 200 from AI simulator, got %d", wTransfer.Code)
	}
	var transferResp struct {
		ToolCalls []AIToolCall `json:"tool_calls"`
	}
	json.Unmarshal(wTransfer.Body.Bytes(), &transferResp)
	if len(transferResp.ToolCalls) == 0 || transferResp.ToolCalls[0].Name != "transfer_call" {
		t.Fatalf("expected tool call 'transfer_call', got: %+v", transferResp.ToolCalls)
	}
}

func TestAIVoicemailAndCDRTranscription(t *testing.T) {
	setupTestEnvironment(t)

	// 1. Voicemail Transcription
	vm := Voicemail{
		Base:      Base{ID: uuid.New()},
		TenantRef: TenantRef{TenantID: tenantA},
		Extension: "1001",
		CallerID:  "+15552221111",
		Duration:  42,
		FilePath:  "/var/lib/pbx/voicemail/default/1001/msg_0001.wav",
	}
	DB.Create(&vm)

	wVM := performRequest(testRouter, "POST", fmt.Sprintf("/api/voicemails/%s/transcribe", vm.ID), nil, adminToken)
	if wVM.Code != http.StatusOK {
		t.Fatalf("expected 200 for voicemail transcription, got %d: %s", wVM.Code, wVM.Body.String())
	}
	var vmResp struct {
		Transcription string `json:"transcription"`
		Summary       string `json:"summary"`
		Sentiment     string `json:"sentiment"`
	}
	json.Unmarshal(wVM.Body.Bytes(), &vmResp)
	if vmResp.Transcription == "" || vmResp.Summary == "" {
		t.Fatalf("expected non-empty transcription and summary, got: %+v", vmResp)
	}

	// 2. CDR Transcription
	cdr := CDR{
		Base:           Base{ID: uuid.New()},
		TenantID:       tenantA,
		UUID:           "test-call-uuid-987",
		CallerIDNumber: "+15552221111",
		CallerIDName:   "Test Caller",
		Destination:    "1001",
		Duration:       60,
		Billsec:        55,
		HangupCause:    "NORMAL_CLEARING",
	}
	DB.Create(&cdr)

	wCDR := performRequest(testRouter, "POST", fmt.Sprintf("/api/cdrs/%s/transcribe", cdr.ID), nil, adminToken)
	if wCDR.Code != http.StatusOK {
		t.Fatalf("expected 200 for CDR transcription, got %d: %s", wCDR.Code, wCDR.Body.String())
	}
	var cdrResp struct {
		Transcription string `json:"transcription"`
		Summary       string `json:"summary"`
		Sentiment     string `json:"sentiment"`
	}
	json.Unmarshal(wCDR.Body.Bytes(), &cdrResp)
	if cdrResp.Transcription == "" || cdrResp.Summary == "" {
		t.Fatalf("expected non-empty CDR transcription and summary, got: %+v", cdrResp)
	}
}

func TestFreeSWITCHAIDialplan(t *testing.T) {
	setupTestEnvironment(t)

	// Seed assistant with extension 890
	asst := AIAssistant{
		Base:            Base{ID: uuid.New()},
		TenantRef:       TenantRef{TenantID: tenantA},
		Name:            "Sarah AI",
		Extension:       "890",
		ExtensionNumber: "890",
		IsActive:        true,
	}
	DB.Create(&asst)

	form := "section=dialplan&context=default"
	req, _ := http.NewRequest("POST", "/fs", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from /fs dialplan, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "ai_assistant_") || !strings.Contains(body, "^890$") {
		t.Fatalf("expected FreeSWITCH dialplan to contain dynamic AI assistant extension 890, got: %s", body)
	}
	if !strings.Contains(body, "ws://pbx_backend:8081/api/ai/stream") {
		t.Fatalf("expected FreeSWITCH dialplan to contain AI WebSocket audio fork stream, got: %s", body)
	}
}

