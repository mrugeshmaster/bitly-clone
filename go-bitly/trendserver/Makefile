
all: clean

clean: 
	find . -name 'trendserver' -type f -exec rm -f {} \; 
	go clean

go-get:
	rm -rf src/github.com
	go get -v github.com/streadway/amqp

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt trendserver

install:
	go install trendserver

build:
	go build trendserver

start:
	./trendserver 

docker-build: 
	docker build -t trendserver .
	docker images

network-create:
	docker network create --driver bridge trendserver

network-inspect:
	docker network inspect trendserver

rabbitmq-run:
	docker run --name rabbitmq \
	           --network trendserver \
			   -p 8080:15672 -p 4369:4369 -p 5672:5672 \
			   -d rabbitmq:3-management

docker-run:
	docker run \
            --network trendserver \
			--name trendserver -p 3000:3000 -td trendserver
	docker ps

docker-network:
	docker network ls

docker-network-inspect:
	docker network inspect host

docker-shell:
	docker exec -it trendserver bash 

docker-clean:
	docker stop rabbitmq
	docker rm rabbitmq	
	docker stop trendserver
	docker rm trendserver
	docker rmi trendserver

docker-ip:
	docker-machine ip

docker-ps:
	 docker ps --all --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t"

docker-ps-ports:
	 docker ps --all --format "table {{.Names}}\t{{.Ports}}\t"

test-ping:
	curl localhost:6060/ping