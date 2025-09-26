# RSS Aggregator API

Este repositório contém uma API REST em Go para agregar e expor posts de RSS/feeds.

## Descrição do projeto

Uma API simples para cadastramento de usuários, feeds (RSS) e seguimento de feeds por usuários. A aplicação também faz scraping periódico dos feeds cadastrados para popular uma tabela de posts.

Principais recursos:
- Criar usuário e gerar uma API key.
- Criar feeds e listar feeds públicos.
- Permitir que usuários sigam (follow) feeds.
- Listar posts dos feeds seguidos por um usuário.
- Endpoints protegidos via API Key enviada no header `Authorization: ApiKey {API_KEY}`.

## Tecnologias utilizadas

- Linguagem: Go
- Router: chi (github.com/go-chi/chi)
- CORS: go-chi/cors
- Dotenv: github.com/joho/godotenv
- Banco de dados: PostgreSQL (driver github.com/lib/pq)
- SQL builder/queries: sqlc (código gerado em `internal/database`)

## Pré-requisitos

- Go 1.20+ instalado
- PostgreSQL rodando
- make/psql (opcional, para rodar as migrations)

## Variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto (o `.env` já está listado em `.gitignore`). Exemplo mínimo:

PORT=8080
DB_URL=postgres://USER:PASSWORD@localhost:5432/rssagg?sslmode=disable

Observação: os nomes usados pela aplicação são `PORT` e `DB_URL`.

## Setup e execução (Windows / PowerShell)

1. Instale dependências (módulos Go serão resolvidos automaticamente ao compilar/rodar):

# Abra PowerShell na pasta do projeto
2. Exporte/coloque o arquivo `.env` na raiz com as variáveis acima.

3. Execute as migrations no banco (exemplo com psql):

# Crie o banco e aplique migrations (ajuste USER e DB conforme necessário)
# psql -U postgres -c "CREATE DATABASE rssagg;"
# psql -U postgres -d rssagg -f sql/schema/001_users.sql
# psql -U postgres -d rssagg -f sql/schema/002_users_api_key.sql
# psql -U postgres -d rssagg -f sql/schema/003_feeds.sql
# psql -U postgres -d rssagg -f sql/schema/004_feed_follows.sql
# psql -U postgres -d rssagg -f sql/schema/005_feeds_lastfetchedat.sql
# psql -U postgres -d rssagg -f sql/schema/006_posts.sql

4. Build e rode a aplicação:

# No PowerShell
go build -o rss-aggregator.exe
.
\rss-aggregator.exe

ou diretamente:

go run .

Ao iniciar, a aplicação irá carregar `.env`, conectar ao Postgres e iniciar um goroutine de scraping periódico.

## Endpoints principais

Base URL: http://localhost:{PORT}/v1

Endpoints implementados:

- GET /v1/healthz
  - Retorno: 200 OK ({}). Usado para readiness.

- POST /v1/users
  - Cria usuário. Payload JSON:
    {
      "name": "Seu Nome"
    }
  - Retorno: 201 Created
    {
      "id": "uuid",
      "created_at": "timestamp",
      "updated_at": "timestamp",
      "name": "Seu Nome",
      "api_key": "chave_gerada"
    }

- GET /v1/users
  - Protegido por API Key (ver seção Autenticação).
  - Retorna os dados do usuário associado à API Key.

- POST /v1/feeds
  - Protegido. Cria um feed associado ao usuário autenticado.
  - Payload JSON:
    {
      "name": "Nome do Feed",
      "url": "https://exemplo.com/feed"
    }
  - Retorno: 201 Created com objeto Feed.

- GET /v1/feeds
  - Lista todos os feeds.

- POST /v1/feed_follows
  - Protegido. Segue um feed (cria feed_follow).
  - Payload JSON:
    {
      "feedId": "uuid-do-feed"
    }
  - Retorno: 201 Created com objeto FeedFollow.

- GET /v1/feed_follows
  - Protegido. Lista os feeds que o usuário autenticado segue.

- DELETE /v1/feed_follows/{feedFollowID}
  - Protegido. Remove o follow do feed para o usuário autenticado.

- GET /v1/user_posts
  - Protegido. Retorna até 10 posts para o usuário autenticado (dos feeds que segue).

## Autenticação

A autenticação usa uma API Key retornada ao criar um usuário. Envie o header:

Authorization: ApiKey {API_KEY}

Exemplo (curl):

curl -X GET "http://localhost:8080/v1/users" -H "Authorization: ApiKey <SUA_API_KEY>"

## Exemplos (JSON)

1) Criar usuário

Request:
POST /v1/users
Content-Type: application/json

{
  "name": "João"
}

Response (201):

{
  "id": "8a7e3b6a-...",
  "created_at": "2025-09-26T12:34:56Z",
  "updated_at": "2025-09-26T12:34:56Z",
  "name": "João",
  "api_key": "abcdef0123456789"
}

2) Criar feed (autenticado)

Request:
POST /v1/feeds
Authorization: ApiKey {api_key}
Content-Type: application/json

{
  "name": "Blog do Exemplo",
  "url": "https://blog.exemplo.com/rss"
}

Response (201): objeto Feed com campos id, created_at, updated_at, closed_at, name, url, user_id

3) Seguir um feed

Request:
POST /v1/feed_follows
Authorization: ApiKey {api_key}
Content-Type: application/json

{
  "feedId": "<uuid-do-feed>"
}

Response (201): objeto FeedFollow

## Observações e próximos passos

- O projeto já inclui as queries geradas pelo `sqlc` em `internal/database`.
- Há um `startScraping` que roda em background (veja `scraper.go`) para manter a tabela de posts atualizada.
- Melhorias sugeridas: documentação OpenAPI/Swagger, testes automatizados, paginação e filtros nos endpoints de posts/feeds.

