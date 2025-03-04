# syntax=docker/dockerfile:1

FROM golang:1.24

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./
COPY billing/* ./billing/
COPY shipment/* ./shipment/
COPY acceptance-test/* ./acceptance-test/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./oms-worker

# Run
CMD [ "./oms-worker" ]
