# Terraform Drift Detection Agent

A cloud-agnostic platform that continuously compares Terraform state files against live cloud infrastructure to detect configuration drift — without requiring `terraform plan` or `terraform apply`.

## Features

- **Read-only** — uses cloud provider APIs directly, no Terraform operations needed
- **Fast** — concurrent resource fetching, results in seconds
- **Multi-interface** — CLI table, JSON output, web dashboard
- **Scheduled scans** — cron-based automated drift checks
- **Extensible** — provider interface makes adding GCP/Azure straightforward

### Supported AWS Resource Types (v1)

| Terraform Type | AWS Service |
|---|---|
| `aws_instance` | EC2 |
| `aws_s3_bucket` | S3 |
| `aws_security_group` | EC2 |
| `aws_vpc` | VPC |
| `aws_subnet` | VPC |
| `aws_iam_role` | IAM |
| `aws_lambda_function` | Lambda |
| `aws_db_instance` | RDS |
| `aws_rds_cluster` | RDS |

---

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- AWS credentials configured (`~/.aws/credentials`, `AWS_PROFILE`, or instance role)

### Install Go (macOS)

```bash
brew install go
```

---

## Quick Start

```bash
# 1. Clone and enter the repo
cd terraform-drift-detection-agent

# 2. Install dependencies
go mod tidy

# 3. Build the binary
go build -o drift ./cmd/drift

# 4. Run a scan against a local state file
./drift scan --state ./testdata/sample.tfstate --region us-east-1

# 5. Output as JSON
./drift scan --state ./testdata/sample.tfstate --output json

# 6. Start the web dashboard (http://localhost:8080)
./drift serve --state ./terraform.tfstate --port 8080

# 7. Run on a schedule (every 6 hours)
./drift schedule --state ./terraform.tfstate --cron "0 */6 * * *"
```

---

## Usage

### `drift scan` — on-demand scan

```
Usage: drift scan [flags]

Flags:
  --state   string   Path to terraform.tfstate (local) or s3://bucket/key
  --region  string   AWS region (default: us-east-1)
  --output  string   Output format: table or json (default: table)
```

#### Example output (table)

```
Terraform Drift Detection Report
State file : ./terraform.tfstate
Provider   : aws
Region     : us-east-1
Scanned at : 2024-01-15 10:23:45 UTC

Summary  Total: 12  |  Missing: 1  |  Modified: 2  |  Tag changed: 1  |  Clean: 8

+---------------------------+--------------------+-------------+-----------+------------------+----------+----------+
| Resource ID               | Type               | Name        | Drift     | Attribute        | Expected | Actual   |
+---------------------------+--------------------+-------------+-----------+------------------+----------+----------+
| i-0abc123def456789a       | aws_instance       | web_server  | MISSING   | -                | -        | -        |
| sg-0abc123def456789a      | aws_security_group | web_sg      | MODIFIED  | description      | old desc | new desc |
| my-app-assets-prod        | aws_s3_bucket      | assets      | TAG_CHANGED| tag:Environment | production| staging |
+---------------------------+--------------------+-------------+-----------+------------------+----------+----------+
```

### `drift serve` — web dashboard

```
Usage: drift serve [flags]

Flags:
  --state   string   State file for on-demand scans via the dashboard
  --region  string   AWS region (default: us-east-1)
  --port    int      HTTP port (default: 8080)
```

Dashboard endpoints:
- `GET /` — HTML dashboard
- `GET /api/reports` — JSON array of all scan reports
- `POST /api/scan` — trigger an on-demand scan

### `drift schedule` — cron scheduler

```
Usage: drift schedule [flags]

Flags:
  --state   string   Path to terraform.tfstate
  --region  string   AWS region (default: us-east-1)
  --cron    string   Cron expression, e.g. "0 */6 * * *"
```

---

## Remote State (S3)

```bash
drift scan --state s3://my-tfstate-bucket/prod/terraform.tfstate --region us-east-1
```

Requires S3 `GetObject` permission on the state bucket.

---

## Architecture

```
terraform-drift-detection-agent/
├── cmd/drift/main.go                   # Cobra CLI (scan / serve / schedule)
├── internal/
│   ├── state/                          # Terraform state parser (local + S3)
│   ├── normalizer/                     # Common resource model
│   ├── providers/
│   │   ├── provider.go                 # Provider interface
│   │   └── aws/                        # AWS adapter (EC2, S3, IAM, Lambda, RDS, VPC)
│   ├── detector/                       # Drift comparison engine
│   ├── reporter/                       # CLI table + JSON output
│   ├── scheduler/                      # Cron scheduler
│   └── dashboard/                      # HTTP server + embedded HTML dashboard
├── pkg/config/                         # Config struct
└── testdata/sample.tfstate             # Sample state for testing
```

### Extending to a new cloud provider

1. Create `internal/providers/gcp/provider.go` implementing the `providers.Provider` interface
2. Add resource fetchers per type
3. Register the provider in `cmd/drift/main.go`

No changes needed to the detector, reporter, or dashboard.

---

## Development

```bash
# Run tests
go test ./...

# Build for multiple platforms
GOOS=linux  GOARCH=amd64 go build -o drift-linux-amd64  ./cmd/drift
GOOS=darwin GOARCH=arm64 go build -o drift-darwin-arm64 ./cmd/drift
```
