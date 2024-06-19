FROM golang:1.19

WORKDIR /go-bike

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o go-bike cmd/main.go

EXPOSE 8080

CMD [ "./go-bike" ]