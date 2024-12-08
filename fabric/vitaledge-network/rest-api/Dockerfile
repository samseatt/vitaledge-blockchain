# Dockerfile for vitaledge REST API

# Base image
FROM golang:1.23.3 as builder

# Set working directory inside the container
WORKDIR /app

# Copy the source code and Go modules
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the REST API binary
RUN go build -o rest-api .

# Create a minimal runtime image
# FROM gcr.io/distroless/base-debian10
# FROM debian:bullseye-slim AS runtime

# Runtime image (match the glibc version of golang:1.23.3)
FROM golang:1.23.3

# Copy the compiled binary
COPY --from=builder /app/rest-api /

# Expose the port used by the REST API
EXPOSE 8082

# Start the REST API
CMD ["/rest-api"]
