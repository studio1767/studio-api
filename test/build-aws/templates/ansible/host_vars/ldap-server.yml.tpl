---
server_name: ${server_name}

users:
%{ for k, user in users ~}
- { uname: "${k}", gname: "${user.given_name}", fname: "${user.family_name}"}
%{ endfor ~}
