FROM alpine:latest

# Install Go
WORKDIR /srv/install

RUN apk add --no-cache git make musl-dev go

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV PATH=/go/bin:${PATH}

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Copy files
WORKDIR /srv

COPY . .

# Clean and build
RUN go mod tidy
RUN go build -o serv .

EXPOSE 8080

# Run
CMD [ "nohup", "./serv", "2>&1", "&" ]