#!/bin/bash

kubectl run kafka-producer -ti --image=quay.io/strimzi/kafka:latest-kafka-4.1.1 --rm=true --restart=Never -- \
  bin/kafka-console-producer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic
