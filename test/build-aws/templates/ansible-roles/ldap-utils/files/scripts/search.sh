#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "$${BASH_SOURCE[0]}" )" && pwd )"
cd $${SCRIPT_DIR}

PWFILE=/etc/ldap/secrets/bind-pw.txt

if [ -r $${PWFILE} ]; then
  ldapsearch -H ldap://${ldap_server} -ZZ -x -D cn=${bind_cn},ou=admin,${domain_dn} -y $${PWFILE} -b ${domain_dn} -LLL
else
  ldapsearch -H ldap://${ldap_server} -ZZ -x -D cn=${bind_cn},ou=admin,${domain_dn} -W -b ${domain_dn} -LLL
fi



