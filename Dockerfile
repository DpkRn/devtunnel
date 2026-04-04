# Tunnel server — multi-stage: compile in Go image, ship only the static binary.
#
# Final image (~10MB): no Alpine, no shell, no Go toolchain — only mytunneld + distroless base.
# During `docker build`, Docker may still cache the *builder* image (golang:…) on the host;
# that is not inside the pushed/runnable image. Remove leftovers: docker image prune -f
#
# Build:  docker build -t mytunneld .
# Run:    docker run -d --name mytunneld -p 3000:3000 -p 9000:9000 mytunneld

# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /mytunneld ./cmd/server

# Smallest practical runtime: no package manager, non-root (uid 65532).
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build --chown=65532:65532 /mytunneld /mytunneld

EXPOSE 3000/tcp 9000/tcp

ENTRYPOINT ["/mytunneld"]
