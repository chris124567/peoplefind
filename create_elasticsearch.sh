#!/bin/sh
set -eu

node_name="peoplefind"

curl -X PUT "localhost:9200/${node_name}?pretty" -H "Content-Type: application/json" -d'
{
    "settings" : {
        "index" : {
            "number_of_shards" : 1, 
            "number_of_replicas" : 0
        }
    }
}
'
