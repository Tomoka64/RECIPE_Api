FROM golang:1.10-alpine

COPY . /go/src/github.com/Tomoka64/RECIPE_Api

WORKDIR /go/src/github.com/Tomoka64/RECIPE_Api

RUN \
        echo $GOPATH && \
        apk add --no-cache bash curl git openssh && \
        curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN dep ensure -v

RUN go install ./...
    
CMD RECIPE_Api
