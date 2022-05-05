#!/usr/bin/env bash

EDV_ID=$1

curl -v \
-H "Content-Type: application/json; charset=utf-8" \
-X POST http://localhost:$EDV_PORT/edvs/$EDV_ID/docs \
-d @- << DATA
{
  "id": "51686fbc-cc46-11ec-9940-784f43a5dedc",
  "sequence": 0,
  "indexed": [
    {
      "sequence": 0,
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
