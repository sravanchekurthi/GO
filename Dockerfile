# syntax=docker/dockerfile:1

# ---------- Stage 1: Build ----------
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install git (required for go get), and build tools
RUN apk add --no-cache git

# Copy go files
COPY go.mod ./
COPY go.sum ./

ENV GOPROXY=https://proxy.golang.org,direct

RUN ls -la
RUN cat go.mod
RUN cat go.sum

RUN go mod download -x

COPY . ./

# Build the binary
RUN go build -o nsctl .

# ---------- Stage 2: Minimal runtime ----------
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/nsctl .

# Run CLI tool by default
ENTRYPOINT ["./nsctl"]    