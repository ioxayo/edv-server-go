#!/bin/bash

curl -v -X POST http://localhost:$EDV_PORT/edvs \
-H "Content-Type: application/json; charset=utf-8" \
-d @- << DATA
{
  "sequence": 0,
  "controller": "did:example:123456789",
  "keyAgreementKey": {
    "id": "https://example.com/kms/12345",
    "type": "X25519KeyAgreementKey2019"
  },
  "hmac": {
    "id": "https://example.com/kms/67891",
    "type": "Sha256HmacKey2019"
  }
}
DATA

