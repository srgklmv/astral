FROM golang:1.24.3-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o ./tmp/main cmd/astral/main.go

FROM alpine:latest AS production

WORKDIR /app

COPY --from=build /app/config.json /app/config.json
COPY --from=build /app/migrations /app/migrations
COPY --from=build /app/tmp/main /app/main

EXPOSE 3000

CMD ["./main"]