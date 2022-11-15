FROM golang:1.9.2
ADD . /go/src/task
WORKDIR /go/src/task
RUN go get task
RUN go install
ENTRYPOINT ["/go/bin/task"]