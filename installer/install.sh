#!/bin/bash
set -e

echo "Starting PBX System Update/Install..."

# Generate Self-Signed Certificates for FreeSWITCH WebRTC
mkdir -p ../telephony-core/certs
if [ ! -f ../telephony-core/certs/wss.pem ]; then
  echo "Generating self-signed WSS certificate..."
  openssl req -x509 -newkey rsa:4096 -keyout ../telephony-core/certs/key.pem -out ../telephony-core/certs/cert.pem -days 365 -nodes -subj "/CN=pbx.local"
  cat ../telephony-core/certs/cert.pem ../telephony-core/certs/key.pem > ../telephony-core/certs/wss.pem
fi

# Pull latest images
docker-compose pull

# Restart services
docker-compose up -d --remove-orphans

echo "Update complete! Services are running."
docker ps
