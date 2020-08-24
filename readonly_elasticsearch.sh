#!/bin/sh
set -eu
index="peoplefind"

curl -H 'Content-Type: application/json'  -XPUT "http://localhost:9200/${index}/_settings" -d'
{
  "index": {
    "blocks.read_only": true
  }
}' 
