CREATE DATABASE IF NOT EXISTS ${db_name};
USE ${db_name};

CREATE USER IF NOT EXISTS ${db_user}@'${db_client}' IDENTIFIED BY '${db_password}';
GRANT ALL PRIVILEGES ON ${db_name}.* TO ${db_user}@'${db_client}';

CREATE TABLE IF NOT EXISTS project (
  id         INT UNSIGNED AUTO_INCREMENT NOT NULL,
  name       VARCHAR(256) NOT NULL,
  code       VARCHAR(64) NOT NULL,
  PRIMARY KEY (`id`)
);

