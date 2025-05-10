# ------------ Stage 1: Build Go binary ------------
FROM golang:1.23.6 AS build
#RUN apk --no-cache add gcc g++ make ca-certificates
# Set working directory
WORKDIR /go/src/github.com/Sumitk99/CloudRunr/build-server

# Copy Go source code
COPY . .
RUN go mod download

RUN GO111MODULE=on go build -o /go/bin/app .



# ------------ Stage 2: Runtime Environment ------------
FROM ubuntu:jammy

# Install dependencies
RUN apt-get update && \
    apt-get install -y curl git && \
    curl -sL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean

# Set working directory
WORKDIR /home/app

# Copy compiled Go binary from builder stage
COPY --from=build /go/bin .
COPY .env .
# Copy and set up the startup script
COPY main.sh .
RUN chmod +x main.sh

# Set entrypoint
ENTRYPOINT ["/home/app/main.sh"]