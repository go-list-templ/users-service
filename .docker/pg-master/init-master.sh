#!/bin/bash

set -e

echo "host replication replica_user 0.0.0.0/0 md5" >> "$PGDATA/pg_hba.conf"

cat >> "$PGDATA/postgresql.conf" << EOF
# Replication settings
wal_level = replica
max_wal_senders = 10
max_replication_slots = 10
synchronous_commit = off

# Performance settings
shared_buffers = 256MB
checkpoint_completion_target = 0.7
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_cache_size = 512MB
EOF

# Create replication user
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER replica_user REPLICATION LOGIN ENCRYPTED PASSWORD 'root';
EOSQL
