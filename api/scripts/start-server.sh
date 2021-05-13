#!/bin/bash

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
DB_SCRIPTS_DIR="$DIR/../db/scripts"

if ! psql "$DB_INTEGRATION_CONN_STR" -l &>/dev/null; then
    "$DB_SCRIPTS_DIR/create-test-db.sh" "$DB_INTEGRATION_NAME" "$DB_INTEGRATION_USER" "$DB_INTEGRATION_PASSWORD"
fi

flyway -user="$DB_INTEGRATION_USER" -password="$DB_INTEGRATION_PASSWORD" -url="jdbc:postgresql://localhost/$DB_INTEGRATION_NAME" -locations="filesystem:$DB_SCRIPTS_DIR/../migrations" migrate

export DB_CONN_STR="$DB_INTEGRATION_CONN_STR"

go run $DIR/../main.go
