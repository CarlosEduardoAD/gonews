FROM golang:bookworm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 3000

CMD ["tail", "-f", "/dev/null"]