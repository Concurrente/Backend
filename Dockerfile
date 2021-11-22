FROM golang:1.15-alpine AS GO_BUILD
COPY . /server
WORKDIR /server/server
RUN go build -o /go/bin/server/server

FROM alpine:3.10
WORKDIR app
COPY --from=GO_BUILD /go/bin/server/ ./
EXPOSE 9000
CMD ./server