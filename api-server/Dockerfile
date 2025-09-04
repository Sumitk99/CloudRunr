FROM golang:1.23.3 AS build
RUN apt-get update && apt-get install -y gcc g++ make ca-certificates
WORKDIR /go/src/github.com/Sumitk99/cloudrunr/api-server
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GO111MODULE=on go build -o /build/api-server ./cmd/api-server

FROM ubuntu:22.04
ENV PORT=8080
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /usr/bin
COPY --from=build /build/api-server .
COPY .env .
COPY client.properties .
EXPOSE 8080
CMD ["./api-server"]
