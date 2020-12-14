
all: clean

clean: 
	find . -name 'database' -type f -exec rm -f {} \; 
	go clean

go-get:
	rm -rf src/github.com
	go get -v github.com/streadway/amqp
	go get -v github.com/go-sql-driver/mysql

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt database

install:
	go install database

build:
	go build database

start:
	./database 

docker-build: 
	docker build -t database .
	docker images

network-create:
	docker network create --driver bridge database

network-inspect:
	docker network inspect database

rabbitmq-run:
	docker run --name rabbitmq \
	           --network database \
			   -p 8080:15672 -p 4369:4369 -p 5672:5672 \
			   -d rabbitmq:3-management

docker-run:
	docker run \
            --network database \
			--name database -p 3000:3000 -td database
	docker ps

docker-network:
	docker network ls

docker-network-inspect:
	docker network inspect host

docker-shell:
	docker exec -it database bash 

docker-clean:
	docker stop rabbitmq
	docker rm rabbitmq	
	docker stop database
	docker rm database
	docker rmi database

docker-ip:
	docker-machine ip

docker-ps:
	 docker ps --all --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t"

docker-ps-ports:
	 docker ps --all --format "table {{.Names}}\t{{.Ports}}\t"