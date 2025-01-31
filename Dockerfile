
FROM golang:latest as build
WORKDIR /usr/src/app
COPY . .
RUN go build -o main .
FROM ubuntu:latest
RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /usr/src/app/main .
CMD ["./main"]