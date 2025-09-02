FROM golang:1.23.3-alpine3.19 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/Sumitk99/cloudrunr/api-server
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GO111MODULE=on go build -o /build/api-server ./cmd/api-server

FROM alpine:3.19
ENV PORT=8080
WORKDIR /usr/bin
COPY --from=build /build/api-server .
COPY .env .
COPY client.properties .
EXPOSE 8080
CMD ["./api-server"]