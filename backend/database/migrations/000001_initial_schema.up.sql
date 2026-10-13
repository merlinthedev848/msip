-- MSIP Multi-Tenant PBX Complete Initial Schema

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 1. Tenants (Organizations & Domains)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Users (Tenant Admins & Superadmins)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'admin',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);

-- 3. Extensions
CREATE TABLE IF NOT EXISTS extensions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    extension_number VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    provisioning_code VARCHAR(8) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, extension_number)
);
CREATE INDEX IF NOT EXISTS idx_extensions_tenant_id ON extensions(tenant_id);

-- 4. Trunks (Carrier Gateways)
CREATE TABLE IF NOT EXISTS trunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    provider_name VARCHAR(100) NOT NULL,
    sip_server VARCHAR(255) NOT NULL,
    auth_method VARCHAR(50) NOT NULL,
    username VARCHAR(100),
    password VARCHAR(100),
    register BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_trunks_tenant_id ON trunks(tenant_id);

-- 5. Inbound Routes
CREATE TABLE IF NOT EXISTS inbound_routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    did_number VARCHAR(50) NOT NULL,
    destination_type VARCHAR(50) NOT NULL,
    destination_target VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, did_number)
);
CREATE INDEX IF NOT EXISTS idx_inbound_routes_tenant_id ON inbound_routes(tenant_id);

-- 6. Outbound Routes
CREATE TABLE IF NOT EXISTS outbound_routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    dial_pattern VARCHAR(100) NOT NULL,
    trunk_id UUID NOT NULL REFERENCES trunks(id) ON DELETE CASCADE,
    priority INT DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_outbound_routes_tenant_id ON outbound_routes(tenant_id);

-- 7. Call Detail Records (CDRs)
CREATE TABLE IF NOT EXISTS cdrs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE SET NULL,
    caller_id_name VARCHAR(100),
    caller_id_number VARCHAR(50),
    destination VARCHAR(100),
    context VARCHAR(50),
    duration INT DEFAULT 0,
    billsec INT DEFAULT 0,
    hangup_cause VARCHAR(100),
    uuid VARCHAR(100) UNIQUE,
    transcription TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_cdrs_tenant_id ON cdrs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cdrs_created_at ON cdrs(created_at DESC);

-- 8. IVR Menus
CREATE TABLE IF NOT EXISTS ivr_menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    extension VARCHAR(50) NOT NULL,
    greeting VARCHAR(255) NOT NULL,
    timeout INT DEFAULT 10000,
    max_failures INT DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, extension)
);
CREATE INDEX IF NOT EXISTS idx_ivr_menus_tenant_id ON ivr_menus(tenant_id);

-- 9. IVR Choices
CREATE TABLE IF NOT EXISTS ivr_choices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_id UUID NOT NULL REFERENCES ivr_menus(id) ON DELETE CASCADE,
    digits VARCHAR(10) NOT NULL,
    action VARCHAR(50) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_ivr_choices_menu_id ON ivr_choices(menu_id);

-- 10. Voicemail
CREATE TABLE IF NOT EXISTS voicemails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    extension VARCHAR(50) NOT NULL,
    caller_id VARCHAR(100),
    duration INT DEFAULT 0,
    file_path VARCHAR(255) NOT NULL,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_voicemails_tenant_id ON voicemails(tenant_id);

-- 11. Video Rooms
CREATE TABLE IF NOT EXISTS video_rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    room_code VARCHAR(50) UNIQUE NOT NULL,
    max_users INT DEFAULT 10,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_video_rooms_tenant_id ON video_rooms(tenant_id);

-- 12. IP Cameras
CREATE TABLE IF NOT EXISTS ip_cameras (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    rtsp_url VARCHAR(255) NOT NULL,
    location VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_ip_cameras_tenant_id ON ip_cameras(tenant_id);

-- 13. Team Chats
CREATE TABLE IF NOT EXISTS team_chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_team_chats_tenant_id ON team_chats(tenant_id);

-- 14. SMS Messages
CREATE TABLE IF NOT EXISTS sms_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    direction VARCHAR(20) NOT NULL,
    sender VARCHAR(50) NOT NULL,
    receiver VARCHAR(50) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_sms_messages_tenant_id ON sms_messages(tenant_id);

-- 15. API Credentials
CREATE TABLE IF NOT EXISTS api_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    token_key VARCHAR(100) UNIQUE NOT NULL,
    secret_hash VARCHAR(64),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_api_credentials_tenant_id ON api_credentials(tenant_id);

-- 16. Webhooks
CREATE TABLE IF NOT EXISTS webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    event VARCHAR(100) NOT NULL,
    target_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_webhooks_tenant_id ON webhooks(tenant_id);

-- 17. Throttle Rules
CREATE TABLE IF NOT EXISTS throttle_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    rate INT DEFAULT 60,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_throttle_rules_tenant_id ON throttle_rules(tenant_id);

-- 18. Firewall Rules
CREATE TABLE IF NOT EXISTS firewall_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    ip_address VARCHAR(50) NOT NULL,
    action VARCHAR(20) NOT NULL,
    notes VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_firewall_rules_tenant_id ON firewall_rules(tenant_id);

-- 19. Fraud Rules
CREATE TABLE IF NOT EXISTS fraud_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    prefix VARCHAR(50) NOT NULL,
    max_cost DOUBLE PRECISION DEFAULT 0,
    block_action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_fraud_rules_tenant_id ON fraud_rules(tenant_id);
