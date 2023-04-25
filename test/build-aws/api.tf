
## create the api server config file

resource "local_file" "api_ca_certificate" {
  content         = tls_self_signed_cert.ca.cert_pem
  filename        = "local/configs/certs/ca.crt"
  file_permission = "0644"
}

resource "local_file" "api_key" {
  content         = tls_private_key.services["api"].private_key_pem
  filename        = "local/configs/certs/api.key"
  file_permission = "0600"
}

resource "local_file" "api_cert" {
  content         = tls_locally_signed_cert.services["api"].cert_pem
  filename        = "local/configs/certs/api.crt"
  file_permission = "0644"
}


resource "local_file" "api_server" {
  content = templatefile("../../configs/config.yaml.tpl", {
    api_server       = local.services_all_ip["api"][0]
    listen_address   = var.services.api.listen_address
    listen_port      = var.services.api.listen_port
    ca_cert_file     = "certs/ca.crt"
    cert_file        = "certs/api.crt"
    key_file         = "certs/api.key"
    db_server        = aws_instance.server["db"].public_ip
    db_port          = var.services.db.listen_port
    db_name          = var.studio_code
    db_user          = var.studio_code
    db_password      = random_password.db_user.result
    ldap_server_uri  = "ldaps://${aws_instance.server["ldap"].public_ip}:${var.services.ldap.listen_port}"
    ldap_search_base = local.studio_domain_dn
    ldap_bind_dn     = "cn=admin,${local.studio_domain_dn}"
    ldap_bind_pw     = random_password.ldap_admin.result
    ldap_start_tls   = false
  })
  filename        = "local/configs/server.yaml"
  file_permission = "0640"
}


## create the api client config file

resource "local_file" "client_user_key" {
  for_each = tls_private_key.users

  content         = each.value.private_key_pem
  filename        = "local/configs/certs/users/${each.key}.key"
  file_permission = "0600"
}

resource "local_file" "client_user_cert" {
  for_each = tls_locally_signed_cert.users

  content         = each.value.cert_pem
  filename        = "local/configs/certs/users/${each.key}.crt"
  file_permission = "0644"
}

resource "null_resource" "user_pkcs12" {
  for_each = var.users

  triggers = {
    certificate = tls_locally_signed_cert.users[each.key].cert_pem
  }

  provisioner "local-exec" {
    command = <<-CMD
      openssl pkcs12 -export -inkey ${basename(local_file.client_user_key[each.key].filename)}   \
                      -in ${basename(local_file.client_user_cert[each.key].filename)}            \
                      -name ${var.users[each.key].given_name}_${var.users[each.key].family_name} \
                      -out ${each.key}.pfx   \
                      -passout pass:
      CMD

    working_dir = "local/configs/certs/users"
  }
}


resource "local_file" "api_client" {
  content = templatefile("../client/config.yaml.tpl", {
    api_server      = local.services_all_ip["api"][0]
    api_port        = var.services.api.listen_port
    ca_cert_file    = "certs/ca.crt"
    admin_cert_file = "certs/users/admin.crt"
    admin_key_file  = "certs/users/admin.key"
    operator_cert_file = "certs/users/user1.crt"
    operator_key_file  = "certs/users/user1.key"
    user_cert_file  = "certs/users/user2.crt"
    user_key_file   = "certs/users/user2.key"
  })
  filename        = "local/configs/client.yaml"
  file_permission = "0640"
}

resource "local_file" "bad_api_client1" {
  content = templatefile("../client/config.yaml.tpl", {
    api_server      = local.services_all_ip["api"][0]
    api_port        = var.services.api.listen_port
    ca_cert_file    = "certs/ca.crt"
    admin_cert_file = "certs/users/admin.crt"
    admin_key_file  = "certs/users/admin.key"
    operator_cert_file = "certs/users/user1.crt"
    operator_key_file  = "certs/users/user1.key"
    user_cert_file  = "certs/users/bad-user.crt"
    user_key_file   = "certs/users/bad-user.key"
  })
  filename        = "local/configs/bad-client-1.yaml"
  file_permission = "0640"
}

resource "local_file" "bad_api_client2" {
  content = templatefile("../client/config.yaml.tpl", {
    api_server      = local.services_all_ip["api"][0]
    api_port        = var.services.api.listen_port
    ca_cert_file    = "certs/bad-ca.crt"
    admin_cert_file = "certs/users/admin.crt"
    admin_key_file  = "certs/users/admin.key"
    operator_cert_file = "certs/users/user1.crt"
    operator_key_file  = "certs/users/user1.key"
    user_cert_file  = "certs/users/user2.crt"
    user_key_file   = "certs/users/user2.key"
  })
  filename        = "local/configs/bad-client-2.yaml"
  file_permission = "0640"
}

resource "local_file" "bad_api_client3" {
  content = templatefile("../client/config.yaml.tpl", {
    api_server      = local.services_all_ip["api"][0]
    api_port        = var.services.api.listen_port
    ca_cert_file    = "certs/bad-ca.crt"
    admin_cert_file = "certs/users/admin.crt"
    admin_key_file  = "certs/users/admin.key"
    operator_cert_file = "certs/users/user1.crt"
    operator_key_file  = "certs/users/user1.key"
    user_cert_file  = "certs/users/bad-user.crt"
    user_key_file   = "certs/users/bad-user.key"
  })
  filename        = "local/configs/bad-client-3.yaml"
  file_permission = "0640"
}
