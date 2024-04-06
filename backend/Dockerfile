FROM golang:1.21

WORKDIR /backend

COPY go.mod go.sum .

RUN go mod download

COPY . .

EXPOSE 8000

RUN CGO_ENABLED=1 GOOS=linux go build -o /backend-root

CMD ["/backend-root"]