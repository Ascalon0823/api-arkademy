FROM golang:alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -v -o /usr/local/bin/app

FROM alpine

WORKDIR /usr/src/app
RUN touch .env && \
    echo MONGO_URI=${MONGO_URI} >> .env && \
    echo JWT_SECRET=${JWT_SECRET} >> .env && \
    echo SERVER_PORT=1027 >> .env && \
    cat .env
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

EXPOSE 1027

CMD ["app"]