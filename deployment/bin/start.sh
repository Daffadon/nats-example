#!/bin/bash

docker compose -p iot \
  -f deployment/docker/docker-compose.nats.yml \
  -f deployment/docker/docker-compose.db.yml \
  -f deployment/docker/docker-compose.monitoring.yml \
  up -d --wait

nsc push -A -u nats://localhost:4221
