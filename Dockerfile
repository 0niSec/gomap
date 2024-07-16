# Stage 1
FROM golang:1.22.5-alpine3.20 AS build
RUN apk add libpcap-dev build-base
WORKDIR /build/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags "-s -w -linkmode external -extldflags='-static'" -o gomap .

FROM scratch
COPY --from=build /build/src/gomap .
ENTRYPOINT [ "./gomap" ]
