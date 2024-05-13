# Start from the official Golang base image
FROM golang:1.21.1-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code and public folder into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file and the public folder from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/public ./public

# Expose port 1323 to the outside world
EXPOSE 1323

# Command to run the executable
CMD ["./main"]
