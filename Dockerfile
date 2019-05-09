FROM golang:1.11

WORKDIR /go/src/eventsourcismo
COPY . .

RUN cd /go/src/eventsourcismo
RUN go get -v
RUN go build

EXPOSE 3031

CMD [ "eventsourcismo" ]


