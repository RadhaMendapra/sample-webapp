FROM golang:latest

# set working directory in the container
WORKDIR $GOPATH/src/github.com/RadhaMendapra/sample-webapp

# copy from current working directory to the present working directory
COPY . .

# download dependencies
RUN go get -d -v ./...

# install packages
RUN go install github.com/go-chi/chi
RUN go install github.com/RadhaMendapra/sample-webapp

# expose port outside the container
EXPOSE 8080

# run sample-webapp executable
CMD ["sample-webapp"]
