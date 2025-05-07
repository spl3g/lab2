# STEP-1
# build app from source

FROM golang:1.24.1-alpine3.21 AS builder

WORKDIR /myapp

COPY ./services/ ./services/

COPY ./vendor* ./vendor/
COPY ./go.mod ./go.sum ./

COPY ./internal ./internal

RUN go build -o app ./services/permissions-service/main.go


# STEP-2
# make container

FROM alpine:3.21

WORKDIR /mysuperapp

COPY --from=builder /myapp/app ./

EXPOSE 10000

CMD [ "/mysuperapp/app" ]
