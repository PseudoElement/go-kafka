FROM golang:1.24

# Set destination for COPY
WORKDIR /app

# Copy the entire project
COPY . .

# Download dependencies
RUN go mod tidy
RUN go mod download

# Build
RUN go build -o go-kafka

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["./go-kafka", "10"]