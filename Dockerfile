FROM resin/raspberrypi3-ubuntu-buildpack-deps as builder

RUN [ "cross-build-start" ]

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository -y ppa:gophers/archive
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y golang

RUN go version

RUN mkdir -p /app/build/
WORKDIR /app/build/
COPY . /app/build/

RUN make build_linux

RUN stat /app/build/bin/rusted

RUN [ "cross-build-end" ]

FROM resin/raspberrypi3-ubuntu

ENV INITSYSTEM on

RUN [ "cross-build-start" ]

RUN apt-get update && apt-get install -y \
	screen \
	ppp \
	vim \
	&& rm -rf /var/lib/apt/lists/*

RUN [ "cross-build-end" ]

COPY resources/modem/ppp/mint /etc/ppp/peers/mint
COPY resources/modem/chatscripts/mint /etc/chatscripts/mint

COPY --from=builder /app/build/bin/rusted /usr/local/bin/rusted

COPY resources/start.sh /start.sh

CMD ["/start.sh"]

