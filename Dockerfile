FROM golang:1.21.5

WORKDIR /app

ADD . /lib

COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go mod download

COPY ./lib ./lib

RUN go build .

CMD [ "./himbot" ]