FROM golang

COPY ./ /app

WORKDIR /app

RUN go get github.com/cactus/go-statsd-client/statsd

CMD ["go","run","main.go"]