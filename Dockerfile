# Stage 1: Build the Go application
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create the final image with Docker-in-Docker
FROM docker:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .

# Set the entrypoint to the Go application
ENTRYPOINT ["./main"]