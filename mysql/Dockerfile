# FROM mysql:latest
FROM --platform=amd64 mysql:latest

COPY ./sqls/my.cnf ./etc/mysql/conf.d/my.cnf

COPY ./sqls/initialize.sql /docker-entrypoint-initdb.d/