#!/bin/bash

set -euo pipefail

_hasher="`dirname ${BASH_SOURCE[0]}`/hash.go"


ADMIN_PASSWD=${ADMIN_PASSWD:-Admin123}

# Create an admin user.
POSTGRES_USER=${POSTGRES_USER:-root}
POSTGRES_HOST=${POSTGRES_HOST:-localhost}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-root}
POSTGRES_DB=${POSTGRES_DB:-template}
POSTGRES_PORT=${POSTGRES_PORT:-5432}

HASH_PEPPER=${HASH_PEPPER:-`openssl rand -base64 64`}

cat <<EOSQL | psql "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB"
INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
VALUES (
	'Admin', 'admin@strv.com', '$(go run ${_hasher} "${ADMIN_PASSWD}")', 'admin', NOW(), NOW()
) ON CONFLICT DO NOTHING;
EOSQL
