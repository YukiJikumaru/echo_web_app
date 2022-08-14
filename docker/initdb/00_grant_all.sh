#!/bin/bash

echo "##### START 00_grant_all.sh #####"

DATABASES='echo_web_app_development echo_web_app_test'
DB_USERS='application'
DEFAULT_PASSWORD='password'

for DB_NAME in $DATABASES; do
    echo "CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\`;"
    echo "CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\`;" | mysql -uroot -p${MYSQL_ROOT_PASSWORD}
done

for DB_USER in $DB_USERS; do
    echo "CREATE USER '${DB_USER}'@'%' IDENTIFIED BY '${DEFAULT_PASSWORD}';"
    echo "CREATE USER '${DB_USER}'@'%' IDENTIFIED BY '${DEFAULT_PASSWORD}';" | mysql -uroot -p${MYSQL_ROOT_PASSWORD}
done

for DB_NAME in $DATABASES; do
    for DB_USER in $DB_USERS; do
        echo "GRANT ALL ON \`${DB_NAME}\`.* TO '${DB_USER}'@'%';"
        echo "GRANT ALL ON \`${DB_NAME}\`.* TO '${DB_USER}'@'%';" | mysql -uroot -p${MYSQL_ROOT_PASSWORD}
    done
done

echo 'FLUSH PRIVILEGES;'
echo 'FLUSH PRIVILEGES;' | "${mysql[@]}"

echo "##### END 00_grant_all.sh #####"
