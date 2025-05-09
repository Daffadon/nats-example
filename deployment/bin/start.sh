#!/bin/bash

docker compose -p factory\
  -f deployment/docker/docker-compose.nats.yml \
  -f deployment/docker/docker-compose.db.yml \
  up -d
