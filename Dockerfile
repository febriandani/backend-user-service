# Docker build GO
FROM golang:1.19-alpine AS builder
ARG USER_SERVICE_ADDRESS
# SET TZ
RUN apk add -U tzdata
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
WORKDIR /app
COPY . .
# Download Go modules
RUN go mod download
# Separate the build steps and handle symbolic links carefully
RUN go build -o ./backend-users-service ./cmd/server/main.go \
    && go build -ldflags  "-X main.userServiceAddr=$USER_SERVICE_ADDRESS" -o ./gateway-service ./cmd/client/main.go


# Docker build backend-users-service
FROM alpine:latest AS backend-users-service
WORKDIR /app
# Copy only the necessary binary file from the builder stage
COPY --from=builder /app/backend-users-service .
# Expose port and define entry point
EXPOSE 50051
ENTRYPOINT ["./backend-users-service"]


# Docker build gateway-service
FROM alpine:latest AS gateway-service
WORKDIR /app
# Copy only the necessary binary file from the builder stage
COPY --from=builder /app/gateway-service .
# Expose port and define entry point
EXPOSE 8080
ENTRYPOINT ["./gateway-service"]
