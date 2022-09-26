init:
	mv .env.example .env && docker-compose up -d

down:
	docker-compose down
up:
	docker-compose up -d --remove-orphans