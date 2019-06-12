#!/bin/bash


echo tcp://$(docker inspect -f "{{ .NetworkSettings.Networks.samplenetwork.IPAddress }}" sawtooth-validator-default):4004