FROM golang:1.24-alpine AS builder

# Copy the source code into the container
RUN mkdir -p /build
COPY . /build

# Set the working directory
WORKDIR /build

# Download dependencies
RUN go mod download

# Build the Go application
RUN go build -o ./build/glance-agent

# Install UPX for binary compression
RUN apk add --no-cache upx

RUN upx --best ./build/glance-agent

# Use Alpine as the base image
FROM alpine:latest

# Set environment variables to avoid prompts
ENV TZ=UTC

# Install packages (e.g., bash and curl)
RUN apk update && \
  apk add --no-cache bash curl ca-certificates

RUN mkdir -p /app

# Copy the application
COPY --from=builder /build/build/glance-agent /app/glance-agent

# Set working directory
WORKDIR /app

# Set default command
CMD ["/app/glance-agent"]
EXPOSE 9012