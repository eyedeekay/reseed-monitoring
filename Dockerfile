FROM debian:stable-backports
ADD . /home/user/go/src/github.com/eyedeekay/reseed-monitoring
RUN apt-get update && \
	apt-get dist-upgrade -y && \
	apt-get install -y golang-any git make ca-certificates && \
	adduser --disabled-password --gecos 'user,,,,' user && \
	chown -R user:user /home/user/ && \
	cd /home/user/go/src/github.com/eyedeekay/reseed-monitoring && \
	git remote set-url origin https://github.com/eyedeekay/reseed-monitoring
USER user
WORKDIR /home/user/go/src/github.com/eyedeekay/reseed-monitoring
RUN make build
CMD /home/user/go/src/github.com/eyedeekay/reseed-monitoring/reseed-monitor/reseed-monitor
