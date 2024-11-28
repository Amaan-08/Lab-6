# Stage 1: Go Application
FROM golang:1.20 AS build

# Set working directory
WORKDIR /app

# Copy Go files
COPY . .

# Install dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Stage 2: MySQL
FROM mysql:8.0

# Set environment variables for MySQL
ENV MYSQL_ROOT_PASSWORD=rootpassword
ENV MYSQL_DATABASE=toronto_time_db

# Expose MySQL port
EXPOSE 3306

# Stage 3: Combine Go Application and MySQL

# Use an empty base to run both services
FROM ubuntu:latest

# Install MySQL client and other dependencies
RUN apt-get update && \
    apt-get install -y mysql-client && \
    apt-get install -y curl && \
    apt-get install -y libmariadb-dev

# Copy Go application from the build stage
COPY --from=build /app/main /app/main

# Set working directory
WORKDIR /app

# Expose Go application port
EXPOSE 8080

# Run Go application and MySQL in the same container
CMD service mysql start && ./main
