---
page_title: "Provider: Tribal"
description: |-
  The Tribal provider manages resources in Tribal, a credential and resource expiry tracking tool.
---

# Tribal Provider

The Tribal provider manages credentials, certificates, teams, and organization settings in [Tribal](https://github.com/seaburr/tribal), a tool for tracking expiring resources and sending renewal notifications.

## Example Usage

```terraform
terraform {
  required_providers {
    tribal = {
      source  = "seaburr/tribal"
      version = "~> 0.1"
    }
  }
}

provider "tribal" {
  host    = "http://localhost:8000"
  api_key = "tribal_sk_..."
}
```

## Authentication

The provider requires an API key for authentication. It can be provided via:

- The `api_key` argument in the provider block
- The `TRIBAL_API_KEY` environment variable

## Argument Reference

- `host` (Optional) - Base URL of the Tribal API. Defaults to `http://localhost:8000`. Can also be set via the `TRIBAL_HOST` environment variable.
- `api_key` (Required) - API key for authenticating with the Tribal API. Can also be set via the `TRIBAL_API_KEY` environment variable.
