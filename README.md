# go-logging-demo
Demo app with structured logging and monitoring written in Go

## Endpoints
```
http://localhost:8080/v1/log
http://localhost:8080/v1/log?extEventId=123
http://localhost:8080/metrics
http://localhost:8080/liveness
http://localhost:8080/readiness
```

## Install packages manually
```
go get -v -u github.com/gorilla/mux
go get github.com/google/uuid
go get github.com/sirupsen/logrus
go get github.com/bshuster-repo/logrus-logstash-hook
```

## Build Dicker image
```
docker build -t go-demo-app:1.0.0 .
```

## Run app in Docker
```
docker run -d --name go-demo-app -p 8080:8080 go-demo-app:1.0.0
```

## Run with Docker Compose
```
docker-compose up -d
```

## Run Fluent-bit locally
Run a Fluent Bit instance that will receive messages over TCP port 5170, and send the messages to the STDOUT interface in JSON format every one second:
```
docker run --rm --name fluent-bit -p 127.0.0.1:5170:5170/tcp fluent/fluent-bit:1.6.10-debug /fluent-bit/bin/fluent-bit -i tcp -o stdout -p format=json_lines -f 1
```
