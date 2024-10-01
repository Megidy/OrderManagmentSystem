FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download

# build
COPY /app ./ 
RUN go build -o ./bin/app cmd/main.go

FROM alpine


WORKDIR /
COPY --from=builder /usr/local/src/bin/app / 



CMD ["/app"]