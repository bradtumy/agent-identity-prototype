FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o agent_runner ./agent_runner

FROM alpine:3.19
COPY --from=build /src/agent_runner /usr/local/bin/agent_runner
CMD ["/usr/local/bin/agent_runner"]
