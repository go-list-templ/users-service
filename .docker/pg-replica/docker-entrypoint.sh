#!/bin/bash

set -e

# Ensure proper ownership
chown -R postgres:postgres "$PGDATA"
chmod 700 "$PGDATA"

# Wait until master is ready
echo "Waiting for master to be ready..."
until pg_isready -h master -p 5432 -U replica_user; do
    echo "Still waiting for master..."
    sleep 2
done

# Initialize replica if needed
if [ ! -s "$PGDATA/PG_VERSION" ]; then
    echo "Running base backup..."
    gosu postgres pg_basebackup -h master -U replica_user -D "$PGDATA" -Fp -Xs -P -R
    chmod 700 "$PGDATA"
fi

# Start postgres as the correct user
exec gosu postgres postgres
