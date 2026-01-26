FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o darwaza

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/darwaza .
EXPOSE 8080
CMD ["./darwaza"]
