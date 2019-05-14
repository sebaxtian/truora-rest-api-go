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

#### Run API REST
> go run main.go

#### Test API REST
> GET: http://localhost:3333/infoserver?domain=truora.com

* Example response
> {
    "ID": "8bd875b5-8d5c-4a2e-8989-cc5b72d201c2",
    "Servers": [
        {
            "IPAddress": "34.193.69.252",
            "Address": "410 Terry Ave N.",
            "SslGrade": "A",
            "Country": "US",
            "Owner": "Amazon Technologies Inc."
        },
        {
            "IPAddress": "34.193.204.92",
            "Address": "410 Terry Ave N.",
            "SslGrade": "A",
            "Country": "US",
            "Owner": "Amazon Technologies Inc."
        }
    ],
    "ServersChanged": false,
    "SslGrade": "A",
    "PreviousSslGrade": "A",
    "Logo": "https://uploads-ssl.webflow.com/5b559a554de48fbcb01fd277/5b97f0ac932c3291fa40d053_icon32.png",
    "IsDown": false,
    "LastUpdated": "2019-05-14T12:42:33-05:00"
}
