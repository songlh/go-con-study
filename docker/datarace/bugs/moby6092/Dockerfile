# This file describes the standard way to build Docker, using docker
#
# Usage:
#
# # Assemble the full dev environment. This is slow the first time.
# docker build -t docker .
#
# # Mount your source in an interactive container for quick testing:
# docker run -v `pwd`:/go/src/github.com/dotcloud/docker --privileged -i -t docker bash
#
# # Run the test suite:
# docker run --privileged docker hack/make.sh test
#
# # Publish a release:
# docker run --privileged \
#  -e AWS_S3_BUCKET=baz \
#  -e AWS_ACCESS_KEY=foo \
#  -e AWS_SECRET_KEY=bar \
#  -e GPG_PASSPHRASE=gloubiboulga \
#  docker hack/release.sh
#
# Note: Apparmor used to mess with privileged mode, but this is no longer
# the case. Therefore, you don't have to disable it anymore.
#

#docker-version	0.6.1
FROM	ubuntu:14.04
MAINTAINER	Tianon Gravi <admwiggin@gmail.com> (@tianon)

# Packaged dependencies
RUN	apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
	apt-utils \
	aufs-tools \
	automake \
	btrfs-tools \
	build-essential \
	curl \
	dpkg-sig \
	git \
	iptables \
	libapparmor-dev \
	libcap-dev \
	libsqlite3-dev \
	lxc=1.0* \
	mercurial \
	pandoc \
	reprepro \
	ruby1.9.1 \
	ruby1.9.1-dev \
	s3cmd=1.1.0* \
	--no-install-recommends

# Get lvm2 source for compiling statically
ENV LVM2_VERSION 2.02.103
RUN mkdir -p /usr/local/lvm2 \
	&& curl -fsSL "https://mirrors.kernel.org/sourceware/lvm2/LVM2.${LVM2_VERSION}.tgz" \
		| tar -xzC /usr/local/lvm2 --strip-components=1
# See https://git.fedorahosted.org/cgit/lvm2.git/refs/tags for release tags

# Compile and install lvm2
RUN cd /usr/local/lvm2 \
	&& ./configure \
		--build="$(gcc -print-multiarch)" \
		--enable-static_link \
	&& make device-mapper \
	&& make install_device-mapper
# See https://git.fedorahosted.org/cgit/lvm2.git/tree/INSTALL

# Compile and install lvm2
RUN	cd /usr/local/lvm2 && ./configure --enable-static_link && make device-mapper && make install_device-mapper
# see https://git.fedorahosted.org/cgit/lvm2.git/tree/INSTALL

# Install Go
RUN	curl -s https://dl.google.com/go/go1.2.2.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV	PATH	/usr/local/go/bin:$PATH
ENV	GOPATH	/go:/go/src/github.com/dotcloud/docker/vendor
RUN	cd /usr/local/go/src && ./make.bash --no-clean 2>&1

# Compile Go for cross compilation
ENV	DOCKER_CROSSPLATFORMS	\
	linux/386 linux/arm \
	darwin/amd64 darwin/386 \
	freebsd/amd64 freebsd/386 freebsd/arm
# (set an explicit GOARM of 5 for maximum compatibility)
ENV	GOARM	5
RUN	cd /usr/local/go/src && bash -xc 'for platform in $DOCKER_CROSSPLATFORMS; do GOOS=${platform%/*} GOARCH=${platform##*/} ./make.bash --no-clean 2>&1; done'

# Grab Go's cover tool for dead-simple code coverage testing
RUN	go get -u golang.org/x/tools/cover

# TODO replace FPM with some very minimal debhelper stuff
RUN	gem install --no-rdoc --no-ri fpm --version 1.0.2

# Get the "busybox" image source so we can build locally instead of pulling
RUN	git clone -b buildroot-2014.02 https://github.com/jpetazzo/docker-busybox.git /docker-busybox

# Setup s3cmd config
RUN	/bin/echo -e '[default]\naccess_key=$AWS_ACCESS_KEY\nsecret_key=$AWS_SECRET_KEY' > /.s3cfg

# Set user.email so crosbymichael's in-container merge commits go smoothly
RUN	git config --global user.email 'docker-dummy@example.com'

# Add an unprivileged user to be used for tests which need it
RUN groupadd -r docker
RUN useradd --create-home --gid docker unprivilegeduser

VOLUME	/var/lib/docker
WORKDIR	/go/src/github.com/dotcloud/docker
ENV	DOCKER_BUILDTAGS	apparmor selinux

# Wrap all commands in the "docker-in-docker" script to allow nested containers
ENTRYPOINT	["hack/dind"]

# Upload docker source
ADD	.	/go/src/github.com/dotcloud/docker
