# FROM golang:1.22.4 as build
# WORKDIR /build/src
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .
# RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gomap .

# We do not need to build the binary in the container anymore
# This is all taken care of by the CI/CD pipeline
FROM scratch
COPY gomap /usr/bin/gomap
ENTRYPOINT [ "/usr/bin/gomap" ]