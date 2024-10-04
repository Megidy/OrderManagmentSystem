# First stage: builder
FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download

# build producer (main program)
COPY /app ./ 
RUN go build -o ./bin/app cmd/main.go

# build consumer (kitchen service)
RUN go build -o ./bin/kitchen pkg/kitchen/consumer/main.go

# Second stage: runtime
FROM alpine

WORKDIR /


COPY --from=builder /usr/local/src/bin/app /app
COPY --from=builder /usr/local/src/bin/kitchen /kitchen


CMD ["/app"]
