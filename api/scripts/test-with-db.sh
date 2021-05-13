#!/bin/bash

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
DB_SCRIPTS_DIR="$DIR/../db/scripts"

if ! psql "$DB_TEST_CONN_STR" -l &>/dev/null; then
    "$DB_SCRIPTS_DIR/create-test-db.sh" "$DB_TEST_NAME" "$DB_TEST_USER" "$DB_TEST_PASSWORD"
fi

flyway -user="$DB_TEST_USER" -password="$DB_TEST_PASSWORD" -url="jdbc:postgresql://localhost/$DB_TEST_NAME" -locations="filesystem:$DB_SCRIPTS_DIR/../migrations" migrate

export DB_CONN_STR="$DB_TEST_CONN_STR"

ginkgo -p -r $@
