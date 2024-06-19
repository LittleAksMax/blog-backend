FROM golang:alpine as build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /app/bin/ ./...

FROM alpine:latest as stage

EXPOSE 8080
WORKDIR /app/

COPY --from=build /app/bin ./bin

ENTRYPOINT ["/app/bin/blog"]
