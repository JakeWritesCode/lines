# Start from the official golang base image
FROM golang:1.22

# Add Maintainer info
LABEL maintainer="Your Name <your.email@example.com>"

# Set the current working directory in the container
WORKDIR /app

# Copy the source from the current directory to the working directory in the container
COPY . .

# Download all the dependencies
RUN go mod download

# Install the package
RUN go build -o main ./cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8100

# Command to run the executable
CMD ["./main"]