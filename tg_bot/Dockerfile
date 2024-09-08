FROM golang:1.23-alpine3.20 AS build

ENV CGO_ENABLED=0

WORKDIR /src

COPY . .

RUN go build -o /go/bin/bff

FROM alpine:3.20

COPY --from=build /go/bin/bff /usr/local/bin/bff

ENTRYPOINT [""]

CMD ["/usr/local/bin/bff"]
