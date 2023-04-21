## render the run script

resource "local_file" "run_playbook" {
  content = templatefile("templates/ansible/run-ansible.sh.tpl", {
      inventory_file = "inventory.ini"
    })
  filename = "local/ansible/run-ansible.sh"
  file_permission = "0755"
}


## render the playbook

resource "local_file" "playbook" {
  content = templatefile("templates/ansible/playbook.yml.tpl", {
      server_role = local.server_role,
      db_server_role = local.db_server_role,
      db_client_role = local.db_client_role,
      db_schema_role = local.db_schema_role,
      ldap_server_role = local.ldap_server_role,
      ldap_utils_role = local.ldap_utils_role,
      user_management_role = local.user_management_role,
      user_creation_role = local.user_creation_role,
    })
  filename = "local/ansible/playbook.yml"
  file_permission = "0640"
}


## render the inventory file

resource "local_file" "inventory" {
  content = templatefile("templates/ansible/inventory.ini.tpl", {
    ldap_servers = [ var.services.ldap.host_names[0] ],
    db_servers = [ var.services.db.host_names[0] ],
  })
  filename = "local/ansible/inventory.ini"
  file_permission = "0640"
}

