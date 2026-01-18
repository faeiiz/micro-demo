FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pehchan

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/pehchan .
EXPOSE 8081
CMD ["./pehchan"]
