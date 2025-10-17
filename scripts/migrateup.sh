#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd internal/sql/schema
goose turso $DB_URL up
