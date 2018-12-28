FROM resin/raspberrypi3-ubuntu-buildpack-deps as builder

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository -y ppa:gophers/archive
RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y golang

RUN go version

FROM resin/raspberrypi3-ubuntu

ENV INITSYSTEM on

RUN apt-get update && apt-get install -y \
	screen \
	ppp \
	vim \
	&& rm -rf /var/lib/apt/lists/*

COPY resources/modem/ppp/twilio /etc/ppp/peers/twilio
COPY resources/modem/chatscripts/twilio /etc/chatscripts/twilio

COPY resources/start.sh /start.sh
CMD ["/start.sh"]

