FROM golang:1.22.4 as build
WORKDIR /build/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gomap .
FROM scratch
COPY --from=build /build/src/gomap /usr/bin/gomap
ENTRYPOINT ["/usr/bin/gomap"]