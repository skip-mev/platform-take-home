terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
      version = "3.24.0"
    }
  }
}

provider "vault" {
  address = "http://localhost:8200"
}

resource "vault_policy" "app_role_policy" {
  name   = "app-role-policy"
  policy = <<EOT
path "transit/encrypt/*" {
  capabilities = ["update"]
}
path "transit/decrypt/*" {
  capabilities = ["update"]
}
path "transit/sign/*" {
  capabilities = ["update"]
}
path "transit/keys/*" {
  capabilities = ["create", "read", "update", "list"]
}
path "transit/export/*" {
  capabilities = ["deny"]
}
path "transit/delete/*" {
  capabilities = ["deny"]
}
EOT
}

resource "vault_auth_backend" "approle" {
  type = "approle"
}

resource "vault_approle_auth_backend_role" "app_role" {
  backend        = vault_auth_backend.approle.path
  role_name      = "app_role"
  token_policies = [vault_policy.app_role_policy.name]
  
  token_ttl     = 3600   # 1 hour in seconds
  token_max_ttl = 14400  # 4 hours in seconds
}

resource "vault_mount" "transit" {
  path        = "transit"
  type        = "transit"
  description = "Transit Secrets Engine for signing service"
}
