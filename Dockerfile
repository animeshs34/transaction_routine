# --- Test Stage ---
FROM golang:1.23-alpine AS test
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

RUN go test -v ./... -coverprofile=coverage.out

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/api ./cmd/api

# Runtime stage
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
COPY --from=builder /bin/api /bin/api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/bin/api"]
