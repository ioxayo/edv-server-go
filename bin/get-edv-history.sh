#!/bin/bash

EDV_ID=$1
AFTER_SEQ=$2
BEFORE_SEQ=$3
curl -v \
-H "Accept: application/json; charset=utf-8" \
-X GET http://localhost:$EDV_PORT/edvs/$EDV_ID/history?afterSequence=$AFTER_SEQ\&beforeSequence=$BEFORE_SEQ
