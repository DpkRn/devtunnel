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

---

## Let’s Encrypt wildcard (`*.clickly.cv`) with Namecheap — DNS-01 flow

**Why DNS-01:** A wildcard cert **must** use the **DNS-01** challenge. **HTTP-01** (webroot on port 80) **cannot** issue `*.clickly.cv`.

**What happens:** Let’s Encrypt asks you to prove you control DNS for `clickly.cv` by publishing a **TXT** record at:

`_acme-challenge.clickly.cv`

(with a **value** they give you — it changes each issuance/renewal). After their servers see the TXT, they issue the cert.

### A) Fully manual (understand the flow; no Namecheap API)

1. On the server (or any machine with certbot):

   ```bash
   sudo certbot certonly --manual --preferred-challenges dns \
     -d clickly.cv -d '*.clickly.cv' \
     --email you@example.com --agree-tos
   ```

2. Certbot prints something like: **add TXT** at `_acme-challenge.clickly.cv` with **value** `xxxxxxxx`.

3. In **Namecheap** → **Domain List** → **Manage** `clickly.cv` → **Advanced DNS** → **Add New Record**:
   - **Type:** `TXT Record`
   - **Host:** `_acme-challenge` (Namecheap often wants this **without** the domain suffix; the UI builds `clickly.cv` for you)
   - **Value:** paste the string certbot gave you (quotes if it says so)
   - **TTL:** 1 min or Automatic

4. Wait for DNS to propagate (often 1–15+ minutes). Check:

   ```bash
   dig TXT _acme-challenge.clickly.cv +short
   ```

5. Press **Enter** in certbot when the TXT is visible.

6. Certs land under `/etc/letsencrypt/live/clickly.cv/` (`fullchain.pem`, `privkey.pem`).

7. Copy into this repo’s nginx mount and reload:

   ```bash
   sudo cp /etc/letsencrypt/live/clickly.cv/fullchain.pem /path/to/devtunnel/nginx/ssl/
   sudo cp /etc/letsencrypt/live/clickly.cv/privkey.pem   /path/to/devtunnel/nginx/ssl/
   cd /path/to/devtunnel && docker compose restart nginx
   ```

**Renewal:** Wildcard + manual = you must repeat DNS TXT when certbot renews, unless you automate (below).

### B) Automated with Namecheap API + `certbot-dns-namecheap`

Namecheap only enables the DNS API if your account meets their rules (e.g. enough domains or purchase history — see [Namecheap API](https://www.namecheap.com/support/api/intro/)).

1. **Namecheap:** Profile → **Tools** → **Business & Dev Tools** → **Namecheap API Access** → turn on API, note **API key**, set **whitelist** to your **EC2 public IP**.

2. **Credentials file** (mode `600`):

   ```ini
   # /root/.secrets/namecheap.ini
   dns_namecheap_username=yournamecheapusername
   dns_namecheap_api_key=yourapikey
   ```

3. Install plugin (example with pip / venv on Ubuntu):

   ```bash
   sudo apt-get install -y python3-pip
   sudo pip install --break-system-packages certbot-dns-namecheap
   ```

   (Use a venv if you prefer not to use `--break-system-packages`.)

4. Issue:

   ```bash
   sudo certbot certonly \
     --authenticator dns-namecheap \
     --dns-namecheap-credentials /root/.secrets/namecheap.ini \
     -d clickly.cv -d '*.clickly.cv' \
     --email you@example.com --agree-tos
   ```

5. Same copy to `nginx/ssl/` and `docker compose restart nginx` as above.

6. **Renewal:** `certbot renew` with the same authenticator (often via cron/systemd timer).

### DNS records you still need for the app

Wildcard TLS only fixes the **lock**. For tunnels you also need:

- **A** record: `*` → EC2 IP (wildcard hostnames), and/or **A** for `@` / `www` as you like.

### Summary

| Step | Action |
|------|--------|
| 1 | Request cert with **both** `clickly.cv` and `*.clickly.cv` |
| 2 | Prove control via **TXT** at `_acme-challenge.clickly.cv` (DNS-01) |
| 3 | Install **fullchain.pem** + **privkey.pem** into **nginx/ssl/** |
| 4 | **Restart nginx** container |
