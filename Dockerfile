FROM golang:alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app

FROM alpine

WORKDIR /usr/src/app

ENV MONGO_URI_LOCAL=${MONGO_URI_LOCAL}
ENV JWT_SECRET=${JWT_SECRET}
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

EXPOSE 8080

CMD ["app"]