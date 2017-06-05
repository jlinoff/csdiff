# This docker file creates a go compilation container that can be used
# to cross-compile go for linux on any platform that supports docker.
#
#   $ cd csdiff
#   $ docker run -it --rm -v $(pwd):/opt/go/project goco make
#
# To create it. Note that gover is an ARG that defines the version of
# go that you want to build. This example shows how to build it for go-1.8.3.
#
#   $ docker build --build-arg gover=1.8.3 -f Dockerfile -t goco:1.8.3 -t goco:latest .
FROM centos:latest

RUN yum clean all && yum update -y && yum install -y git make

ARG gover
ENV GO_VERSION=$gover
ENV GOROOT=/opt/go/latest
ENV GOPATH=/opt/go/project
ENV GO_PROG=/opt/go/latest/bin/go

# Setup the volume.
RUN mkdir -p ${GOPATH}
VOLUME ${GOPATH}

# Install go in /opt/go
RUN mkdir -p /opt/go/${GO_VERSION}/dl && \
    cd /opt/go/${GO_VERSION}/dl && \
    curl -k -O -L https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz && \
    cd /opt/go/${GO_VERSION} && \
    tar zxf dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/${GO_VERSION}/go /opt/go/latest && \
    ${GO_PROG} version
    
# Install golint
RUN cd /opt/go/${GO_VERSION}/dl && \
    GOPATH=/opt/go/${GO_VERSION}/dl ${GOROOT}/bin/go get -u github.com/golang/lint/golint && \
    cp bin/* ${GOROOT}/bin

# Wrapper for the go command that makes it
# natural for the user to run something like
# docker run -it --rm -v $(pwd):/opt/go/project goco go build myprog.go
RUN /bin/echo '#!/bin/bash'                           > /opt/go/goco.sh && \
    /bin/echo 'export PATH="${GOROOT}/bin:${PATH}"'  >> /opt/go/goco.sh && \
    /bin/echo 'cd /opt/go/project'                   >> /opt/go/goco.sh && \
    /bin/echo '$*'                                   >> /opt/go/goco.sh && \
    chmod a+rx /opt/go/goco.sh && \
    /opt/go/goco.sh go version

# Run in go environment.
ENTRYPOINT ["/opt/go/goco.sh"]
CMD ["/opt/go/latest/bin/go", "version"]
