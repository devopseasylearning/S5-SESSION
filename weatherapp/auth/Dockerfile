#build stage
FROM golang:alpine AS builder
RUN apk add git
WORKDIR /go/src/app
RUN git clone https://github.com/abohmeed/weatherapp
WORKDIR /go/src/app/weatherapp/auth/src/main
RUN go build -o /go/bin/app -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=auth Version=0.0.1
EXPOSE 8080
