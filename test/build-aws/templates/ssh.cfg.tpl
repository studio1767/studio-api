
# servers
%{ for key, value in servers ~}
Host ${key}
  Hostname ${value}
%{ endfor ~}


Host *
  User ubuntu
  IdentityFile ${ssh_key_file}
  IdentitiesOnly yes

