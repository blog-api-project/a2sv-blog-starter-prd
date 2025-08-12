## Blog API Documentation

---

## Overview

Welcome to the Blog Platform API. This guide explains how to use each endpoint with clear descriptions, example request/response bodies, and common error codes.

---

**Base URL**: `http://localhost:8080` (or your configured server URL)  
**Base path**: `/api`

All protected endpoints require JWT authentication.

- Header: `Authorization: Bearer <token>`

## Authentication (JWT)

Our API uses JWT (JSON Web Tokens) for authentication. Include your access token in the `Authorization` header for protected endpoints.

**Format**: `Authorization: Bearer <your_access_token>`

**Token types**

- **Access Token**: Valid for ~15 minutes (used for API requests)
- **Refresh Token**: Valid for ~7 days (used to obtain new access tokens)

---

## Table of Contents

1. User Management (register/login/logout/etc.)
2. Token Management
3. OAuth Integration
4. Admin Operations
5. Blogs
6. Comments
7. Models (reference)
8. Error Handling
9. Security Best Practices
10. Testing

## User Management (summary)

### Register

- `POST /api/users/register`
- Body: `username`, `first_name`, `last_name`, `email`, `password`
- 200 Response:

```json
{ "message": "User registered successfully" }
```

### Login

- `POST /api/users/login`
- Body: `email_or_username`, `password`
- 200 Response:

```json
{
  "message": "Login successful",
  "access_token": "<token>",
  "refresh_token": "<token>",
  "token_type": "Bearer",
  "expires_in": 900
}
```

### Logout

- `POST /api/users/logout`
- Body: `{ "access_token": "<token>" }`
- 200 Response:

```json
{ "message": "Logout successful" }
```

### Forgot Password / Reset Password

- `POST /api/users/forgot-password` — request password reset (returns success regardless to avoid enumeration)
- `POST /api/users/reset-password` — provide reset token + new password

### Update Profile

- `PUT /api/users/profile` (auth required)
- Fields: `first_name`, `last_name`, `bio`, `profile_picture`, `contact_info`
- 200: returns updated user object

---

## Token Management

### Validate Access Token

- `POST /api/auth/validate`
- Body: `{ "access_token": "<token>" }`
- Response: token validity and claims

### Refresh Access Token

- `POST /api/auth/refresh`
- Body: `{ "refresh_token": "<token>" }`
- Response: new access token (and expiry)

Notes:

- Refresh tokens valid for ~7 days
- Use the refresh endpoint when access token is near expiry

---

## OAuth Integration

Supports OAuth login with `google`, `github`, `facebook`.

### Initiate OAuth Flow

- `GET /api/auth/{provider}/login`
- Providers: `google`, `github`, `facebook`
- Response includes `auth_url` to redirect user to.

### Handle OAuth Callback

- `GET /api/auth/{provider}/callback?code={authorization_code}&state={state}`
- Response contains access and refresh tokens and user info

### Link OAuth Account

- `POST /api/auth/{provider}/link` (auth required)
- Body: `{ "code": "authorization_code_from_oauth_provider" }`

---

## Admin Operations

> Admin-only endpoints — require admin access token

### Promote User to Admin

- `POST /api/admin/users/{userID}/promote`
- 200 Response:

```json
{
  "message": "User promoted to admin successfully",
  "userID": "target_user_id",
  "newRole": "admin"
}
```

Business rules: only admins can promote, cannot promote self, target must exist, cannot promote if already admin.

### Demote Admin to User

- `POST /api/admin/users/{userID}/demote`
- 200 Response:

```json
{
  "message": "User demoted to user successfully",
  "userID": "target_user_id",
  "newRole": "user"
}
```

Business rules: only admins can demote others, cannot demote self, cannot demote last admin.

---

### Blogs

#### Create Blog

- Method: `POST`
- Path: `/api/blogs/create`
- Auth: required
- Content-Type: `multipart/form-data`
- Form fields:
  - `title` (string, required)
  - `content` (string, required)
  - `tags` (string[], optional; send multiple keys: `tags=go&tags=web`)
  - `images` (file[], optional; multiple files allowed; field name: `images`)
- 200 Response:

```json
{ "message": "Blog created successfully" }
```

- Errors: `400` / `500` with:

```json
{ "error": "..." }
```

---

#### List Blogs

- Method: `GET`
- Path: `/api/blogs/`
- Auth: required
- Query params (optional):
  - `page` (int, default 1)
  - `page_size` (int, default 10)
  - `sort_by` (string: `recent` (default), `popular`, `discussed`, `shared`, `oldest`)
  - `title`, `author`, `tags` (accepted; used by search only)
- 200 Response:

```json
{
  "blog": [
    /* Blog objects */
  ],
  "pagination": {
    "total_pages": 1,
    "current_page": 1,
    "total_posts": 1,
    "page_size": 10
  }
}
```

---

#### Update Blog

