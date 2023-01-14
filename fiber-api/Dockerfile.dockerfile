FROM golang:1.18-alpine3.17

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./


RUN go build -o /contact-api

EXPOSE 3000

CMD [ "/contact-api" ]