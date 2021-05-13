#!/bin/bash

set -euo pipefail

dbname=${1:?db name missing}
uid=${2:?user name missing}
passwd=${3:?password missing}

sudo -u postgres createdb "$dbname"
sudo -u postgres createuser "$uid"

sudo -u postgres psql <<EOF
alter user $uid with encrypted password '$passwd';
grant all privileges on database '$dbname' to '$uid';
EOF