- Method: `PUT`
- Path: `/api/blogs/:id`
- Auth: required (must be the author)
- Content-Type: `multipart/form-data`
- Form fields:
  - `title` (string, required)
  - `content` (string, required)
  - `tags` (string[], optional)
- 200 Response:

```json
{
  "message": "Blog updated successfully",
  "updated_blog": {
    /* Blog */
  }
}
```

- 400 Responses include:

```json
{ "error": "unauthorized access: you are not permitted to update this blog" }
{ "error": "blog title must not be empty" }
{ "error": "blog content must not be empty" }
```

---

#### Delete Blog

- Method: `DELETE`
- Path: `/api/blogs/:id`
- Auth: required (must be the author)
- 200 Response:

```json
{ "message": "Blog deleted successfully" }
```

- Errors: `400` / `502` with:

```json
{ "error": "..." }
```

---

#### Search Blogs

- Method: `GET`
- Path: `/api/blogs/search`
- Auth: required
- Query params (optional):
  - `title` (string; regex match, case-insensitive)
  - `author` (string; username; resolved to user ID)
  - `page`, `page_size`, `sort_by`, `tags` accepted but only `title` and `author` are applied
- 200 Response:

```json
{
  "count": 2,
  "data": [
    /* Blog objects */
  ]
}
```

---

#### Like Blog

- Method: `POST`
- Path: `/api/blogs/:id/like`
- Auth: required
- 200 Response:

```json
{ "message": "Blog liked successfully" }
```

- 400 Response example:

```json
{ "error": "user has already liked this blog" }
```

---

#### Dislike Blog

- Method: `POST`
- Path: `/api/blogs/:id/dislike`
- Auth: required
- 200 Response:

```json
{ "message": "You dislike the blog" }
```

- 400 Response example:

```json
{ "error": "user has already disliked this blog" }
```

---

### Comments

#### Create Comment

- Method: `POST`
- Path: `/api/comments/create/:id`
- Description: Create a comment on blog with `id = :id`
- Auth: required
- Content-Type: `application/json`
- Body:

```json
{ "content": "Nice post!" }
```

- 200 Response:

```json
{ "message": "Comment created successfully" }
```

- 400 Response example:

```json
{ "error": "content cannot be empty" }
```

---

#### Update Comment

- Method: `PUT`
- Path: `/api/comments/:id`
- Auth: required
- Content-Type: `application/json`
- Body:

```json
{ "content": "Edited content" }
```

- 200 Response:

```json
{ "message": "Comment updated successfully" }
```

- 400 Response example:

```json
{ "error": "Editing comment failed" }
```

---

#### Delete Comment

- Method: `DELETE`
- Path: `/api/comments/:id`
- Auth: required
- 200 Response:

```json
{ "message": "Comment deleted successfully" }
```

- 400/502 Response example:

```json
{ "error": "comment deletion failed" }
```

---

### Models (reference)

#### Blog

```json
{
  "ID": "string",
  "AuthorID": "string",
  "Title": "string",
  "Content": "string",
  "ImageURL": ["string"],
  "Tags": ["string"],
  "PostedAt": "ISO datetime",
  "LikeCount": 0,
  "DislikeCount": 0,
  "CommentCount": 0,
  "ShareCount": 0,
  "AISuggestion": "string",
  "CreatedAt": "ISO datetime",
  "UpdatedAt": "ISO datetime"
}
```

#### Comment

```json
{
  "ID": "string",
  "BlogID": "string",
  "UserID": "string",
  "Content": "string",
  "CreatedAt": "ISO datetime",
  "UpdatedAt": "ISO datetime"
}
```

---

## Error Handling

### Common HTTP Status Codes

- `200 OK`
- `201 Created`
- `400 Bad Request`
- `401 Unauthorized`
- `403 Forbidden`
- `404 Not Found`
- `409 Conflict`
- `500 Internal Server Error`

### Error Response Format

```json
{
  "error": "Human-readable error message"
}
```

### Common Error Messages

- Authentication: `"Invalid credentials"`, `"Invalid token"`, `"Token not found or expired"`
- Validation: `"Invalid input"`, `"Email format is invalid"`, `"Password must be at least 8 characters"`
- Business: `"Admin cannot promote themselves"`, `"Only admins can promote users"`

---

## Rate Limiting

- Currently, the API does not implement rate limiting.

---

## Security Best Practices

For API consumers:

1. Store tokens securely (avoid localStorage / JS-accessible cookies)
2. Use HTTPS in production
3. Validate tokens on client and server
4. Handle token expiration and implement automatic refresh
5. Sanitize and validate inputs on the server

---

## Testing the API (Postman recommended)

Recommended flow:

1. Register a user
2. Login and store `access_token` and `refresh_token`
3. Call protected endpoints with `Authorization: Bearer <access_token>`
4. When access token expires, call refresh and update `access_token`
5. Test admin endpoints with an admin account

---

For additional support, refer to the project documentation or contact the development team.
