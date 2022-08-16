# CockroachDB (CRDB) Evaluation
This project follows the [Quickstart with CockroachDB](https://www.cockroachlabs.com/docs/cockroachcloud/quickstart.html) guide. Shared credentials are for example purposes. 

## Download CRDB
Enables the ability to execute CLI commands using [Cockroach Commands](https://www.cockroachlabs.com/docs/stable/cockroach-commands.html).

We are also able to run a local CRDB server. This is not covered within the project.

```
curl https://binaries.cockroachdb.com/cockroach-v22.1.5.linux-amd64.tgz | tar -xz && sudo cp -i cockroach-v22.1.5.linux-amd64/cockroach /usr/local/bin/
```

## Create DB Table using local file
```
cat dbinit.sql | cockroach sql --url postgresql://wonkymic:AEOLKPUjsaDjYS2UYDe60w@free-tier4.aws-us-west-2.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dmultipass-43
postgresql://wonkymic:AEOLKPUjsaDjYS2UYDe60w@free-tier4.aws-us-west-2.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dmultipass-43
```

## Connection Properties
| | |
| :--- | :--- |
| host | free-tier4.aws-us-west-2.cockroachlabs.cloud:26257 |
| database | defaultdb |
| user | wonkymic | 
| password | `<password>` |

### Application connection
postgresql://wonkymic:`<password>`@
free-tier4.aws-us-west-2.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dmultipass-43

### Command connection to create Account table
cat dbinit.sql | cockroach sql --url "postgresql://wonkymic:AEOLKPUjsaDjYS2UYDe60w@free-tier4.aws-us-west-2.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dmultipass-43"

## User Schema
```
{
    Id <uuid>,
    Name <string>
}
```

## API Validation
| endpoint | verb | curl | 
| :--- | :--- | :--- |
| ping | GET | curl http://localhost:8080/ping |
| user | GET | curl http://localhost:8080/user/1 |
| user | POST | curl http://localhost:8080/user/ |