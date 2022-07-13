#!/usr/bin/env bash

EDV_ID=$1
VERBOSE=false

for i in "$@"; do
  case $i in
    -v|--verbose) VERBOSE=true;;
    -*)
      echo "Invalid option: $i"
      exit 1
      ;;
    *)
      ;;
  esac
done

curl -v \
-H "Content-Type: application/json; charset=utf-8" \
-X POST http://localhost:5000/edvs/$EDV_ID/query \
-d @- << DATA
{
  "index": "https://example.com/kms/z7BgF536GaR",
  "equals": [
    {
      "DUQaxPtSLtd8L3WBAIkJ4DiVJeqoF6bdnhR7lSaPloZ": "RV58Va4904K-18_L5g_vfARXRWEB00knFSGPpukUBro"
    },
    {
      "AarngVIZLl0kIp2xEHUH5o5uVc-470roQaOIbqMUD7DFQQypWQ==": "AYubg9VnEitQBxlhjVFnYRlfQ5UHWe3ia4aMiQ6srhcrXtEK2Q=="
    }
  ],
  "returnFullDocuments": $VERBOSE
}
DATA
