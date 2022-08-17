# CockroachDB (CRDB) Evaluation
This project follows the [Quickstart with CockroachDB](https://www.cockroachlabs.com/docs/cockroachcloud/quickstart.html) guide. Shared credentials are for example purposes. 

## Status
| Task | Status |
| :--- | :--- |
| CRDB Instance Created | Complete |
| GoLang & CRDB interaction | Complete |
| GoLang REST | Complete |
| REST & CRDB combination | Complete |
| CRDB Connection Pooling | Complete |
| Performance Testing | Incomplete |
| ACID (Isolation) Testing | Incomplete |

## Download CRDB
Enables the ability to execute CLI commands using [Cockroach Commands](https://www.cockroachlabs.com/docs/stable/cockroach-commands.html).

We are also able to run a local CRDB server. This is not covered within the project.

```
curl https://binaries.cockroachdb.com/cockroach-v22.1.5.linux-amd64.tgz | tar -xz && sudo cp -i cockroach-v22.1.5.linux-amd64/cockroach /usr/local/bin/
```

## Connection Properties
| | |
| :--- | :--- |
| host | free-tier4.aws-us-west-2.cockroachlabs.cloud:26257 |
| database | defaultdb |
| user | wonkymic | 
| password | `<password>` |

### Application connection
```
postgresql://wonkymic:<password>@
free-tier4.aws-us-west-2.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dmultipass-43
```

### Command connection to create Account table
```
cat dbinit.sql | cockroach sql --url "postgresql://wonkymic:<password>@free-tier4.aws-us-west-2.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full&options=--cluster%3Dmultipass-43"
```
## User Schema/Struct
```
{
    Id <uuid>,
    Name <string>
}
```

## API Validation
| endpoint | verb | url | Summary |
| :--- | :--- | :--- | :--- |
| ping | GET | `http://localhost:8080/ping` | Pong |
| user | GET | `http://localhost:8080/user/<uuid>` | Returns User Struct |
| user | POST | `http://localhost:8080/user` | Returns a list of Users |
