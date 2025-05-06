# Base image with Go support and necessary build tools
FROM golang:1.24-bullseye

# Install C build tools and dependencies for CGO/Turso
RUN apt-get update && apt-get install -y \
    build-essential \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Enable CGO and set ldflags for proper linking
ENV CGO_ENABLED=1
ENV CGO_LDFLAGS="-ldl"

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o main .

# Expose the port your application runs on
# (Replace 8080 with your actual port)
EXPOSE 8080

# Command to run the application
CMD ["./main"]
