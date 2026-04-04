# devtunnel

Expose local ports to the internet using a simple CLI.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/DpkRn/devtunnel/master/install.sh | bash
```

This will download the latest `mytunnel` binary for your OS (Linux or macOS) and install it to `/usr/local/bin`.

## Usage

```bash
mytunnel http 3000
```

## Build from Source

```bash
go build -o mytunnel ./cmd/client
sudo mv mytunnel /usr/local/bin/
```

For a completely fresh build:

```bash
go clean -cache -modcache -i -r
go build -a -o mytunnel ./cmd/client
```

## Docker (tunnel server + nginx)

- **443** — HTTPS (nginx → mytunneld :3000, `X-Forwarded-Proto: https`).
- **80** — redirects to HTTPS; `/.well-known/acme-challenge/` is served for Let’s Encrypt.
- **9000** — tunnel control (mytunneld), not nginx.

First run creates a **self-signed** cert in `nginx/ssl/` (browser warning). Replace with real certs (`fullchain.pem`, `privkey.pem`) for production. **Let’s Encrypt (webroot):** with the stack up, `sudo certbot certonly --webroot -w /home/ubuntu/devtunnel/nginx/certbot -d yourdomain.com -d '*.yourdomain.com'` (wildcard needs DNS challenge), then copy/symlink PEMs into `nginx/ssl/` and `docker compose restart nginx`.

Production domain is set in `internal/config/config.go` (`PublicHostSuffix`, `PublicURLScheme`). The CLI dials **`clickly.cv:9000`** for the tunnel; override with `DEVTUNNEL_SERVER=localhost:9000` for local dev.

```bash
chmod +x scripts/docker-server.sh && ./scripts/docker-server.sh
```

The script uses **`docker compose`** (plugin) if available, otherwise **`docker-compose`**.

Ubuntu’s default apt **does not** ship `docker-compose-plugin`. Install the Compose v2 binary once:

```bash
chmod +x scripts/install-docker-compose.sh && ./scripts/install-docker-compose.sh
```

Or add [Docker’s official apt repo](https://docs.docker.com/engine/install/ubuntu/) and install `docker-compose-plugin` from there.

Plain Docker (no nginx):

```bash
docker build -t mytunneld .
docker run --rm -p 3000:3000 -p 9000:9000 mytunneld
```

The mytunneld image is **distroless** (~10 MB). `docker image prune -f` clears dangling build layers.
