# Changelog

## [Unreleased]

## [0.1.0] - 2026-03-30

### Added

- Binary distribution for Linux and macOS
- Command-line argument parsing for client — supports `http <port>` command with usage instructions
- Initial project setup with client and server implementations, protocol definitions, and Go module files

### Fixed

- Added missing config values and read time limit
- Minor fix and removed overwriting of headers
- Removed byte reading size dependency

## [0.0.1] - 2026-03-28 / 2026-03-29

### Added

- Request and response protocol implementation
- Unique public URL generated per connection
- Install script
- Pre-built binaries
- README for tunnel client; server now outputs dynamic port forwarding info
- Minimal HTTP request forwarding (`go run server.go http <port>`)
- TCP tunneling support
- Initial TCP bidirectional persistent connection
