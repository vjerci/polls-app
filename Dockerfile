# FROM node:20-alpine3.17 as frontbuilder
# WORKDIR /app
# COPY ./front ./
# COPY ./front/.env.docker ./.env.local
# RUN npm i && \
#     npm run build


FROM golang:1.20-alpine as backendbuilder
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /app/server ./cmd/server/server.go


FROM alpine:3.18
WORKDIR /app
# RUN mkdir -p /app/front/dist
# COPY --from=frontbuilder /app/dist /app/front/dist/
COPY ./config.yml /app/config.yml
COPY --from=backendbuilder /app/server /app/server
RUN chmod 755 /app/server

EXPOSE 8000
CMD [ "/app/server" ]


