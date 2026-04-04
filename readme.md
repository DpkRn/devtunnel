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

## Docker (tunnel server)

On EC2 (or any host with Docker), after cloning the repo:

```bash
chmod +x scripts/docker-server.sh
./scripts/docker-server.sh
```

That stops/removes any existing container and image named `devtunnel-server`, rebuilds, and runs detached with `--restart unless-stopped` on ports **3000** and **9000**.

Manual one-off:

```bash
docker build -t devtunnel-server .
docker run --rm -p 3000:3000 -p 9000:9000 devtunnel-server
```

- **3000** — public HTTP edge (subdomain routing)
- **9000** — tunnel control (yamux; `mytunnel` clients dial this)

Override the Go image version if `go.mod` needs a newer toolchain:

```bash
docker build --build-arg GO_VERSION=1.25 -t devtunnel-server .
```
