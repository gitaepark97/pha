DB_URL=mysql://hugo:hugo_drowssap@tcp(localhost:3330)/pha?parseTime=true

createdb:
	docker run --name mysql -p 3330:3306 -e MYSQL_ROOT_PASSWORD=root_drowssap -e MYSQL_USER=hugo -e MYSQL_PASSWORD=hugo_drowssap -e MYSQL_DATABASE=pha -e TZ=Asia/seoul --platform linux/amd64 -d mysql:5.7 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

removedb:
	docker container stop mysql && docker container rm mysql

migrateup:
	migrate -path migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db/database.dbml --project pha

db_dbml:
	sql2dbml doc/db/schema.sql --mysql -o doc/db/database.dbml

sqlc:
	sqlc generate

server:
	go run main.go

mock:
	mockgen -package mockrepository -destination repository/mock/repository.go github.com/gitaepark/pha/repository Repository
	mockgen -package mockservice -destination service/mock/service.go github.com/gitaepark/pha/service Service

test:
	go test -v -cover ./...

.PHONY: mysql migrateup migrateup1 migratedown migratedown1 db_docs db_dbml sqlc server mock test