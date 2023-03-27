
listener: ${listener}

service:
  id: ${client_id}
  secret: "${client_secret}"
  state_secret: "${state_secret}"
  cookie_hash_secret: "${cookie_hash_secret}"
  cookie_enc_secret: "${cookie_enc_secret}"
  redirect_urls:
  %{ for url in redirect_urls ~}
  - ${url}
  %{ endfor ~}

https:
  cert_file: ${https_cert_file}
  key_file: ${https_key_file}

idp:
  issuer_url: ${idp_issuer_url}
  ca_cert_file: ${idp_ca_cert_file}

db:
  server: ${db_server}
  port: ${db_port}
  database: ${db_name}
  user: ${db_user}
  password: "${db_password}"
  ca_cert_file: ${db_ca_cert_file}
