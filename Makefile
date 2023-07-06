# //TODO: Explain to Daniel how to use makefile and how we will implement the CI/CD using Makefile at server level

user := root
database := test

dev:
	docker-compose up -d
	@echo "docker development setup started."
	
dev-build:
	docker-compose up --build
	@echo "docker compose image rebuilded."

stop:
	@echo "Stopping and clear all"
	docker-compose down
	@echo "Docker compose Stopped"

swagger:
	docker-compose exec -i app swag init
	@echo "Swagger doc generated"