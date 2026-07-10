# AGENTS.md

## How this project works

This is a Go CLI application using the Cobra framework. The main purpose is to generate boilerplate code for new application features as sub-modules.

## Directory structure

- `cmd/` - CLI command entry points
  - `root.go` - Main CLI root command
  - `generate.go` - Code generation for new sub-modules
- `payment/`, `produk/`, `transaksi/` - Generated feature modules (each has controller/repo/service for MVP pattern)

## Getting started

### Install the CLI tool locally:
```bash
go install ./...  # Compiles the alat_CLI binary to $GOPATH/bin/alat_CLI
```

### Test the CLI:
```bash
alat_CLI help                # See available commands
alat_CLI generate -h         # See generate command options
alat_CLI generate cache-scan -h  # See cache-scan options
```

### Scan and clean cache files:
```bash
alat_CLI generate cache-scan                        # Scan default cache dirs, interactive delete
alat_CLI generate cache-scan --min-size 50           # Only show dirs >= 50 MB
alat_CLI generate cache-scan --dry-run               # Preview only, no deletion
alat_CLI generate cache-scan --all                   # Delete all without prompt
```

Scans `~/Library/Caches` (macOS) and `~/.cache` (Linux) for large cache directories. Displays sizes and lets you select which to delete.

### Generate a new feature:
```bash
# Example: Generate a new module called "auth" using PostgreSQL
alat_CLI generate -n auth -d postgres
# Creates: auth/controller.go, auth/repository.go

# Example: Generate a new module called "user" using MySQL  
alat_CLI generate -n user -d mysql
# Creates: user/controller.go, user/repository.go
```

## CLI generated files contain:

Each generated module includes basic structure:
```go
// controller.go - Request/response handler
// repository.go - Database access layer (depends on *sql.DB)
```

The repository.go includes:
- SQL driver configuration comment (postgres/mysql)
- DSN setup template comment
- Ready-to-use SQL DB connection setup

## Build and test process

**No tests exist in this repository** - this is a CLI tool generator, not a library.

**Build and install:**
```bash
cd /path/to/alat_CLI
make build && make install  # If available
# or:
go build -o alat_CLI ./...  # Create binary locally
```

**Features are added via `alat_CLI generate`** not via git pull/test cycles.

## Development workflow

1. Use `alat_CLI generate <name> -d postgres|mysql` to create new features
2. Complete generated code with business logic
3. Build with standard Go tooling
4. No CI/CD pipeline required (current repo setup minimal)

## Key Go tooling commands

```bash
# Build the binary locally
make build  # If available
# or:
go build -o alat_CLI ./...

# To install the binary to GOPATH/bin
make install  # If available  
# or:
go install ./...

# Run /tmp for testing the tool locally
/tmp/alat_CLI help
```

## Dependencies

- `github.com/spf13/cobra` v1.10.2 - CLI framework
- `database/sql` - Standard library SQL package
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/lib/pq` - PostgreSQL driver

## Generated module structure quirks

Cobra is used for code generation. Generated modules assume:
- You will create a separate database connection for each module
- Service objects will hold business logic layered between controller and repository
- No built-in migrations or connection pooling - you must create these

Generated database setup examples:
```go
// PostgreSQL (in generated repository.go):
// Driver: github.com/lib/pq
// dsn := "host=localhost user=postgres password=secret dbname=mydb sslmode=disable"
// db, err := sql.Open("postgres", dsn)

// MySQL (in generated repository.go):
// Driver: github.com/go-sql-driver/mysql  
// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=true"
// db, err := sql.Open("mysql", dsn)
```