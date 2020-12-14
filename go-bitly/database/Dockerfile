FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
ENV GOBIN /app/bin
RUN cd /app/src/database ; go install
CMD ["/app/bin/database"]