---
server_name: ${server_name}

ldap_groups:
%{ for group in groups ~}
- "${group}"
%{ endfor ~}

ldap_users:
%{ for k, user in users ~}
- uname: "${k}"
  gname: "${user.given_name}"
  fname: "${user.family_name}"
  groups:
%{ for group in user.groups ~}
  - "${group}"
%{ endfor ~}
%{ endfor ~}

ldap_user_groups:
%{ for k, user in users ~}
%{ for group in user.groups ~}
- user: ${k}
  group: ${group}
%{ endfor ~}
%{ endfor ~}
