DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres

create-db:
	createdb -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) d0
	psql -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) -d d0 -f create_table.sql

drop-db:
	dropdb -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) d0
