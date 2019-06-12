#!/bin/bash
docker-compose -f ./sawtooth-default.yaml down >> /tmp/out.log 2>> /tmp/out.log && docker-compose -f ./sawtooth-default.yaml up --force-recreate >> /tmp/out.log 2>> /tmp/out.log &
