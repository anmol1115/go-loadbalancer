FROM golang:1.23-alpine

WORKDIR /app

COPY ./src/backend/go.* .
RUN go mod download

COPY ./src/backend/ .
RUN go build -o backend .

RUN chmod +x ./backend
CMD ["./backend"]
