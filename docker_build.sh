#! /bin/bash

docker buildx build --platform linux/amd64,linux/arm/v7,linux/arm64 -t adriencos/pidarr_bot:latest --push .
