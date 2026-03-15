---
page_title: "tribal_team (Resource) - tribal"
description: |-
  Manages a team in Tribal.
---

# tribal_team

Manages a team in Tribal. Teams are used to group tracked resources under a named organizational unit.

-> **Note:** Tribal has no team deletion API. Running `terraform destroy` on a `tribal_team` resource is a no-op; the team will remain in Tribal.

## Example Usage

```terraform
resource "tribal_team" "platform" {
  name = "Platform Team"
}

resource "tribal_resource" "prod_api_key" {
  name            = "Production API Key"
  dri             = "platform-team@example.com"
  type            = "API Key"
  expiration_date = "2026-12-31"
  purpose         = "Authenticates the payment service with Stripe"
  generation_instructions = "Log into Stripe dashboard > Developers > API keys > Create new key"
  slack_webhook   = "https://hooks.slack.com/services/TEAM/CHANNEL/WEBHOOK"
  team_id         = tribal_team.platform.id
}
```

## Argument Reference

### Required

- `name` - Name of the team.

## Attribute Reference

In addition to the arguments above, the following computed attributes are exported:

- `id` - Numeric team ID.
- `created_at` - Timestamp when the team was created.

## Import

`tribal_team` resources can be imported using the numeric team ID:

```shell
terraform import tribal_team.platform 1
```
