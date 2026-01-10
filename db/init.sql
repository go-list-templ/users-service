CREATE DATABASE papa
    WITH
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';

\c papa;

create table outbox
(
    message_id uuid primary key,
    message    jsonb
);

alter table outbox
    replica identity full;

create publication sequin_pub for table outbox;

select pg_create_logical_replication_slot('sequin_slot', 'pgoutput');