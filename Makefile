DB_HOST = localhost
DB_PORT = 5432
DB_USER = hlucie
DB_NAME = camagru

all: build

build:
	docker compose up --build --remove-orphans 

down:
	docker compose down -v

clean:
	@if [ -n "$$(docker ps -a -q)" ]; then \
		docker rm -f $$(docker ps -a -q); \
	else \
		echo "No containers to remove"; \
	fi

prune: down
	yes | docker system prune && yes | docker volume prune

re: prune build

db:
	docker exec -it postgres psql -U hlucie -d camagru
