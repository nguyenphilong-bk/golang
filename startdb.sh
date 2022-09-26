#!/bin/bash

docker-compose up -d --force-recreate

sleep 5

docker exec mongodb /scripts/rs-init.sh