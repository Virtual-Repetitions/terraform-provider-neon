terraform {
  required_providers {
    neon = {
      source = "example.org/virtual-repetitions/neon"
    }
  }
}

data "neon_projects" "example" {}