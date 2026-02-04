FROM quay.io/strimzi/kafka:latest-kafka-4.1.1
USER root:root
RUN mkdir -p /opt/kafka/custom-config/ && touch /opt/kafka/custom-config/log4j.properties
RUN mkdir -p /opt/kafka/plugins/s3-connector
COPY ./files/confluentinc-kafka-connect-s3-12.0.1.zip /tmp/
RUN unzip /tmp/confluentinc-kafka-connect-s3-12.0.1.zip -d /opt/kafka/plugins/s3-connector
RUN curl -s https://repo1.maven.org/maven2/org/apache/kafka/connect-file/4.1.1/connect-file-4.1.1.jar -o /opt/kafka/plugins/connect-file-4.1.1.jar
USER 1001
