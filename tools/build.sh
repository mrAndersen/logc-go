#!/usr/bin/env bash

docker build -f ./docker/Dockerfile -t logcgo ./
docker rm --force loggo || true; docker create --name loggo logcgo