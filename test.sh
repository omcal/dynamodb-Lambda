#!/bin/bash


for ((i = 0; i < 100; i++)); do
    id="Life_is_$i"

    curl -X POST \
      -H "Content-Type: application/json" \
      -d "{\"id\": \"$id\", \"data\": \"42\"}" \
      localhost:8080/picus/put
    echo "{\"id\": \"$id\", \"data\": \"42\"}" 
    echo
done

