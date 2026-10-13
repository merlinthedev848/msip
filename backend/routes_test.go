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

func performRequest(r *gin.Engine, method, path string, body []byte, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
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

	// Valid login
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
		"Name": "Newcastle Support",
		"InternalNumber": "995",
		"SiteName": "Agilico Systems",
		"RingCadence": "Default",
		"CallerNamePrefix": "Support (9)",
		"NoAnswerTimeout": 15,
		"NoAnswerDestination": "Front Door Failover",
		"Members": [
			{"Member": "Nicola Simpson", "Extension": "1001", "DelaySec": 0},
			{"Member": "Anne Mitchell Hall", "Extension": "1002", "DelaySec": 5}
		]
	}`
	wGrpCreate := performRequest(testRouter, "POST", "/api/ring-groups", []byte(groupPayload), adminToken)
	if wGrpCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 for ring group create, got %d: %s", wGrpCreate.Code, wGrpCreate.Body.String())
	}

	var createdGrp CallDistributionGroup
	json.Unmarshal(wGrpCreate.Body.Bytes(), &createdGrp)
	if createdGrp.Name != "Newcastle Support" || len(createdGrp.Members) != 2 {
		t.Fatalf("unexpected group content: %+v", createdGrp)
	}

	wGrpList := performRequest(testRouter, "GET", "/api/ring-groups", nil, adminToken)
	if wGrpList.Code != http.StatusOK || !strings.Contains(wGrpList.Body.String(), "Newcastle Support") {
		t.Fatalf("expected group in list, got: %s", wGrpList.Body.String())
	}

	// 2. Call Queues Test
	queuePayload := `{
		"Name": "Tier 2 Helpdesk",
		"InternalNumber": "992",
		"Extension": "8002",
		"MaxCapacity": 20,
		"RingStrategy": "Round Robin",
		"MemberRingTimeout": 25,
		"MemberRetryDelay": 10,
		"MemberWrapUpDelay": 5,
		"PlayMusicWhileWaiting": true,
		"Music": "Harmony",
		"CallerNamePrefix": "Tier2",
		"NoAnswerTimeout": 30,
		"NoAnswerDestination": "Voicemail Box",
		"Members": [
			{"Member": "Sharon Collinson", "Extension": "1002", "Penalty": 0, "Disabled": false},
			{"Member": "Julie Schubeler", "Extension": "1003", "Penalty": 1, "Disabled": true}
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
