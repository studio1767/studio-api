[ldap_servers]
%{ for name in ldap_servers ~}
${name}
%{ endfor ~}

[db_servers]
%{ for name in db_servers ~}
${name}
%{ endfor ~}
