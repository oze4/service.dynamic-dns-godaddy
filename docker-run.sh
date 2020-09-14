#!/bin/bash

docker run                   \
  -d                         \
  -e GODADDY_APIKEY='-'      \
  -e GODADDY_APISECRET='-'   \
  -e GODADDY_DOMAIN='-'      \
  -e BASELINE_RECORD='-'     \
  --name godaddy_dynamic_dns \
  oze4/godaddy-dynamic-dns:latest