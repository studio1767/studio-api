## Studio

variable "studio_name" {
  default = "Studio1767"
}

variable "studio_code" {
  default = "s1767"
}

variable "studio_domain" {
  default = "example.xyz"
}

locals {
  studio_domain_dn = join(",", formatlist("dc=%s", split(".", var.studio_domain)))
}


## Services

variable "services" {
  type = map(object({
    host_names = list(string)
    ip_addresses = list(string)
    uris = list(string)
    listen_address = string
    listen_port = number
  }))
  default = {
    ldap = {
      host_names = [
        "ldap",
        "ldap-00",
        "ldap-01"
      ]
      ip_addresses = [
        "127.0.0.1"
      ]
      uris = []
      listen_address = "0.0.0.0"
      listen_port = 636
    },
    db = {
      host_names = [
        "db",
        "db-00",
        "db-01"
      ]
      ip_addresses = [
        "127.0.0.1"
      ]
      uris = []
      listen_address = "0.0.0.0"
      listen_port = 3306
    },
    api = {
      host_names = [
        "api",
        "api-00",
        "api-01"
      ]
      ip_addresses = [
        "127.0.0.1"
      ]
      uris = []
      listen_address = "127.0.0.1"
      listen_port = 8443
    }
  }
}

locals {
  services_public_ip = { for k, v in aws_instance.server: k => v.public_ip }
  
  services_all_ip = { for k, v in var.services: k => [
    for ip in concat(v.ip_addresses, [lookup(local.services_public_ip, k, "none")]) : ip if ip != "none"
  ]}
}


## Users and Groups

variable "groups" {
  type = list(string)
  default = [
    "users", "operators", "admins"
  ]
}

variable "users" {
  type = map(object({
    given_name = string
    family_name = string
    groups = list(string)
    cert_groups = list(string)
  }))
  default = {
    admin = {
      given_name = "admin"
      family_name = "admin"
      groups = ["users", "admins"]
      cert_groups = ["users"]
    },
    user1 = {
      given_name = "user"
      family_name = "one"
      groups = ["users", "operators"]
      cert_groups = ["operators"]
    },
    user2 = {
      given_name = "user"
      family_name = "two"
      groups = ["users"]
      cert_groups = ["users"]
    }
  }
}


## Network

data "external" "my_public_ip" {
  program = ["scripts/my-public-ip.sh"]
}

locals {
  management_ip  = data.external.my_public_ip.result["my_public_ip"]
  management_net = "${local.management_ip}/32"
}

variable "cidr_block" {
  type = string
  default = "10.16.0.0/24"
}

## AWS settings

variable "aws_profile" {
  default = ""
}

variable "aws_region" {
  default = ""
}

variable "aws_availability_zone" {
  default = ""
}

data "aws_availability_zones" "site" {
  state = "available"
}

locals {
  num_zones = length(data.aws_availability_zones.site.names)
  aws_availability_zone = var.aws_availability_zone == "" ? element(data.aws_availability_zones.site.names, local.num_zones - 1) : var.aws_availability_zone
}

variable "aws_instance_type" {
  default = "t3a.micro"
}

