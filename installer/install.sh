#!/usr/bin/env bash
set -euo pipefail

# Determine repository root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${PROJECT_ROOT}"

echo "=========================================================="
echo "           MSIP Telephony Platform Installer             "
echo "=========================================================="
echo "Project Root: ${PROJECT_ROOT}"

# Detect Docker Compose command (v2 plugin vs v1 standalone)
if docker compose version &>/dev/null; then
    COMPOSE_CMD="docker compose"
elif command -v docker-compose &>/dev/null; then
    COMPOSE_CMD="docker-compose"
else
    echo "ERROR: Neither 'docker compose' nor 'docker-compose' was found."
    echo "Please install Docker and Docker Compose before running this script."
    exit 1
fi
echo "Using Docker Compose command: ${COMPOSE_CMD}"

# Create default .env if missing
if [ ! -f .env ]; then
    echo "Creating default .env configuration..."
    cat <<EOF > .env
DB_USER=pbxuser
DB_PASS=pbxpass
DB_NAME=pbx
EXTERNAL_IP=127.0.0.1
EOF
fi

# Ensure telephony-core certificates directory exists
CERTS_DIR="${PROJECT_ROOT}/telephony-core/certs"
mkdir -p "${CERTS_DIR}"

if [ ! -f "${CERTS_DIR}/wss.pem" ]; then
    echo "Generating WebRTC WSS TLS certificates..."
    if command -v openssl &>/dev/null; then
        openssl req -x509 -newkey rsa:4096 \
            -keyout "${CERTS_DIR}/key.pem" \
            -out "${CERTS_DIR}/cert.pem" \
            -days 365 -nodes \
            -subj "/CN=pbx.local" \
            -addext "subjectAltName=DNS:localhost,DNS:pbx.local,IP:127.0.0.1" 2>/dev/null || \
        openssl req -x509 -newkey rsa:4096 \
            -keyout "${CERTS_DIR}/key.pem" \
            -out "${CERTS_DIR}/cert.pem" \
            -days 365 -nodes \
            -subj "/CN=pbx.local"
        cat "${CERTS_DIR}/cert.pem" "${CERTS_DIR}/key.pem" > "${CERTS_DIR}/wss.pem"
        echo "WebRTC TLS certificates generated successfully."
    else
        echo "WARNING: openssl not found in PATH; skipping certificate generation."
    fi
fi

# Build local containers and pull external dependencies
echo "Building local services (Backend & Frontend)..."
${COMPOSE_CMD} build --pull

echo "Pulling latest upstream service images..."
${COMPOSE_CMD} pull --ignore-pull-failures || true

echo "Starting all telephony containers..."
${COMPOSE_CMD} up -d --remove-orphans

echo ""
echo "=========================================================="
echo "   MSIP Telephony Platform Deployed Successfully!         "
echo "=========================================================="
echo " Web Management Dashboard : http://localhost"
echo " API Gateway (KrakenD)    : http://localhost:8080"
echo " SIP Signaling (Kamailio) : UDP/TCP 5060"
echo " FreeSWITCH Media / WSS   : TCP 5066 / WSS 7443"
echo " Homer SIP Packet Capture : http://localhost:9080"
echo " Turn / STUN Server       : UDP 3478"
echo " Default Superadmin       : admin@pbx.local / admin123"
echo "=========================================================="
echo ""
${COMPOSE_CMD} ps
