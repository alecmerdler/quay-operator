#!/bin/sh

echo "attempting to install required pg_trgm extension"

psql -d $POSTGRESQL_DATABASE -c 'CREATE EXTENSION pg_trgm;'
