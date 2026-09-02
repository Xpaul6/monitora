FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY controllers/ ./controllers/
COPY models/ ./models/
COPY utils/ ./utils/

RUN CGO_ENABLED=0 GOOS=linux go build -o /monitora

EXPOSE 8080

CMD ["/monitora"]
