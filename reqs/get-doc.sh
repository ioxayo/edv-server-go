#!/bin/bash

EDV_ID=$1
DOC_ID=$2
curl -v \
-H "Accept: application/json; charset=utf-8" \
-X GET http://localhost:$EDV_PORT/edvs/$EDV_ID/docs/$DOC_ID
