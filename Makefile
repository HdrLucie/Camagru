all: build

build:
	docker compose up --build --remove-orphans 

down:
	docker compose down -v

clear:
	@if [ -n "$$(docker ps -a -q)" ]; then \
		docker rm -f $$(docker ps -a -q); \
	else \
		echo "No containers to remove"; \
	fi

prune: down
	yes | docker system prune && yes | docker volume prune

re: prune build