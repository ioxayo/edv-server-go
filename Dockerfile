# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /edv-server-go

ENV PORT=5000
EXPOSE 5000

CMD [ "/edv-server-go" ]
