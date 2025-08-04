# ------------ Stage 1: Build Go binary ------------
FROM golang:1.23.6 AS build
WORKDIR /go/src/github.com/Sumitk99/CloudRunr/build-server

# Copy all source code (backend and frontend) to the build container
COPY . .
RUN go mod download
RUN GO111MODULE=on go build -o /go/bin/app .

# ------------ Stage 2: Runtime Environment ------------
FROM ubuntu:jammy

# Install dependencies for Go binary and Node.js for frontend build
RUN apt-get update && \
    apt-get install -y curl git && \
    curl -sL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean

WORKDIR /home/app

# Copy compiled Go binary from builder stage
COPY --from=build /go/bin .
COPY .env .
COPY main.sh .
RUN chmod +x main.sh

# Copy all project files (including frontend and package.json) to the runtime container
#COPY . .

# Install frontend dependencies and build Angular app using local CLI
#RUN npm install
#RUN npx ng build --configuration=production

# Set entrypoint
ENTRYPOINT ["/home/app/main.sh"]

