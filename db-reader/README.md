# PoC: Parquet file generator based on Kappa architecture

## DB reader

Tool to read all results from relational database. PostgreSQL is the default
variant of such database.

### Description

This tool reads all results stored in PostgreSQL database. It is possible to
use command line flags to select the database and its options (user, etc.).

### How to build the tool

```
go build
```

### Usage

```
Usage of db-reader:
  -db-host string
        database host name (default "localhost")
  -db-name string
        database name (default "d0")
  -db-password string
        database password for given user (default "postgres")
  -db-port int
        database port (default 5432)
  -db-user string
        database user (default "postgres")
  -show-clusterid
        ClusterID (default true)
  -show-id
        ID of the record (default true)
  -show-key
        Key of the record (default true)
  -show-org
        external organization (default true)
  -show-path
        path (default true)
  -show-report
        the whole report
```

### PostgreSQL table structure

Currently, the simplified SQL table structure is used:

```
CREATE TABLE reports (
    id                    SERIAL PRIMARY KEY,
    key                   CHAR(1) NOT NULL,
    cluster_id            CHAR(36) NOT NULL,
    path                  VARCHAR(200) NOT NULL,
    external_organization VARCHAR(20) NOT NULL,
    report                TEXT NOT NULL
);
```
