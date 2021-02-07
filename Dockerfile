FROM golang:1.15

LABEL maintainer="OMNI-3 http://wiki.corp.dev.vtb/pages/viewpage.action?pageId=67138642"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]