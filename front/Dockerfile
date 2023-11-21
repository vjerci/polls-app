FROM node:20-alpine3.17 as frontbuilder
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_GOOGLE_CLIENT_ID
WORKDIR /app
COPY . .
RUN npm i && \
    npm run build

FROM nginx:alpine
WORKDIR /app
COPY --from=frontbuilder ./app/out ./static
COPY ./nginx.conf /etc/nginx/nginx.conf