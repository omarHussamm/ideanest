# Use an official Golang runtime as a base image
FROM golang:1.21.5

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the entire project to the container
COPY . .

# Download and install dependencies
RUN go mod init

# Build the Go application
# RUN go build -o app ./cmd

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./app"]