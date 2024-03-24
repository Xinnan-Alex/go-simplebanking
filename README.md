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
   - Durability: successfully transaction data must be recorded in persistant storage.
9. Overcome deadlock by always updating the account with smaller ID first.
10. Learnt how to use Github Action to build and test code with CI.
11. Implemented REST API Endpoint using [Gin Framework](https://github.com/gin-gonic/gin), used built-in validator to validate JSON content in request body, query params and uri params.
12. Use [Golang Viper](https://github.com/spf13/viper) Library to load config using file/environmental variable
13. Use [GoMock](https://github.com/uber-go/mock) Library to create a mockdb for API Testing
14. Create a custom validator and register it into Gin Framework, currency validator
15. Create custom matcher for creating unit test for Create User API Endpoint using GoMock.
16. Learn the differences between JWT and PASETO and identify why PASETO is much more secured than using JWT.
17. Implement JWT/Paseto Token Maker using [JWT](https://github.com/golang-jwt/jwt) and [PASETO](https://github.com/o1egl/paseto) Library.
18. Login API will now create a JWT/Paseto token depending on which TokenMaker is initialised.
19. Implement an authentication middleware to authenticate API endpoints.
20. Implement unit tests for testing authentication middleware and authenticated API endpoints.
21. Create a Dockerfile Image to build and run the app - optimised with only copying the go binary file into the image
22. Learn to connect two seperated container using network with command `docker network create networkname` and connect them with command `docker network connect networkname containername`. You can also connect it when running a container by specifying in the docker run command `docker run --name containername --network networkname -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:password@postgrescontainername:5432/simple_bank?sslmode=disable" simplebank:latest`. By doing so, you don't have to specify the IP of the postgres container, you can just provide the container name.
23. Learn to use docker composes to spin up multiple containers at the same time (stack) and also modify the startup sequence of each container using [wait-for](https://github.com/eficode/wait-for) script.

## Diagrams

### Simple Banking System Database Schema

![Simple Banking System Database Schema](/docs/SimpleBank_Schema.png)
