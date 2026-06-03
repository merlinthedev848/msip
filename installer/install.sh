#!/bin/bash
set -e

echo "Starting PBX System Update/Install..."

# Pull latest images
docker-compose pull

# Restart services
docker-compose up -d --remove-orphans

echo "Update complete! Services are running."
docker ps
