## Certificate Authority

resource "tls_private_key" "ca" {
  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_self_signed_cert" "ca" {
  private_key_pem = tls_private_key.ca.private_key_pem

  subject {
    common_name  = "${var.studio_name} CA"
    organization = var.studio_name
    country = "AU"
  }
  
  validity_period_hours = 480
  early_renewal_hours = 48

  allowed_uses = [
    "cert_signing",
    "crl_signing"
  ]
  is_ca_certificate = true
}


## Users

resource "tls_private_key" "users" {
  for_each = var.users

  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_cert_request" "users" {
  for_each = var.users

  private_key_pem = tls_private_key.users[each.key].private_key_pem

  subject {
    common_name  = "${each.key}@${var.studio_domain}"
    organization = var.studio_name
    country = "AU"
  }

  uris = [for group in each.value.cert_groups: "group:${group}"]
}

resource "tls_locally_signed_cert" "users" {
  for_each = tls_cert_request.users

  ca_private_key_pem = tls_private_key.ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem
  cert_request_pem   = each.value.cert_request_pem

  validity_period_hours = 240
  early_renewal_hours = 48

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}


## Services

resource "tls_private_key" "services" {
  for_each = var.services
  
  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_cert_request" "services" {
  for_each = var.services

  private_key_pem = tls_private_key.services[each.key].private_key_pem

  subject {
    common_name  = "${each.key}.${var.studio_domain}"
    organization = var.studio_name
    country = "AU"
  }

  dns_names = concat(each.value.host_names, [ for host in each.value.host_names: "${host}.${var.studio_domain}" ])
  ip_addresses = local.services_all_ip[each.key]
  uris = each.value.uris
}

resource "tls_locally_signed_cert" "services" {
  for_each = tls_cert_request.services
  
  ca_private_key_pem = tls_private_key.ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem
  cert_request_pem   = each.value.cert_request_pem

  validity_period_hours = 240
  early_renewal_hours = 48

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}


resource "local_file" "services" {
  for_each = tls_locally_signed_cert.services
  
  content         = each.value.cert_pem
  filename        = "local/configs/certs/${each.key}.crt"
  file_permission = "0644"
}

