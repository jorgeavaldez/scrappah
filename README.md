# Scrappah

A wireproxy wrapper with database-backed VPN configuration management.

## Project Structure

- `cmd/db/` - Database helper tool for managing VPN configurations
- `cmd/scrappah/` - Main server application
- `pkg/db/` - Database layer with Repository pattern
- `pkg/wp_config.go` - Wireproxy configuration parsing and validation

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
- **Config validation**: Validates wireproxy configs before storage using `ParseConfig()`
- **Context support**: All database operations use Go context for timeout/cancellation
- **Repository pattern**: Clean database abstraction layer

## Configuration

The database helper uses a local SQLite database (`local.db`) to store VPN configurations. Each config includes:

- Name (derived from filename)
- Active status
- Configuration content (validated wireproxy format)