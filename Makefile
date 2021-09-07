DBNAME:=todolist
# https://docs.docker.com/docker-for-mac/networking/#use-cases-and-workarounds
DOCKER_DNS:=db
FLYWAY_CONF?=-url=jdbc:mysql://$(DOCKER_DNS):3306/$(DBNAME) -user=root -password=password

export DATABASE_DATASOURCE:=root:password@tcp($(DOCKER_DNS):3306)/$(DBNAME)

docker-compose/build:
	docker-compose build

docker-compose/up:
	docker-compose up

docker-compose/upd:
	docker-compose up -d

docker-compose/up/service:
	docker-compose up $(service)

docker-compose/down:
	docker-compose down

docker-compose/logs:
	docker-compose logs -f

DB_SERVICE:=db
mysql/client:
	docker-compose exec $(DB_SERVICE) mysql -uroot -hlocalhost -ppassword $(DBNAME)

mysql/init:
	docker-compose exec $(DB_SERVICE) \
		mysql -u root -h localhost -ppassword \
		-e "create database \`$(DBNAME)\`"

__mysql/drop:
	docker-compose exec $(DB_SERVICE) \
		mysql -u root -h localhost -ppassword \
		-e "drop database \`$(DBNAME)\`"

# ローカルのDBを全部クリアして、マイグレーションを1から実行します
# -j オプションは禁止！
_mysql_initialize: __mysql/drop mysql/init flyway/migrate

MIGRATION_SERVICE:=migration
flyway/info:
	docker-compose run --rm $(MIGRATION_SERVICE) $(FLYWAY_CONF) info

flyway/validate:
	docker-compose run --rm $(MIGRATION_SERVICE) $(FLYWAY_CONF) validate

flyway/migrate:
	docker-compose run --rm $(MIGRATION_SERVICE) $(FLYWAY_CONF) migrate

flyway/repair:
	docker-compose run --rm $(MIGRATION_SERVICE) $(FLYWAY_CONF) repair

flyway/baseline:
	docker-compose run --rm $(MIGRATION_SERVICE) $(FLYWAY_CONF) baseline

# 以下 prod rds 用 flyway ( local で叩いても繋がらないよ
FLYWAY_DOCKER:=flyway/flyway
CODEBUILD_SRC_DIR?=$(shell pwd)
PROD_FLYWAY_CMD = \
  prod/flyway/info \
  prod/flyway/validate \
  prod/flyway/migrate \
  prod/flyway/repair \
  prod/flyway/baseline

$(PROD_FLYWAY_CMD):
	@echo run $(@F) target in prod
	@docker run -v $(CODEBUILD_SRC_DIR)/database/schema:/flyway/sql -i --rm $(FLYWAY_DOCKER) $(FLYWAY_CONF) $(@F)
