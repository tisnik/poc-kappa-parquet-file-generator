.PHONY: default clean build

DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres

create-db:
	createdb -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) d0
	psql -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) -d d0 -f create_table.sql

drop-db:
	dropdb -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) d0

topic-cleaner/topic-cleaner:	topic-cleaner/topic-cleaner.go
	cd topic-cleaner && go build topic-cleaner.go

topic-constructor/topic-constructor:	topic-constructor/topic-constructor.go
	cd topic-constructor && go build topic-constructor.go

message-producer/message-producer:	message-producer/message-producer.go
	cd message-producer && go build message-producer.go

db-writer/db-writer:	db-writer/db-writer.go
	cd db-writer && go build db-writer.go

build:	topic-cleaner/topic-cleaner topic-constructor/topic-constructor message-producer/message-producer db-writer/db-writer

clean:
	rm topic-cleaner/topic-cleaner
	rm topic-constructor/topic-constructor
	rm message-producer/message-producer
	rm db-writer/db-writer
