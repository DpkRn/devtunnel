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

**80** → nginx → mytunneld **:3000** (HTTP, `Host` kept for subdomains). **9000** → mytunneld (tunnel control), not nginx.

```bash
chmod +x scripts/docker-server.sh && ./scripts/docker-server.sh
# same as: docker compose up -d --build
```

Set `PublicHostSuffix` in `internal/config/config.go` to your domain when using port 80.

Plain Docker (no nginx):

```bash
docker build -t mytunneld .
docker run --rm -p 3000:3000 -p 9000:9000 mytunneld
```

The mytunneld image is **distroless** (~10 MB). `docker image prune -f` clears dangling build layers.
