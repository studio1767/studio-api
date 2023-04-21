## Bad Certificate Authority

resource "tls_private_key" "bad_ca" {
  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_self_signed_cert" "bad_ca" {
  private_key_pem = tls_private_key.bad_ca.private_key_pem

  subject {
    common_name  = "${var.studio_name} Bad CA"
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

resource "local_file" "bad_ca" {
  content         = tls_self_signed_cert.bad_ca.cert_pem
  filename        = "local/configs/certs/bad-ca.crt"
  file_permission = "0644"
}

## Bad User

resource "tls_private_key" "bad_user" {
  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_cert_request" "bad_user" {
  private_key_pem = tls_private_key.bad_user.private_key_pem

  subject {
    common_name  = "baduser"
    organization = var.studio_name
    country = "AU"
  }

  uris = ["email:baduser@example.xyx"]
}

resource "tls_locally_signed_cert" "bad_user" {
  ca_private_key_pem = tls_private_key.bad_ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.bad_ca.cert_pem
  cert_request_pem   = tls_cert_request.bad_user.cert_request_pem

  validity_period_hours = 240
  early_renewal_hours = 48

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

resource "local_file" "bad_user_key" {
  content         = tls_private_key.bad_user.private_key_pem
  filename        = "local/configs/certs/users/bad-user.key"
  file_permission = "0600"
}

resource "local_file" "bad_user_cert" {
  content         = tls_locally_signed_cert.bad_user.cert_pem
  filename        = "local/configs/certs/users/bad-user.crt"
  file_permission = "0644"
}

