FROM	ubuntu:14.04
MAINTAINER	Seongyeol Lim <seongyeol37@gmail.com>

COPY	.	/go/src/github.com/dwake/docker-trust/external/github.com/docker/docker
ADD		.	/
ADD		null /
COPY	nullfile /tmp
ADD		[ "vimrc", "/tmp" ]
COPY	[ "bashrc", "/tmp" ]
COPY	[ "test file", "/tmp" ]
ADD		[ "test file", "/tmp/test file" ]
