# Start by building the application.
FROM golang:1.23.3-bookworm AS builder

WORKDIR /workdir
COPY . /workdir

ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
ENV CGO_ENABLED=0

# Build the binaries
RUN go build -o main .


FROM alpine:3.19.1

# app's parameter
EXPOSE 8080

# Copy assets and binaries
COPY --from=builder /workdir/main /tetris

WORKDIR /
# Start the app
ENTRYPOINT ["/tetris"]
