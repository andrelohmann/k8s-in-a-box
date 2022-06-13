FROM golang:alpine

# Expose HTTP port and set necessary environment variables
EXPOSE 8000

# copy source code into the $GOPATH and switch to that directory
COPY . ${GOPATH}/src/github.com/andrelohmann/k8s-in-a-box
WORKDIR ${GOPATH}/src/github.com/andrelohmann/k8s-in-a-box

# compile source code and copy into $PATH
RUN go mod init; \
    go install;

# the default command runs the service in the foreground
CMD ["k8s-in-a-box"]
