#!/bin/bash

uuid=${1:?db name missing}

sudo -u postgres psql <<EOF
drop database $uuid;
drop user $uuid;
EOF
