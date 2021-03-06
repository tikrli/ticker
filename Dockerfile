FROM golang:alpine AS build-env
WORKDIR /go/src/github.com/systemli/ticker
ADD . /go/src/github.com/systemli/ticker
RUN go build -o /ticker

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /ticker /ticker

EXPOSE 8080
ENTRYPOINT ["/ticker"]
