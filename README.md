# 📝 Blog Platform API (Go + Gin + MongoDB)

Welcome! This is a backend API for a modern blog platform. It focuses on clean endpoints, security, and solid foundations for future features like search, tags, and AI-assisted content. Built with Go, Gin, and MongoDB. ✨

## 🚀 What You Get

- ✅ User management: register, login, logout, profile update
- 🔐 Security: bcrypt, JWT access/refresh tokens, RBAC
- 🧑‍💻 Admin: promote/demote users
- 🔗 OAuth2: Google, GitHub, Facebook
- 🧪 Tests: unit + integration
- 📄 Postman-friendly API docs in `docs/api_documentation.md`
- 📰 Blog: full CRUD, pagination, tags, search, filtering
- 📈 Popularity: views, likes/dislikes, comments count
- 🤖 AI: content suggestions and enhancements

## 🧰 Tech Stack

- Language/Framework: Go (Gin)
- Database: MongoDB
- Auth: JWT (access + refresh), OAuth2 (google/github/facebook)
- Email: SMTP (password reset)
- Testing: `testing` + `testify`

## 🎯 Goals (from the PRD)

- RESTful API with clear endpoints
- Blog CRUD (create, read, update, delete)
- Authentication & authorization (JWT, RBAC, OAuth2)
- Tags, filtering, search (extensible design)
- AI hooks for content suggestions (future-ready)
- Performance & scalability in mind

### Feature Matrix (status)

- Users: registration, login, forgot/reset password, logout, profile update — ✅ implemented
- Tokens: validate, refresh — ✅ implemented
- RBAC and Admin: promote/demote — ✅ implemented
- OAuth2 login (Google/GitHub/Facebook) — ✅ implemented
- Blog posts (CRUD) — ✅ implemented
- Tags, filtering, search — ✅ implemented
- Popularity tracking (likes/views/comments) — ✅ implemented
- AI content suggestions — ✅ implemented

## 🧱 Architecture at a Glance

- Delivery: controllers, router
- Domain: models, contracts (repositories, services, use cases)
- Usecases: business logic
- Infrastructure: JWT, OAuth, validation, email, middleware
- Repositories: MongoDB adapters

## 🔧 Setup (3 steps)

1. Clone

```bash
git clone <your-repo-url>
cd blog_api
go mod tidy
```

2. Configure env (create `.env`)

```env
# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB_NAME=blog_platform

# JWT
JWT_SECRET_KEY=change_me_dev_only_please_use_long_random

# Email (for password reset)
SMTP_HOST=
SMTP_PORT=
SMTP_USERNAME=
SMTP_PASSWORD=

# OAuth
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URI=http://localhost:8080/api/auth/google/callback

GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
GITHUB_REDIRECT_URI=http://localhost:8080/api/auth/github/callback

FACEBOOK_CLIENT_ID=
FACEBOOK_CLIENT_SECRET=
FACEBOOK_REDIRECT_URI=http://localhost:8080/api/auth/facebook/callback

# Server
PORT=8080
```

3. Run

```bash
go run Delivery/main.go
```

The API is available at http://localhost:8080

## 🧪 Run Tests

```bash
go test ./...
```

## 🔐 Security Highlights

- Passwords hashed with bcrypt (never plain text)
- JWT access (15m) and refresh (7d) tokens
- RBAC middleware checks roles from JWT
- OAuth2 least-privilege scopes

## 🔌 Key Endpoints (overview)

- Users: register, login, logout, forgot/reset password, update profile
- Tokens: validate, refresh
- Admin: promote, demote
- OAuth: login URL, callback, link account
- Blogs: create, list (with pagination/filters), get, update, delete
- Blog Interactions: like, dislike, metrics
- AI: suggest content

Full details: see `docs/api_documentation.md` ✅

## ✍️ Blog Module (implemented)

Endpoints:

- POST `/api/blogs` — Create blog (auth)
  - Body: { title, content, tags: [string], cover_image?, published_at? }
- GET `/api/blogs` — List blogs with pagination/filtering/sorting
  - Query: page, page_size, sort=recent|popular|commented, tags, author, q (search)
  - Response: items[], pagination meta
- GET `/api/blogs/:id` — Get single blog (increments view count)
- PUT `/api/blogs/:id` — Update blog (author-only)
- DELETE `/api/blogs/:id` — Delete blog (author or admin)

Permissions:

- Create/Update/Delete: author; Admin can delete any
- Read: public

## 🔎 Search, Filter, Tags (implemented)

- Full-text search on title/content via `q` query param
- Filters: `tags` (comma-separated), `author`, `date_from`, `date_to`
- Sorting: `sort=recent|popular|commented`
- Pagination: `page`, `page_size` on list endpoints

## 📈 Popularity Tracking (implemented)

- Views incremented on GET `/api/blogs/:id`
- Likes/dislikes: one reaction per user, toggle behavior
- Metrics endpoint: GET `/api/blogs/:id/metrics` returns views, likes, comments count

## 🤖 AI Integration (implemented)

- POST `/api/ai/suggest`
  - Body: { title?, keywords?: [string], draft_content?: string }
  - Response: { suggestions: [string], outline?: [string], improvements?: [string] }
- Provider is abstracted; configure via environment or adapter

## 📐 Non-Functional Requirements

- Scalability: goroutines + channels for background tasks (email, metrics)
- Security: bcrypt; signed JWT; RBAC; OAuth2; careful error messages
- Performance: pagination on list endpoints; indexes on frequent queries (author, tags, created_at); consider caching hot reads
- Reliability: timeouts on external calls; clear error handling; logs

## 🧭 Development Tips

- Use Postman with an environment: `baseUrl`, `accessToken`, `refreshToken`
- Start with register → login → validate → profile update → refresh
- For admin endpoints, use an admin token

## 🗺️ Roadmap Checklist (PRD)

- [x] User registration, login, logout, profile update
- [x] JWT auth (access/refresh), validation/refresh flow
- [x] Forgot/reset password via email
- [x] RBAC + Admin promote/demote
- [x] OAuth2 login/link (Google/GitHub/Facebook)
- [x] Blog CRUD endpoints
- [x] Search, tags, filtering
- [x] Popularity tracking (views, likes, comments)
- [x] AI content suggestions

## 🧩 Extending to Full Blog Features

This starter is built to expand into:

- Blog CRUD (posts, comments)
- Search and tag-based filtering
- Popularity metrics (likes, views, comments)
- AI-assisted content suggestions

## 🛠️ Project Structure

```
blog_project/
├── Delivery/              # HTTP layer (controllers, router)
├── Domain/                # Core models and interfaces
├── Infrastructure/        # JWT, OAuth, validation, middleware
├── Repositories/          # MongoDB adapters
├── Usecases/              # Business logic
├── docs/                  # API documentation (Postman-ready)
└── .env             # Environment variables (local only)
```
