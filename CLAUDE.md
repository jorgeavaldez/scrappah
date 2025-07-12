# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Scrappah is a wireproxy wrapper that manages multiple VPN configurations through a SQLite database. It automatically creates SOCKS5 proxies for all stored VPN configurations, enabling simultaneous WireGuard tunnel management with sequential port allocation starting from 8001.

## Architecture

The codebase follows a standard Go project structure:

- **cmd/scrappah/** - Main server application that loads VPN configs from database, dynamically injects SOCKS5 proxy sections, and starts WireGuard tunnels
- **cmd/db/** - Database utility CLI for adding, listing, and validating VPN configurations  
- **pkg/db/** - Repository pattern for SQLite database operations using libsql driver
- **pkg/wp_config.go** - Wireproxy configuration parser that validates WireGuard configs and handles INI format parsing

The main application (`cmd/scrappah/main.go`) uses a strings.Builder pattern to dynamically append `[Socks5]` sections to each VPN configuration before validation, enabling proxy functionality without modifying stored configs.

## Common Commands

### Build
```bash
make build          # Build both executables
make build-db       # Build database helper only  
make build-server   # Build main server only
```

### Database Management
```bash
make run-db ARGS="list"                        # List all VPN configs
make run-db ARGS="add /path/to/config.conf"    # Add VPN config from file
make run-db ARGS="revalidate"                  # Validate all stored configs
```

### Running and Testing
```bash
make run-server     # Start proxy server (creates SOCKS5 proxies from port 8001+)
make test-proxy     # Test connection through first proxy (port 8001)
```

### Development
```bash
make test           # Run all tests
make fmt            # Format Go code
make vet            # Run go vet
make lint           # Run golangci-lint (fallback to vet)
```

## Key Implementation Details

- **Dynamic Proxy Injection**: VPN configs are loaded from database and SOCKS5 sections are dynamically appended using strings.Builder before wireproxy validation
- **Sequential Port Allocation**: Each VPN config gets assigned a unique SOCKS5 port starting from 8001
- **Validation Pipeline**: All configs are validated through wireproxy parser before tunnel creation
- **Context-based Cancellation**: Uses context.WithCancel for graceful shutdown on SIGINT/SIGQUIT
- **Repository Pattern**: Database operations abstracted through Repository interface with proper connection management

The wireproxy integration requires configs to have `[Interface]` and `[Peer]` sections, with proxy sections (`[Socks5]`, `[http]`, etc.) added dynamically by scrappah rather than stored in the database.