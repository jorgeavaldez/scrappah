# Scrappah

A wireproxy wrapper with database-backed VPN configuration management that automatically creates SOCKS5 proxies for all stored VPN configurations.

## Project Structure

- `cmd/db/` - Database helper tool for managing VPN configurations
- `cmd/scrappah/` - Main server application
- `pkg/db/` - Database layer with Repository pattern
- `pkg/wp_config.go` - Wireproxy configuration parsing and validation

## Quick Start

```bash
# Build and add VPN configurations
make build
make run-db ARGS="add /path/to/your/wireguard.conf"

# Start the proxy server (creates SOCKS5 proxies starting from port 8001)
make run-server

# Test your connection through the first proxy
make test-proxy
```

## Build

```bash
make build          # Build both executables
make build-db       # Build database helper only
make build-server   # Build main server only
```

## Database Helper Usage

The database helper (`cmd/db/`) provides commands for managing VPN configurations:

### List all VPN configs
```bash
make run-db ARGS="list"
```

### Add VPN config from file
```bash
make run-db ARGS="add /path/to/config.conf"
```

### Validate all stored configs
```bash
make run-db ARGS="revalidate"
```

## Main Server Usage

The main server (`cmd/scrappah/`) automatically:

1. **Loads all VPN configurations** from the database
2. **Dynamically adds SOCKS5 proxy sections** to each config (starting from port 8001)
3. **Starts multiple WireGuard tunnels** - one for each valid configuration
4. **Creates SOCKS5 proxies** accessible on sequential ports

### Testing Proxy Connections

```bash
# Test IP through first proxy (port 8001)
make test-proxy

# Test with specific proxy port
curl --socks5 localhost:8001 https://api.ipify.org
curl --socks5 localhost:8002 https://api.ipify.org
# ... etc for each configured VPN
```

## Development

### Format and lint code
```bash
make fmt            # Format Go code
make vet            # Run go vet
make lint           # Run golangci-lint (or vet if not available)
```

### Run tests
```bash
make test           # Run all tests
make test-verbose   # Run tests with verbose output
make test-coverage  # Run tests with coverage
```

### Clean up
```bash
make clean          # Remove built binaries
```

## Features

- **Database-backed config storage**: Store and manage multiple VPN configurations
- **Automatic SOCKS5 proxy creation**: Dynamically adds SOCKS5 proxy sections to each VPN config
- **Multi-tunnel support**: Runs multiple WireGuard tunnels simultaneously (one per VPN config)
- **Sequential port allocation**: Automatically assigns unique ports starting from 8001

## Configuration

Uses a local SQLite database (`local.db`) to store VPN configurations. Each config includes:

- Name (derived from filename)
- Active status
- Configuration content (validated wireproxy format)

## Credits
This is really just a small wrapper around [wireproxy](https://github.com/pufferffish/wireproxy) to make it easier to manage VPN configurations and spin up many. All credits to them!
