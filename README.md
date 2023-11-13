# Simple Banking System using Golang

Things learnt during developing this system:

1. Design a database scheme and generate postgres SQL code using [dbdiagram](dbdiagram.io)
2. Install docker engine and docker desktop on my machine (Debian 12)
3. Use docker to run a PostgreSQL container and run it on port 5432 on host machine
4. Using [TablePlus](https://tableplus.com/) to run some SQL and check database information
5. Learn to do migration using [golang-migrate](https://github.com/golang-migrate/migrate)
6. Generate CRUD code using [SQLC](https://github.com/sqlc-dev/sqlc)
7. Write unit tests for CRUD operations with the help of [Testify](https://github.com/stretchr/testify)
8. Implement database transaction following the ACID Property
   - Atomicity: Either all operations complete successfully or transaction fails and database is unchanged.
   - Consistency: Database state must be valid after transaction.
   - Isolation: concurrent transaction must not affect each other.
   - Duratability: successfully transation data must be recorded in persistant storage.
9. Overcome deadlock by always updating the account with smaller ID first
10. Learnt how to use Github Action to build and test code with CI

## Diagrams

### Simple Banking System Database Schema

![Simple Banking System Database Schema](/docs/SimpleBank_Schema.png)
