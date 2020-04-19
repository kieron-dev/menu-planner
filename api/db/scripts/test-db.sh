#!/bin/bash

set -u

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

name=$(uuidgen)
name=u${name//-/}
passwd=$(uuidgen)
passwd=${passwd//-/}

"$DIR/create-test-db.sh" "$name" "$passwd"

flyway -user="$name" -password="$passwd" -url="jdbc:postgresql://localhost/$name" -locations="filesystem:$DIR/../migrations" migrate

export DB_CONN_STR="host=localhost dbname=$name user=$name password=$passwd"

ginkgo db

"$DIR/delete-test-db.sh" "$name"
