---
- hosts: all
  become: yes
  gather_facts: no
  vars:
    ansible_python_interpreter: "/usr/bin/env python3"
  tasks:
  - import_role:
      name: ${server_role}


- hosts: db_servers
  become: yes
  gather_facts: no
  vars:
    ansible_python_interpreter: "/usr/bin/env python3"
  tasks:
  - import_role:
      name: ${db_server_role}
  - import_role:
      name: ${db_client_role}
  - import_role:
      name: ${db_schema_role}


- hosts: ldap_servers
  become: yes
  gather_facts: no
  vars:
    ansible_python_interpreter: "/usr/bin/env python3"
  tasks:
  - import_role:
      name: ${ldap_server_role}
  - import_role:
      name: ${ldap_utils_role}
  - import_role:
      name: ${user_management_role}
  - import_role:
      name: ${user_creation_role}

