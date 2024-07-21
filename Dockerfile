FROM node:22-alpine

WORKDIR /osrsflip

COPY . .

# Angular RUN Commands
RUN npm install -g @angular/cli
RUN npm install

# Go
WORKDIR /osrsflip/srv

RUN apk add --no-cache git make musl-dev go

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV PATH=/go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

RUN go mod tidy
RUN go build -o serv .
RUN ./serv

# Run project
WORKDIR /osrsflip

CMD [ "ng", "serve", "--host", "0.0.0.0" ]