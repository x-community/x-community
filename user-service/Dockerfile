FROM golang:1.14.4 AS builder

WORKDIR /x-community

COPY  . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service .

FROM alpine:3.12.0

COPY --from=builder /x-community/service /usr/local/bin

CMD [ "/usr/local/bin/service" ]