# Use the official golang:alpine image as the base image
FROM golang:alpine 

# Install necessary build tools and dependencies
RUN apk update && apk add --no-cache git gcc g++ sqlite-dev

# Set the working directory inside the container
WORKDIR /app

# Copy Go application source code into the container
COPY . .

# Build Go application with CGO enabled
RUN CGO_ENABLED=1 go build -o main ./cmd/api

# Expose the port application runs on
EXPOSE 8080

# Run the Go application
CMD ["./main"]
