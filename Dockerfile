FROM golang:1.24.1-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/swap-ai App/*.go

FROM alpine:3.19

EXPOSE 8080

RUN mkdir -p /usr/local/bin/ /server/ 

COPY --from=build /app/swap-ai /server/swap-ai

CMD ["/server/swap-ai"]