FROM golang:1.14

RUN mkdir /app

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 8080

RUN go build

CMD /app/go_clean_arc