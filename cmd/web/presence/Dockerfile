# Use the official Golang image as a base image
FROM golang:1.24.2-alpine AS builder

# Set the working directory inside the container
WORKDIR /build

# Copy go module files
# COPY go.mod go.sum ./
COPY . .
# Download dependencies
#RUN go mod download

# Copy the source code into the container
#COPY . .
# Build the Go application
#RUN cd cmd/web/presence
RUN go -C cmd/web/presence build -o main 

# Create a new stage for the final image
FROM alpine:latest
VOLUME /presence
USER root
RUN mkdir /etc/tarantula 
COPY cmd/web/presence/presence-conf.json /etc/tarantula/presence-conf.json

# Set the working directory inside the container
WORKDIR /cmd/web/presence

# Copy the binary from the builder stage
COPY --from=builder /build/cmd/web/presence/main .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the application
CMD ["./main"]