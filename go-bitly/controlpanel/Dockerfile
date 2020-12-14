FROM golang:latest 
EXPOSE 6060
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
RUN cd /app ; go install controlpanel
CMD ["/app/bin/controlpanel"]