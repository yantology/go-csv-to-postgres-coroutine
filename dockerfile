# Use the official Golang image as the base image
FROM golang:1.22.4-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o gin-go-PostgresSQL-Bookstore-Management-Api main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./gin-go-PostgresSQL-Bookstore-Management-Api/cmd"]