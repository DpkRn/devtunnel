#!/usr/bin/env bash
# One-time self-signed cert for local/EC2 testing (browsers will warn). Replace with Let's Encrypt for production.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SSL="$ROOT/nginx/ssl"
mkdir -p "$SSL"
if [[ -f "$SSL/fullchain.pem" && -f "$SSL/privkey.pem" ]]; then
  echo "SSL files already exist in nginx/ssl/"
  exit 0
fi
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout "$SSL/privkey.pem" \
  -out "$SSL/fullchain.pem" \
  -subj "/CN=devtunnel"
echo "Wrote $SSL/fullchain.pem and privkey.pem"

sudo apt-get update
sudo apt-get install -y certbot
certbot --version

sudo certbot certonly --webroot -w /home/ubuntu/devtunnel/nginx/certbot \
  -d clickly.cv -d www.clickly.cv

sudo cp /etc/letsencrypt/live/clickly.cv/fullchain.pem /home/ubuntu/devtunnel/nginx/ssl/
sudo cp /etc/letsencrypt/live/clickly.cv/privkey.pem   /home/ubuntu/devtunnel/nginx/ssl/
cd /home/ubuntu/devtunnel && docker compose up -d --build