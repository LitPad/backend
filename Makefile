ifneq (,$(wildcard ./.env))
include .env
export 
ENV_FILE_PARAM = --env-file .env

endif

build:
	docker-compose up --build -d --remove-orphans

buildp:
	docker-compose -f docker-compose-prod.yml up --build -d --remove-orphans

up:
	docker-compose up -d

upp:
	docker-compose -f docker-compose-prod.yml up -d

down:
	docker-compose down

downp:
	docker-compose -f docker-compose-prod.yml down

logs:
	docker-compose logs

logsp:
	docker-compose -f docker-compose-prod.yml logs
	
test:
	go test -timeout 30m ./tests -v -count=1 -run TestAdminDashboard

swag:
	swag init --md .