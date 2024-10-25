# Use Ubuntu as the base image
FROM ubuntu:latest

# Install necessary dependencies
RUN apt-get update && apt-get install -y \
    golang-go \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source file to the working directory
COPY . .

# Build the Go file
RUN go build -o pingpong1 pingpong1.go

# Run the compiled Go binary
CMD ["./pingpong1"]
