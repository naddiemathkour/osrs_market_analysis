FROM node:22-alpine

WORKDIR /osrsflip

COPY . .

# Go
WORKDIR /osrsflip/srv/install

RUN apk add --no-cache git make musl-dev go

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV PATH=/go/bin:${PATH}

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR /osrsflip/srv

# Run Go server
RUN go mod tidy
RUN go build -o serv .
RUN nohup ./serv &

# Run
WORKDIR /osrsflip

# Angular RUN Commands
RUN npm install -g @angular/cli
RUN npm install

# Host Angular server
CMD [ "ng", "serve", "--host", "0.0.0.0" ]