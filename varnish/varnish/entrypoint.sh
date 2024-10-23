#!/bin/bash

set -exu

export TINYPROXY_IP="$(dig +short tinyproxy | head -n 1)"
sed "s/__TINYPROXY_IP__/${TINYPROXY_IP}/g" /config/main.vcl.template > /config/main.vcl

varnishd \
    -s file,/storage,1G \
    -a :8080,PROXY \
    -n /var/lib/varnish/varnishd \
    -f /config/main.vcl \
    -F
