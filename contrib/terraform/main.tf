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

# Put your Terraform code here