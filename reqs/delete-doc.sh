#!/bin/bash

EDV_ID=$1
DOC_ID=$2
curl -v \
-X DELETE http://localhost:$EDV_PORT/edvs/$EDV_ID/docs/$DOC_ID \
