FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
RUN cd /app ; go install trendserver
CMD ["/app/bin/trendserver"]