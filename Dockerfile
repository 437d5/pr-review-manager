FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app ./cmd/reviewer

FROM alpine:3.20
COPY --from=build /app /app

ENTRYPOINT [ "/app" ]