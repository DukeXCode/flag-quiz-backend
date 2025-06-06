#!/bin/sh
echo "Running SQLite migrations..."
sqlite3 /app/data/database.db < /app/migration/schema.sql
if [ $? -ne 0 ]; then
    echo "Schema migration failed!"
    exit 1
fi

exec ./main
