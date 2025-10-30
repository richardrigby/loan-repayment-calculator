# Multi-stage Dockerfile for building the Go app
FROM golang:1.25-alpine AS builder

ARG TARGETARCH
ARG TARGETOS

ENV CGO_ENABLED=0 \
    GOARCH=$TARGETARCH \
    GOOS=$TARGETOS \
    DOCKER_BUILDKIT=1

WORKDIR /src

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the repo and build the binary
COPY . .
RUN go build -o loan-repayment-calculator cmd/loan-repayment-calculator/main.go

FROM alpine:3.19
RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /src/loan-repayment-calculator /usr/local/bin/loan-repayment-calculator
RUN chmod +x /usr/local/bin/loan-repayment-calculator

USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/loan-repayment-calculator"]
