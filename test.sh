#!/bin/sh

set -x

rm -f db.sqlite3

sqlite3 db.sqlite3 <schema.sql
sqlite3 db.sqlite3 <populate.sql
sqlite3 db.sqlite3 <test_queries.sql

