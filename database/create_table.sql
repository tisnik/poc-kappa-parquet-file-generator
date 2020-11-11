CREATE TABLE reports (
    id                    SERIAL PRIMARY KEY,
    key                   CHAR(1) NOT NULL,
    cluster_id            CHAR(36) NOT NULL,
    path                  VARCHAR(200) NOT NULL,
    external_organization VARCHAR(20) NOT NULL,
    report                TEXT NOT NULL
);
