FROM golang:1.10.4-alpine3.8 as builder

ENV BASE_APP_DIR /go/src/github.com/mszostok/oauth-less
WORKDIR ${BASE_APP_DIR}

COPY . ${BASE_APP_DIR}/

RUN go build -v -o main .
RUN mkdir /app && mv ./main /app/main

FROM alpine:3.8
LABEL source = git@github.com:mszostok/oauth-less.git
WORKDIR /app

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app /app

CMD ["/app/main"]