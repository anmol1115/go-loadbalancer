FROM golang:1.23-alpine

WORKDIR /app

COPY ./src/loadbalancer/go.* .
RUN go mod download

COPY ./src/loadbalancer/ .
RUN go build -o loadBalancer .

RUN chmod +x ./loadBalancer
CMD ["./loadBalancer"]
