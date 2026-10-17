FROM golang:1.24-alpine AS build

WORKDIR /src
COPY api/go.mod ./go.mod
RUN go mod download
COPY api/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /wallace-api ./cmd/server

FROM alpine:3.22
RUN addgroup -S wallace && adduser -S -G wallace wallace
COPY --from=build /wallace-api /usr/local/bin/wallace-api
USER wallace
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/wallace-api"]
