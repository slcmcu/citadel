FROM crosbymichael/golang

RUN go get github.com/codegangsta/cli
ADD . /go/src/github.com/citadel/citadel

RUN cd /go/src/github.com/citadel/citadel && \
    go get -d && go install . ./...

