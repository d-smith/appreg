## Application Registration

Implement simple app registration services using event sourcing and 
CQRS


### Code Generation

In domain... 

<pre>
protoc --go_out=. *.proto
</pre>


### Build for docker

Building for docker is complicated by the use of the OCI driver to connect
to oracle. What this means is we have to run a docker container that extends
the golang image with the libraries and configuration to build natively 
for the container.

A Makefile has been supplied for this (Makefile.docker), which does 
the build. To run it:

<pre>
cd $GOPATH
docker run --rm -it -v "$PWD":/go -w /go/src/github.com/xtraclabs xtracdev/goora bash
cd appreg
make -f Makefile.docker
</pre>


Note that if you are sharing your say mac gopath with the docker container, the
go get commands run from the makefile will overwrite your native cgo built
stuff with the container native stuff, so you'll need to do a go get
back on the native side if you want to run stuff there.

After building the binary in the docker image, exit the shell then build
the image via make.