FROM golang:alpine3.13

RUN apk --update upgrade
RUN apk add build-base
RUN apk add sqlite

RUN rm -rf /var/cache/apk/*
WORKDIR src
COPY *.go .
COPY suss.sqlite .

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o ./out/suss-app .

EXPOSE 5000

CMD ["./out/suss-app"]