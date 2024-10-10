FROM golang:alpine AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /app/bin/ ./...

FROM alpine:latest AS stage

EXPOSE 8080
WORKDIR /app/

COPY --from=build /app/bin ./bin

# Create system user with no home directory
RUN adduser -g -S -D -H -s /sbin/nologin drb

# Own the application directory
RUN chown drb /app && chgrp drb /app

# Use system user for process
USER drb

ENTRYPOINT ["/app/bin/blog"]
