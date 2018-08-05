#!/usr/bin/env bash
file=/var/www/html/funcionarios/deploy/funcionarios
if [ -e "$file" ]; then
    systemctl stop funcionarios.service
    mv "$file" /var/www/html/funcionarios
    chmod 0774 /var/www/html/funcionarios/funcionarios
    systemctl start funcionarios.service
fi