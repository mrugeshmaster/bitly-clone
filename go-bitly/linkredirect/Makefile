
all: clean

clean: 
	find . -name 'linkredirect' -type f -exec rm -f {} \; 
	go clean

go-get:
	rm -rf src/github.com
	go get -v github.com/streadway/amqp
	go get -v github.com/codegangsta/negroni
	go get -v github.com/gorilla/mux
	go get -v github.com/unrolled/render
	go get -v github.com/go-sql-driver/mysql

go-bin:
	go env -w GOBIN=$(shell pwd)/bin

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt linkredirect

install:
	go install linkredirect

build:
	go build linkredirect

start:
	./linkredirect 

docker-build: 
	docker build -t linkredirect .
	docker images

network-create:
	docker network create --driver bridge linkredirect

network-inspect:
	docker network inspect linkredirect

rabbitmq-run:
	docker run --name rabbitmq \
	           --network linkredirect \
			   -p 8080:15672 -p 4369:4369 -p 5672:5672 \
			   -d rabbitmq:3-management
mysql-run:
	docker run --name mysql \
			   --network linkredirect \ 
			   -p 3306:3306 -e MYSQL_ROOT_PASSWORD=bitly \
			   -td mysql:5.7

mysql-shell:
	docker exec -it mysql bash 

docker-run:
	docker run \
            --network linkredirect \
			--name linkredirect -p 3000:3000 -td linkredirect
	docker ps

docker-network:
	docker network ls

docker-network-inspect:
	docker network inspect host

docker-shell:
	docker exec -it linkredirect bash 

docker-clean:
	docker stop rabbitmq
	docker rm rabbitmq
	docker stop mysql
	docker rm mysql
	docker stop linkredirect
	docker rm linkredirect
	docker rmi linkredirect

docker-ip:
	docker-machine ip

docker-ps:
	 docker ps --all --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t"

docker-ps-ports:
	 docker ps --all --format "table {{.Names}}\t{{.Ports}}\t"

test-ping:
	curl localhost:7070/ping
