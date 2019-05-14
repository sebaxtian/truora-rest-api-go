# truora-rest-api-go
API REST usando Go creado para resolver technical exercise propuesto por Truora https://www.truora.com/

#### Config CockroachDB

> WARNING: Use this config only at develop environment

##### Use Docker
[https://www.cockroachlabs.com/docs/stable/install-cockroachdb-linux.html#use-docker]
[https://www.cockroachlabs.com/docs/stable/start-a-local-cluster-in-docker.html#os-linux]

  * > docker pull cockroachdb/cockroach:v19.1.0

##### Run Container
> docker run --name docker-cockroachdb -p 26257:26257 -p 8080:8080 -d cockroachdb/cockroach:v19.1.0 start --insecure

##### Open SQL Shell
> docker exec -it docker-cockroachdb ./cockroach sql --insecure

##### Create User and Database

* In SQL Shell execute

  > CREATE USER IF NOT EXISTS truora;
  > CREATE DATABASE truora;
  > GRANT ALL ON DATABASE truora TO truora;
