# Start the Go app build
FROM golang:latest AS build

# Copy source
WORKDIR /go/src/CityBikeApi/db-table
COPY . .

# Get required modules (assumes packages have been added to ./vendor)
RUN go get -d -v ./...

# Build a statically-linked Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -o cb-api-table .

# New build phase -- create binary-only image
FROM alpine:latest

# Add support for HTTPS and time zones
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates

WORKDIR /root/

# Copy files from previous build container
COPY --from=build /go/src/CityBikeApi/db-table/cb-api-table ./

ENV AWS_ACCESS_KEY_ID=AKIA34XNLPJYLZGHNQGI
ENV AWS_SECRET_ACCESS_KEY=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
# Add environment variables
# ENV ...

# Check results
# RUN env && pwd && find .

# Start the application
# CMD ["go", "run", "cb-api-table.go"]
CMD ["./cb-api-table"]