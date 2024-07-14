FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM alpine:edge

WORKDIR /app

COPY --from=build /app/myapp .

COPY .env /app

EXPOSE 80

ENTRYPOINT ["/app/myapp"]
