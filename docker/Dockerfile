FROM debian:buster

RUN apt update && \
    apt install -y git gccgo-8 golang

# Neosilka
RUN ln -s /usr/bin/gccgo-8 /usr/bin/gccgo

RUN mkdir -p /var/logc-go
WORKDIR /var/logc-go
COPY . /var/logc-go

ENV GOROOT=/usr/lib/go
ENV GOPATH=/var

RUN go get -u github.com/jmoiron/sqlx
RUN go get -u github.com/kshvakov/clickhouse

RUN go build -compiler=gccgo -gccgoflags='-O2' -o bin/logc

#CMD tail -f /dev/null
CMD /var/logc-go/bin/logc