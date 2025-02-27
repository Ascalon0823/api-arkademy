FROM golang:alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app

FROM alpine

WORKDIR /usr/src/app
RUN touch /usr/local/bin/.env && \
    echo MONGO_URI_LOCAL=${MONGO_URI_LOCAL} >> /usr/local/bin/.env && \
    echo JWT_SECRET_LOCAL=${JWT_SECRET} >> /usr/local/bin/.env && \
    cat /usr/local/bin/.env
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

EXPOSE 8080

CMD ["app"]