FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway ./cmd/main.go

FROM gcr.io/distroless/static-debian12
COPY --from=builder /gateway /gateway
COPY routing.yaml /routing.yaml
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/gateway"]
