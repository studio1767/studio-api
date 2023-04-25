#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "$${BASH_SOURCE[0]}" )" && pwd )"
cd $${SCRIPT_DIR}

PWFILE=/etc/ldap/secrets/bind-pw.txt

if [ -r $${PWFILE} ]; then
  ldapwhoami -H ldap://${ldap_server} -ZZ -x -D cn=${bind_dn} -y $${PWFILE}
else
  ldapwhoami -H ldap://${ldap_server} -ZZ -x -D cn=${bind_dn} -W
fi

