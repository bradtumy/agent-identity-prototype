# Broker service
FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o broker ./broker

FROM alpine:3.19
WORKDIR /app
COPY --from=build /src/broker /usr/local/bin/broker
COPY data /data
EXPOSE 8081
CMD ["/usr/local/bin/broker"]
