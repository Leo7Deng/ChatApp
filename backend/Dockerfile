# syntax=docker/dockerfile:1

FROM golang:1.23.3

# Set destination for COPY
WORKDIR /backend

RUN apt update && apt install -y postgresql-client

# Download Go modules
COPY go.mod go.sum ./
COPY . .
RUN go mod tidy

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /chatapp

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/chatapp"]