FROM golang:1.10-alpine

ADD ./web /home

WORKDIR /home

RUN \
       apk add --no-cache bash git openssh && \
       go get -u github.com/julienschmidt/httprouter

CMD ["go","run","main.go"]
