
# CMS API Multi-tenant (Go)

API backend completa para CMS simples com isolamento por tenant, autenticação JWT, CRUD de conteúdo com `jsonb` e upload de imagens.

## Estrutura

```
/cmd/api/main.go
/internal/
  /config
  /database
  /models
  /repositories
  /services
  /handlers
  /middlewares
  /routes
/database/migrations
/docs
/scripts
/pkg/
```

## Requisitos

- Go 1.22+
- PostgreSQL (local via Docker Compose ou externo)
- Variáveis de ambiente configuradas

- `psql` instalado para executar migrations SQL


## Variáveis

- `DB_HOST`
- `DB_PORT` (use `5435` quando subir o Postgres local com Docker Compose)
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `JWT_SECRET`
- `PORT` (opcional, padrão 8080)

## Migrations (criação do banco + tabelas)

As migrations SQL estão em `database/migrations`:

- `000001_create_database.up.sql`
- `000002_create_tables.up.sql`

Para executar tudo:

```bash
cp .env.example .env
source .env

# sobe somente o PostgreSQL local na porta 5435
docker compose up -d

# executa migrations (cria database + tabelas)
./scripts/migrate.sh up
```

Para rollback:

```bash
./scripts/migrate.sh down
```

## Rodando a API localmente

```bash
go run ./cmd/api
```

## Swagger

Com a API em execução:

- UI: `http://localhost:8080/swagger`
- OpenAPI JSON: `http://localhost:8080/swagger/openapi.json`

## Rotas

### Auth
- `POST /auth/register`
- `POST /auth/login`

### Content (autenticado)
- `GET /content?page=home` (buscar por página)
- `GET /content` (listar conteúdos do tenant)
- `POST /content` (criar/atualizar por página)
- `PUT /content` (criar/atualizar por página)
- `DELETE /content/:id`

### Upload (autenticado)
- `POST /upload` (`multipart/form-data`, campo `file`)

## Segurança implementada

- JWT com `tenant_id` no token
- Isolamento por tenant em middleware e queries
- Senha com bcrypt
- Validação de payload (`binding` Gin)
- Sanitização de strings
- Rate limit por IP
- CORS configurado
- Headers de segurança
- Mensagens de erro genéricas
