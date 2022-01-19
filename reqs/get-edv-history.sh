#!/bin/bash

EDV_ID=$1
curl -v \
-H "Accept: application/json; charset=utf-8" \
-X GET http://localhost:$EDV_PORT/edvs/$EDV_ID/history \
