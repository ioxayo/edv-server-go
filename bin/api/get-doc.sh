#!/usr/bin/env bash

EDV_ID=$1
DOC_ID=$2

curl -v \
-H "Accept: application/json; charset=utf-8" \
-X GET http://localhost:5000/edvs/$EDV_ID/docs/$DOC_ID
