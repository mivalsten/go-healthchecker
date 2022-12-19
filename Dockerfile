############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
COPY . /go
WORKDIR /go/src
# Build the binary.
RUN CGO_ENABLED=0 go build .
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/src/healthchecker /healthchecker
# Run the hello binary.
ENTRYPOINT ["/healthchecker"]
#CMD ["sh"]