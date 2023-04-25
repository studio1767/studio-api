
service:
  server: ${api_server}
  listen_address: ${listen_address}
  listen_port: ${listen_port}
  ca_cert_file: ${ca_cert_file}
  cert_file: ${cert_file}
  key_file: ${key_file}

db:
  server: ${db_server}
  port: ${db_port}
  database: ${db_name}
  user: ${db_user}
  password: "${db_password}"

ldap:
  server_uri: ${ldap_server_uri}
  search_base: ${ldap_search_base}
  bind_dn: ${ldap_bind_dn}
  bind_pw: "${ldap_bind_pw}"
  start_tls: ${ldap_start_tls}

