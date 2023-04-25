#!/usr/bin/env bash

/usr/bin/mariadb -h 127.0.0.1 -u root -p$( cat /etc/mysql/secrets/password.txt )

