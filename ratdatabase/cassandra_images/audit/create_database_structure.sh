#!/bin/bash
(sleep 45 && cqlsh -f create_audit_structure.cql) &
cd / && docker-entrypoint.sh