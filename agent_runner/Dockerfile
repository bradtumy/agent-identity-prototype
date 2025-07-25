FROM golang:1.24-alpine

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ✅ Build the runner binary
RUN CGO_ENABLED=0 go build -o /usr/local/bin/agent_runner ./agent_runner/main.go

CMD ["/usr/local/bin/agent_runner"]
