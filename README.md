# Fut App

## Description

Aplicação que comunica-se e gerencia os dados da API em que é comunicado(football-data.org). Funciona como gerenciador de autenticação limitando a visualização por endpoint.

### Como rodar a aplicação?

Primeiramente precisamos subir a dependência que é o postgres e o client para o baco que é o OminiDB.
- Executar o comando `make up-database`
- Utilizando o VSCode, ir na parte de Run and Debug e alterar o projeto para migrator e executar
- Executar o curl de criação de usuário:
`curl --location 'localhost:8000/auth/create' \
--header 'Content-Type: application/json' \
--data '{
    "user": "admin",
    "password": "xpto123"
}'`

- Executar o cur de login:
`curl --location 'localhost:8000/auth/login' \
--header 'Content-Type: application/json' \
--data '{
    "user": "usuario-1",
    "password": "xpto123"
}'`

- Pegar o JWT Token e colocar após o Bearer espaço e executar o curl:
`curl --location '127.0.0.1:8000/campeonatos/' \
--header 'Authorization: Bearer *****'`


### Frameworks

> - Echo: HTTP Handler
> - Gorm: ORM for Database
> - Postgres: Database
> - Netflix/Go env: Load app environments
> - Golang JWT: Json Web Token encryption
> - Stretchr/Testify: For tests

### Architecture

> - Standard Go Project Layout / Clean Architecture Simplified

## OminiDB - SQL Client

> [!NOTE]
> - User: admin
> - Password: admin