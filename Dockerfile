# Using Debian-based Go image 
FROM golang:1.23.2

# Enable CGO
ENV CGO_ENABLED=1

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum before copying source
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the full project source
COPY . .

# Build the Go application
RUN go build -o server ./cmd/main.go

# Expose port
EXPOSE 8081

# Run the app
CMD ["./server"]
