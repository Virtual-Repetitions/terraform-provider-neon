terraform {
  required_providers {
    neon = {
      source = "example.org/virtual-repetitions/neon"
    }
  }
}

provider "neon" {}

resource "neon_project" "example" {
  name            = "example-project-with-branches"
  instance_handle = "scalable"
  platform_id     = "aws"
  region_id       = "aws-us-west-2"
}

resource "neon_branch" "example_branch" {
  parent_project_id = neon_project.example.id
}

