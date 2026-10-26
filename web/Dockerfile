FROM node:24-alpine AS build

ARG VITE_API_TOKEN
ENV VITE_API_TOKEN=$VITE_API_TOKEN

WORKDIR /app
COPY web/package.json ./package.json
RUN npm install
COPY web/ ./
RUN npm run build

FROM caddy:2-alpine
COPY web/Caddyfile /etc/caddy/Caddyfile
COPY --from=build /app/dist /srv
EXPOSE 80
