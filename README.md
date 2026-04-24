# driftwatch

A CLI tool that detects configuration drift between deployed services and their declared infrastructure-as-code definitions.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at your IaC definition and a target environment to scan for drift:

```bash
driftwatch scan --config ./infra/prod.yaml --env production
```

Compare a specific service against its Terraform state:

```bash
driftwatch scan --provider terraform --state ./terraform.tfstate --service api-gateway
```

Output drift report as JSON:

```bash
driftwatch scan --config ./infra/prod.yaml --env staging --output json
```

### Flags

| Flag | Description |
|------|-------------|
| `--config` | Path to IaC definition file |
| `--env` | Target environment to inspect |
| `--provider` | IaC provider (`terraform`, `pulumi`, `cloudformation`) |
| `--output` | Output format: `text` (default), `json`, `yaml` |
| `--service` | Limit scan to a specific service |

---

## License

[MIT](LICENSE)