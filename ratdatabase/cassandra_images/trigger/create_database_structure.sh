#!/bin/bash
(sleep 45 && cqlsh -f create_trigger_structure.cql) &
cd / && docker-entrypoint.sh