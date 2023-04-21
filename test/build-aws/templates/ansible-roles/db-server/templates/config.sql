SET PASSWORD FOR 'root'@'localhost' = PASSWORD('${root_password}');

{% for server in admin_servers %}
GRANT ALL PRIVILEGES ON *.* TO 'root'@'{{server}}' IDENTIFIED BY '${root_password}' WITH GRANT OPTION;
{% endfor %}
