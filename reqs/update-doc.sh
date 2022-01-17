#!/bin/bash

EDV_ID=$1
DOC_ID=$2
curl -v \
-H "Accept: application/json; charset=utf-8" \
-X POST http://localhost:$EDV_PORT/edvs/$EDV_ID/docs/$DOC_ID \
-d @- << DATA
{
  "id": "94684128-c42c-4b28-adb0-aec77bf76044",
  "sequence": 1,
  "jwe": {
    "protected": "eyJlbmMiOiJDMjBQIn0",
    "recipients": [{
      "header": {
        "alg": "A256KW",
        "kid": "https://example.com/kms/zSDn2MzzbxmX"
      },
      "encrypted_key": "OR1vdCNvf_B68mfUxFQVT-vyXVrBembuiM40mAAjDC1-Qu5iArDbug"
    }],
    "iv": "i8Nins2vTI3PlrYW",
    "ciphertext": "Cb-963UCXblINT8F6MDHzMJN9EAhK3I",
    "tag": "pfZO0JulJcrc3trOZy8rjA"
  }
}
DATA
