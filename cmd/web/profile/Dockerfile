# Use the official Golang image as a base image
FROM golang:1.24.2-alpine AS builder

# Set the working directory inside the container
WORKDIR /build

# Copy go module files
#COPY go.mod go.sum ./

# Download dependencies
#RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Create a new stage for the final image
FROM alpine:latest
VOLUME /profile
# Set the working directory inside the container
WORKDIR /cmd/web/profile
USER root
RUN mkdir /etc/tarantula 
COPY profile_conf.json /etc/tarantula/profile_conf.json

# Copy the binary from the builder stage
COPY --from=builder /build/main .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the application
CMD ["./main"]