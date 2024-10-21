#!/bin/bash

set -exu

WORK_DIR="$(pwd)"
DATA_DIR="${DATA_DIR:=${WORK_DIR}/data}"
STORAGE_DIR="${STORAGE_DIR:=${WORK_DIR}/storage}"

mkdir -p "$DATA_DIR"
mkdir -p "$STORAGE_DIR"

varnishd \
    -n "$DATA_DIR" \
    -f "${WORK_DIR}/config/main.vcl" \
    -s "file,${STORAGE_DIR},1G" \
    -T 127.0.0.1:2000 \
    -a 0.0.0.0:8080 \
    -F


#FIXME: start tinyproxy at the same time
# tinyproxy -d
