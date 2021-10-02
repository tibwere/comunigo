#!/bin/sh

curl -X POST http://localhost:$1/sign -H "Accept: application/json" -d "username=$2"