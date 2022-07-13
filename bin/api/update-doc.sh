#!/usr/bin/env bash

EDV_ID=$1
DOC_ID=$2

curl -v \
-H "Content-Type: application/json; charset=utf-8" \
-X POST http://localhost:5000/edvs/$EDV_ID/docs/$DOC_ID \
-d @- << DATA
{
  "id": "cc34bc28-f7cf-11ec-a1a2-784f43a5dedc",
  "sequence": 1,
  "indexed": [
    {
      "sequence": 1,
      "hmac": {
        "id": "https://example.com/kms/z7BgF536GaR",
        "type": "Sha256HmacKey2019"
      },
      "attributes": [
        {
          "name": "DUQaxPtSLtd8L3WBAIkJ4DiVJeqoF6bdnhR7lSaPloZ",
          "value": "RV58Va4904K-18_L5g_vfARXRWEB00knFSGPpukUBro",
          "unique": true
        },
        {
          "name": "CUQaxPtSLtd8L3WBAIkJ4DiVJeqoF6bdnhR7lSaPloZ",
          "value": "QV58Va4904K-18_L5g_vfARXRWEB00knFSGPpukUBro",
          "unique": false
        },
        {
          "name": "NewIndexName",
          "value": "NewIndexValue",
          "unique": true
        }
      ]
    }
  ],
  "jwe": {
    "protected": "eyJlbmMiOiJDMjBQIn0",
    "recipients": [
      {
        "header": {
          "alg": "A256KW",
          "kid": "https://example.com/kms/zSDn2MzzbxmX"
        },
        "encrypted_key": "OR1vdCNvf_B68mfUxFQVT-vyXVrBembuiM40mAAjDC1-Qu5iArDbug"
      }
    ],
    "iv": "i8Nins2vTI3PlrYW",
    "ciphertext": "Cb-963UCXblINT8F6MDHzMJN9EAhK3I",
    "tag": "pfZO0JulJcrc3trOZy8rjA"
  }
}
DATA
