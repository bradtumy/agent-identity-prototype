FROM golang:1.24-alpine

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ✅ Correctly build a binary from the main source file
RUN CGO_ENABLED=0 go build -o /usr/local/bin/broker ./broker/main.go

CMD ["/usr/local/bin/broker"]
