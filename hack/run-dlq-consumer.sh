kubectl run kafka-consumer -ti --image=quay.io/strimzi/kafka:latest-kafka-3.7.0 --rm=true --restart=Never -- \
  bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-dlq-topic --from-beginning
