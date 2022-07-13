#!/usr/bin/env bash

EDV_ID=$1
DOC_ID=$2

curl -v \
-X DELETE http://localhost:5000/edvs/$EDV_ID/docs/$DOC_ID \
