#!/bin/bash

set -euo pipefail

uuid=${1:?db name missing}
passwd=${2:?password missing}

sudo -u postgres createuser "$uuid"
sudo -u postgres createdb "$uuid"

sudo -u postgres psql <<EOF
alter user $uuid with encrypted password '$passwd';
grant all privileges on database $uuid to $uuid;
EOF
