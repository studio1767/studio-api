## role names

locals {
  server_role = "server"
  db_server_role = "db-server"
  db_client_role = "db-client"
  db_schema_role = "db-schema"
  ldap_server_role = "ldap-server"
  ldap_utils_role = "ldap-utils"
  user_management_role = "user-management"
  user_creation_role = "user-creation"
}

## create the servers

data "aws_ami" "ami" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "server" {
  for_each = toset([ for k, v in var.services: k if k != "api" ])
  
  instance_type = var.aws_instance_type
  ami           = data.aws_ami.ami.id

  disable_api_termination     = false
  associate_public_ip_address = true
  source_dest_check           = false

  subnet_id                   = aws_subnet.subnet.id
  key_name                    = var.studio_code
  vpc_security_group_ids      = [aws_security_group.subnet.id]

  tags = {
    Name = "${var.studio_code}-${each.key}"
  }
  
  depends_on = [
    aws_key_pair.ssh_key
  ]
  
  user_data = <<-EOF
  #!/usr/bin/env bash
  hostnamectl set-hostname ${var.studio_code}-${each.key}
  EOF
}


## the server role

resource "template_dir" "server" {
  source_dir      = "templates/ansible-roles/${local.server_role}"
  destination_dir = "local/ansible/roles/${local.server_role}"

  vars = {}
}


## db server

resource "local_file" "db_host_vars" {
  content = templatefile("templates/ansible/host_vars/db-server.yml.tpl", {
    server_name = var.services.db.host_names[0]
    admin_servers = [ local.management_ip ]
  })
  filename        = "local/ansible/host_vars/${var.services.db.host_names[0]}.yml"
  file_permission = "0640"
}

resource "random_password" "db_root" {
  length           = 16
  special          = true
}

resource "random_password" "db_user" {
  length           = 16
  special          = true
}

resource "template_dir" "db_server" {
  source_dir      = "templates/ansible-roles/${local.db_server_role}"
  destination_dir = "local/ansible/roles/${local.db_server_role}"

  vars = {
    ca_cert       = tls_self_signed_cert.ca.cert_pem
    server_cert = tls_locally_signed_cert.services["db"].cert_pem
    server_key = tls_private_key.services["db"].private_key_pem
    root_password = random_password.db_root.result
  }
}

resource "template_dir" "db_client" {
  source_dir      = "templates/ansible-roles/${local.db_client_role}"
  destination_dir = "local/ansible/roles/${local.db_client_role}"

  vars = {
    ca_cert       = tls_self_signed_cert.ca.cert_pem
    root_password = random_password.db_root.result
  }
}

resource "template_dir" "db_schema" {
  source_dir      = "templates/ansible-roles/${local.db_schema_role}"
  destination_dir = "local/ansible/roles/${local.db_schema_role}"

  vars = {
    db_server = "127.0.0.1"
    db_admin = "root"
    db_admin_password = random_password.db_root.result
    db_client = local.management_ip
    db_name = var.studio_code
    db_user = var.studio_code
    db_password = random_password.db_user.result
  }
}


## ldap server

resource "local_file" "ldap_host_vars" {
  content = templatefile("templates/ansible/host_vars/ldap-server.yml.tpl", {
    server_name = var.services.ldap.host_names[0]
    users       = var.users
    groups      = var.groups
  })
  filename        = "local/ansible/host_vars/${var.services.ldap.host_names[0]}.yml"
  file_permission = "0640"
}

resource "random_password" "ldap_admin" {
  length           = 16
  special          = true
}

resource "random_password" "ldap_search" {
  length           = 16
  special          = true
}

resource "template_dir" "ldap_server" {
  source_dir      = "templates/ansible-roles/${local.ldap_server_role}"
  destination_dir = "local/ansible/roles/${local.ldap_server_role}"

  vars = {
    admin_password = random_password.ldap_admin.result
    bind_cn = "search"
    bind_password = random_password.ldap_search.result
    ca_cert = tls_self_signed_cert.ca.cert_pem
    server_ca_cert = tls_self_signed_cert.ca.cert_pem
    server_cert = tls_locally_signed_cert.services["ldap"].cert_pem
    server_key = tls_private_key.services["ldap"].private_key_pem
    organization = lower(var.studio_name)
    domain_dns = var.studio_domain
    domain_dn = local.studio_domain_dn
    domain_dc = split(".", var.studio_domain)[0]
  }
}

resource "template_dir" "ldap_utils" {
  source_dir      = "templates/ansible-roles/${local.ldap_utils_role}"
  destination_dir = "local/ansible/roles/${local.ldap_utils_role}"

  vars = {
    ldap_server = "127.0.0.1"
    ca_cert     = tls_self_signed_cert.ca.cert_pem
    domain_dns  = var.studio_domain
    domain_dn   = local.studio_domain_dn
    bind_dn = "cn=search,ou=admin,${local.studio_domain_dn}"
    bind_password = random_password.ldap_search.result
  }
}

resource "template_dir" "admin_user_management" {
  source_dir      = "templates/ansible-roles/${local.user_management_role}"
  destination_dir = "local/ansible/roles/${local.user_management_role}"

  vars = {
    ldap_uri = "ldap://127.0.0.1:389"
    bind_dn = "cn=admin,${local.studio_domain_dn}"
    bind_pw = random_password.ldap_admin.result
    root_dn = local.studio_domain_dn
    start_tls = true
    ca_cert        = tls_self_signed_cert.ca.cert_pem
  }
}

resource "template_dir" "admin_user_creation" {
  source_dir      = "templates/ansible-roles/${local.user_creation_role}"
  destination_dir = "local/ansible/roles/${local.user_creation_role}"

  vars = {}
}

