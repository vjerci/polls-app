FROM golang:1.20-alpine as apibuilder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /app/server ./cmd/server/server.go

FROM alpine:3.18
WORKDIR /app
COPY ./config.yaml /app/config.yaml
COPY --from=apibuilder /app/server /app/server
RUN chmod 755 /app/server

CMD [ "/app/server" ]


