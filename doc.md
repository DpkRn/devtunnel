# devtunnel — Command Reference

Everything used to build, release, and maintain this project.

---

## Go Build Commands

### Build without cache (`-a`)

```bash
GOOS=darwin GOARCH=arm64 go build -a -o mytunnel-mac-arm64 ./cmd/client
```

| Part | Meaning |
|------|---------|
| `GOOS=darwin` | Target OS: macOS |
| `GOARCH=arm64` | Target CPU: Apple Silicon (M1/M2/M3) |
| `go build` | Compile the Go program |
| `-a` | Force rebuild of all packages — skips the build cache |
| `-o mytunnel-mac-arm64` | Output binary filename |
| `./cmd/client` | Package to build (entry point) |

### All three platform builds

```bash
# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -a -o mytunnel-mac-arm64 ./cmd/client

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -a -o mytunnel-mac ./cmd/client

# Linux x86_64
GOOS=linux GOARCH=amd64 go build -a -o mytunnel-linux ./cmd/client
```

### Full clean rebuild (nuclear option)

```bash
go clean -cache -modcache -i -r
go build -a -o mytunnel ./cmd/client
```

| Flag | Meaning |
|------|---------|
| `-cache` | Delete the build cache |
| `-modcache` | Delete the downloaded module cache |
| `-i` | Remove installed packages |
| `-r` | Apply recursively to all dependencies |

---

## Install Locally

```bash
sudo cp mytunnel-mac-arm64 /usr/local/bin/mytunnel
sudo chmod +x /usr/local/bin/mytunnel
```

| Command | Meaning |
|---------|---------|
| `sudo cp` | Copy with root privileges |
| `/usr/local/bin/` | Standard location for user-installed binaries — already in `$PATH` |
| `chmod +x` | Mark the file as executable |

---

## GitHub CLI (`gh`)

### Install

```bash
brew install gh
```

### Authenticate

```bash
gh auth login
```

Follow the prompts: GitHub.com → HTTPS → Login with a web browser.

### Upload binaries to a release

```bash
gh release upload devtunnel mytunnel-mac mytunnel-mac-arm64 mytunnel-linux --clobber
```

| Part | Meaning |
|------|---------|
| `gh release upload` | Upload assets to an existing GitHub release |
| `devtunnel` | The release tag to upload to |
| `mytunnel-mac mytunnel-mac-arm64 mytunnel-linux` | Files to upload |
| `--clobber` | Overwrite existing assets with the same name |

### Create a new release

```bash
gh release create v0.2.0 mytunnel-mac mytunnel-mac-arm64 mytunnel-linux \
  --title "v0.2.0" \
  --notes "Release notes here"
```

---

## Git Commands

### Check status

```bash
git status
```

Shows modified, staged, and untracked files.

### Stage and commit

```bash
git add .
git commit -m "your message"
```

### Push to GitHub

```bash
git push origin master
```

### View commit history

```bash
git log --oneline
```

---

## Install Script (one-liner)

```bash
curl -fsSL https://raw.githubusercontent.com/DpkRn/devtunnel/master/install.sh | bash
```

| Flag | Meaning |
|------|---------|
| `-f` | Fail silently on HTTP errors (non-zero exit) |
| `-s` | Silent mode — no progress bar |
| `-S` | Show error even in silent mode |
| `-L` | Follow redirects |

The script detects OS and CPU architecture automatically:
- `uname` → gets the OS (`Darwin` or `Linux`)
- `uname -m` → gets the CPU arch (`arm64` or `x86_64`)
- Downloads the correct binary and moves it to `/usr/local/bin/`

---

## Debugging a Bad Install

```bash
# Check what's actually installed
file /usr/local/bin/mytunnel

# Check CPU architecture of your Mac
uname -m

# Manually download and inspect before installing
curl -fsSL <URL> -o /tmp/test-binary
file /tmp/test-binary
xxd /tmp/test-binary | head -3
```

---

## Docker

The `Dockerfile` builds **mytunneld** (`./cmd/server`) as a static binary on **`/mytunneld`** in a **distroless** runtime image.

**Recommended:** `docker-compose.yml` runs **mytunneld** + **nginx** (Alpine). Nginx listens on **80**, proxies to **mytunneld:3000**, and forwards **`Host`** so subdomain routing in Go still works. **9000** stays on mytunneld for tunnel clients.

| Port (host) | Service | Role |
|-------------|---------|------|
| 80 | nginx | Public HTTP |
| 9000 | mytunneld | Control plane (yamux) |
| — | mytunneld:3000 | Edge HTTP (internal to Compose network only) |

```bash
./scripts/docker-server.sh   # same as: docker compose up -d --build
```

Without nginx (direct edge):

```bash
docker build -t mytunneld .
docker run --rm -p 3000:3000 -p 9000:9000 mytunneld
```

Build arg for Go toolchain:

```bash
docker build --build-arg GO_VERSION=1.25 -t mytunneld .
```

Nginx config lives in **`nginx/nginx.conf`**. `.dockerignore` trims build context.
