#!/usr/bin/env bash

curl -v \
-H "Content-Type: application/json; charset=utf-8" \
-X POST http://localhost:5000/edvs \
-d @- << DATA
{
  "sequence": 0,
  "controller": "did:example:123456789",
  "referenceId": "abc5a436-21f9-4b4c-857d-1f5569b2600d",
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
