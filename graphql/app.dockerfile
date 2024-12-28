FROM golang:1.22.5-alpine AS builder

# Install required packages and set the working directory
RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /go/src/github.com/Anshu-rai89/go-ms

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY account account
COPY catalog catalog
COPY order order
COPY graphql graphql
RUN go build -o /go/bin/app ./graphql

# Step 2: Runtime stage
FROM alpine:3.11

# Set working directory for the runtime container
WORKDIR /usr/bin

# Copy the binary from the build stage
COPY --from=builder /go/bin/app .

# Expose port 8080 for the application
EXPOSE 8080

# Ensure the application is executed with the right command
CMD ["/usr/bin/app"]