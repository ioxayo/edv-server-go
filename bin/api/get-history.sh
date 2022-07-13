#!/usr/bin/env bash

EDV_ID=$1

for i in "$@"; do
  case $i in
    -a=*|--after=*) AFTER_SEQ="${i#*=}";;
    -b=*|--before=*) BEFORE_SEQ="${i#*=}";;
    -*)
      echo "Invalid option: $i"
      exit 1
      ;;
    *)
      ;;
  esac
done

curl -v \
-H "Accept: application/json; charset=utf-8" \
-X GET http://localhost:5000/edvs/$EDV_ID/history?afterSequence=${AFTER_SEQ}\&beforeSequence=${BEFORE_SEQ}
