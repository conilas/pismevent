FROM golang:1.11

WORKDIR /go/src/pismevent
COPY . .

RUN cd /go/src/pismevent
RUN go get -v
RUN go build

EXPOSE 3031

CMD [ "pismevent" ]


