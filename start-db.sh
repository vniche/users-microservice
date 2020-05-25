#!/bin/bash

docker run -d --rm --name cete-node1 \
    -p 7000:7000 \
    -p 8000:8000 \
    -p 9000:9000 \
    vniche/cete:latest cete start \
      --id=node1 \
      --raft-address=:7000 \
      --grpc-address=:9000 \
      --http-address=:8000 \
      --data-directory=/tmp/cete/node1