# SPDX-FileCopyrightText: © 2022-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: GPL-3.0-or-later

terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.55.0"
    }
  }
  required_version = ">= 0.14"
}

// leave token commented out to pull ENV value from BASH
provider "digitalocean" {
  // token = var.do_token
}

data "digitalocean_ssh_key" "terraform" {
  name = "franklin"
}
