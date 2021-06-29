terraform {
  required_version = "0.15.3"
  required_providers {
    local = {
      source = "hashicorp/local"
      version = "2.1.0"
    }
    time = {
      source = "hashicorp/time"
      version = "0.7.1"
    }

  }
}
