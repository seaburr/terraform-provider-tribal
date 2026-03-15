# Tribal Terraform Provider

A Terraform provider for managing resources in [Tribal](https://github.com/cberndt/tribal), a credential and resource expiry tracking tool.

## Resources

| Resource | Description |
|---|---|
| `tribal_resource` | A tracked credential or certificate (API key, SSH key, TLS cert, etc.) |
| `tribal_team` | A team that resources can be grouped under |
| `tribal_admin_settings` | Organization-wide notification settings (singleton) |

## Provider Configuration

```hcl
provider "tribal" {
  host    = "http://localhost:8000"  # or TRIBAL_HOST env var
  api_key = "tribal_sk_..."         # or TRIBAL_API_KEY env var
}
```

| Attribute | Type | Description |
|---|---|---|
| `host` | string | Base URL of the Tribal API. Defaults to `http://localhost:8000`. |
| `api_key` | string | API key for authentication. **Required.** |

## Resource: `tribal_resource`

Tracks an expiring credential or certificate.

```hcl
resource "tribal_resource" "example" {
  name                    = "Production API Key"
  dri                     = "platform-team@example.com"
  type                    = "API Key"  # Certificate | API Key | SSH Key | Other
  expiration_date         = "2026-12-31"
  purpose                 = "Authenticates the payment service with Stripe"
  generation_instructions = "Log into Stripe dashboard > Developers > API keys > Create new key"
  slack_webhook           = "https://hooks.slack.com/services/TEAM/CHANNEL/WEBHOOK"

  # Optional
  secret_manager_link = "https://vault.example.com/secret/stripe-api-key"
  team_id             = tribal_team.platform.id
  certificate_url     = "https://api.example.com"  # poll endpoint for auto expiry refresh
  auto_refresh_expiry = true
}
```

### Arguments

| Argument | Required | Description |
|---|---|---|
| `name` | yes | Display name of the resource. |
| `dri` | yes | Directly Responsible Individual (email or team name). |
| `type` | yes | One of: `Certificate`, `API Key`, `SSH Key`, `Other`. |
| `expiration_date` | yes | Expiry date in `YYYY-MM-DD` format. |
| `purpose` | yes | What this credential is used for. |
| `generation_instructions` | yes | Steps to renew or regenerate this credential. |
| `slack_webhook` | yes | Slack webhook for expiration alerts. |
| `secret_manager_link` | no | URL or ARN pointing to this secret in a secret manager. |
| `team_id` | no | Team to assign this resource to. Defaults to the organization default team. |
| `certificate_url` | no | TLS endpoint URL to poll for automatic certificate expiry detection. |
| `auto_refresh_expiry` | no | When `true`, `expiration_date` is updated automatically by polling `certificate_url`. |

### Computed Attributes

| Attribute | Description |
|---|---|
| `id` | Numeric resource ID. |
| `public_key_pem` | PEM-encoded public certificate (if uploaded). |
| `created_at` | Creation timestamp. |
| `updated_at` | Last update timestamp. |

### Import

```shell
terraform import tribal_resource.example 42
```

## Resource: `tribal_team`

Groups resources under a named team.

```hcl
resource "tribal_team" "platform" {
  name = "Platform Team"
}
```

### Arguments

| Argument | Required | Description |
|---|---|---|
| `name` | yes | Name of the team. |

### Computed Attributes

| Attribute | Description |
|---|---|
| `id` | Numeric team ID. |
| `created_at` | Creation timestamp. |

### Import

```shell
terraform import tribal_team.platform 1
```

> **Note:** Tribal has no team deletion API. Running `terraform destroy` on a `tribal_team` resource is a no-op.

## Resource: `tribal_admin_settings`

Manages organization-wide notification settings. This is a singleton — only one instance should exist per Terraform workspace.

```hcl
resource "tribal_admin_settings" "org" {
  reminder_days    = [60, 30, 14, 7, 1]
  notify_hour      = 9     # UTC hour
  alert_on_overdue = true
  alert_on_delete  = true

  # Optional
  slack_webhook = "https://hooks.slack.com/services/YOUR/ORG/WEBHOOK"
}
```

### Arguments

| Argument | Required | Description |
|---|---|---|
| `reminder_days` | yes | List of days before expiration to send reminders. |
| `notify_hour` | yes | UTC hour (0–23) at which daily reminders are sent. |
| `alert_on_overdue` | yes | Send alerts for already-expired resources. |
| `alert_on_delete` | yes | Send an admin Slack alert when a resource is deleted. |
| `slack_webhook` | no | Organization-wide Slack webhook URL. |

> **Note:** `terraform destroy` on `tribal_admin_settings` is a no-op (settings cannot be deleted, only updated).

## Development

### Requirements

- Go 1.21+
- Terraform 1.5+

### Build

```shell
go build -o terraform-provider-tribal .
```

### Run Acceptance Tests

Requires a running Tribal instance at `http://localhost:8000`.

```shell
TF_ACC=1 go test ./internal/provider/... -v -timeout 120s
```

## Publishing to the Terraform Registry

To publish this provider to the [Terraform Registry](https://registry.terraform.io):

1. **Create a GitHub repository** named `terraform-provider-tribal` under your GitHub account or organization.

2. **Sign your releases with GPG.** Generate a key pair and add the public key to your Terraform Registry account under *Signing Keys*.

3. **Tag a release** following the `vX.Y.Z` convention:
   ```shell
   git tag v0.1.0
   git push origin v0.1.0
   ```

4. **Set up GoReleaser.** Create `.goreleaser.yml` at the repo root — HashiCorp provides an [official template](https://developer.hashicorp.com/terraform/registry/providers/publishing#github-actions-workflow). Add the `GITHUB_TOKEN` and `GPG_PRIVATE_KEY` secrets to your repository.

5. **Add the GitHub Actions workflow** (`.github/workflows/release.yml`) that runs GoReleaser on tag push. The workflow produces the signed binaries and checksums that the Registry expects.

6. **Connect to the Terraform Registry.** Go to [registry.terraform.io](https://registry.terraform.io), sign in with GitHub, click *Publish > Provider*, and select the repository. The Registry will pick up the release automatically.

Once published, users can reference the provider as:

```hcl
terraform {
  required_providers {
    tribal = {
      source  = "seaburr/tribal"
      version = "~> 0.1"
    }
  }
}
```

See the [HashiCorp provider publishing guide](https://developer.hashicorp.com/terraform/registry/providers/publishing) for full details.
