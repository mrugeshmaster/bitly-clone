
all: clean

clean: 
	find . -name 'controlpanel' -type f -exec rm -f {} \; 
	go clean

go-get:
	rm -rf src/github.com
	go get -v github.com/streadway/amqp
	go get -v github.com/codegangsta/negroni
	go get -v github.com/gorilla/mux
	go get -v github.com/unrolled/render

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt controlpanel

install:
	go install controlpanel

build:
	go build controlpanel

start:
	./controlpanel 

docker-build: 
	docker build -t controlpanel .
	docker images

network-create:
	docker network create --driver bridge controlpanel

network-inspect:
	docker network inspect controlpanel

rabbitmq-run:
	docker run --name rabbitmq \
	           --network controlpanel \
			   -p 8080:15672 -p 4369:4369 -p 5672:5672 \
			   -d rabbitmq:3-management

docker-run:
	docker run \
            --network controlpanel \
			--name controlpanel -p 3000:3000 -td controlpanel
	docker ps

docker-network:
	docker network ls

docker-network-inspect:
	docker network inspect host

docker-shell:
	docker exec -it controlpanel bash 

docker-clean:
	docker stop rabbitmq
	docker rm rabbitmq	
	docker stop controlpanel
	docker rm controlpanel
	docker rmi controlpanel

docker-ip:
	docker-machine ip

docker-ps:
	 docker ps --all --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t"

docker-ps-ports:
	 docker ps --all --format "table {{.Names}}\t{{.Ports}}\t"

test-ping:
	curl localhost:6060/ping