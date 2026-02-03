all: image-all

image-all:
	docker build -t ttl.sh/kafka-connect-manual-5e3e6cec:24h .
	docker push ttl.sh/kafka-connect-manual-5e3e6cec:24h
