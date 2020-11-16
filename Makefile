#
# Copyright © 2020 Pavel Tisnovsky
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

.PHONY: default clean build

DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres

default: help

create-db: ## Create database and initialize table(s)
	createdb -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) d0
	psql -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) -d d0 -f database/create_table.sql

drop-db: ## Drop database completely
	dropdb -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) d0

topic-cleaner/topic-cleaner: topic-cleaner/topic-cleaner.go  ## Build topic-cleaner tool
	cd topic-cleaner && go build topic-cleaner.go

topic-constructor/topic-constructor: topic-constructor/topic-constructor.go  ## Build topic-constructor tool
	cd topic-constructor && go build topic-constructor.go

message-producer/message-producer: message-producer/message-producer.go  ## Build message producer tool
	cd message-producer && go build message-producer.go

db-reader/db-reader: db-reader/db-reader.go  ## Build DB-reader tool
	cd db-reader && go build db-reader.go

db-writer/db-writer: db-writer/db-writer.go  ## Build DB-writer tool
	cd db-writer && go build db-writer.go

parquet-generator/parquet-generator: parquet-generator/parquet-generator.go  ## Build parquet-generator tool
	cd parquet-generator && go build parquet-generator.go

parquet-reader/parquet-reader: parquet-reader/parquet-reader.go  ## Build parquet-reader tool
	cd parquet-reader && go build parquet-reader.go

build:	topic-cleaner/topic-cleaner \
	topic-constructor/topic-constructor \
	message-producer/message-producer \
	db-reader/db-reader \
	db-writer/db-writer \
	parquet-generator/parquet-generator \
	parquet-reader/parquet-reader

clean: ## Perform cleanup
	rm topic-cleaner/topic-cleaner
	rm topic-constructor/topic-constructor
	rm message-producer/message-producer
	rm db-reader/db-reader
	rm db-writer/db-writer
	rm parquet-generator/parquet-generator
	rm parquet-reader/parquet-reader

fmt: ## Run Go formatter over all source files
	gofmt -d db-reader/*.go
	gofmt -d db-writer/*.go
	gofmt -d topic-constructor/*.go
	gofmt -d topic-cleaner/*.go
	gofmt -d message-producer/*.go
	gofmt -d parquet-generator/*.go
	gofmt -d parquet-reader/*.go

linter: ## Run Go linter over all source files
	golint db-reader/...
	golint db-writer/...
	golint topic-constructor/...
	golint topic-cleaner/...
	golint message-producer/...
	golint parquet-generator/...
	golint parquet-reader/...

errcheck: ## Run Go error checker over all source files
	pushd db-reader; errcheck ./...; popd
	pushd db-writer; errcheck ./...; popd
	pushd topic-constructor; errcheck ./...; popd
	pushd topic-cleaner; errcheck ./...; popd
	pushd message-producer; errcheck ./...; popd
	pushd parquet-generator; errcheck ./...; popd
	pushd parquet-reader; errcheck ./...; popd

help: ## Show this help screen
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-/]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-36s\033[0m %s\n", $$1, $$2}'
	@echo ''
