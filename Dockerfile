# Use the official Go image as the base image
FROM golang:1.23.1-alpine3.19

RUN apk add --no-cache chromium \
    && apk add --no-cache bash \
    && apk add --no-cache --virtual .build-deps gcc g++ make

# Set environment variables for Chrome
ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/


# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]