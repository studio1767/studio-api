---
server_name: ${server_name}

admin_servers:
%{ for address in admin_servers ~}
- ${address}
%{ endfor ~}
