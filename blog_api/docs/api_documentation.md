## Blog API Documentation

---

## Overview

Welcome to the Blog Platform API. This guide explains how to use each endpoint with clear descriptions, example request/response bodies, and common error codes.

---

**Base URL**: `http://localhost:8080` (or your configured server URL)  
**Base path**: `/api`

### Using Postman

- Create a Postman collection and add each request from this document
- Create a Postman Environment with variables:
  - `baseUrl`: `http://localhost:8080`
  - `accessToken`: set after login
  - `refreshToken`: set after login
- For protected routes, set the Authorization header to: `Bearer {{accessToken}}`

## Table of Contents

1. Authentication
2. User Management (register/login/logout/etc.)
3. Token Management
4. OAuth Integration
5. Admin Operations
6. Blogs
7. Comments
8. Models (reference)
9. Error Handling
10. Security Best Practices
11. Testing

---

## Authentication

Our API uses JWT (JSON Web Tokens) for authentication. You'll need to include your access token in the Authorization header for protected endpoints.

**Format**: `Authorization: Bearer <your_access_token>`

### Token Types

- **Access Token**: Valid for 15 minutes, used for API requests
- **Refresh Token**: Valid for 7 days, used to get new access tokens

---

## User Management

### 1. User Registration

Create a new user account on the platform.

**Endpoint**: `POST /api/users/register`

**Request Body**:

```json
{
  "username": "zufan_Gebrehiwot",
  "first_name": "Zufan",
  "last_name": "Gebrehiwot",
  "email": "zufan@example.com",
  "password": "securePassword123!"
}
```

**Response** (200 OK):

```json
{
  "message": "User registered successfully"
}
```

**Validation Rules**:

- Username: 3-50 characters, alphanumeric and underscores only
- Email: Valid email format, must be unique
- Password: Minimum 8 characters, must contain uppercase, lowercase, number, and special character
- First/Last name: 1-50 characters

**Error Responses**:

- `400 Bad Request`: Invalid input data or validation errors
- `409 Conflict`: Username or email already exists

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/users/register`
- Body: raw JSON (use the example above)

---

### 2. User Login

Authenticate a user and receive access tokens.

**Endpoint**: `POST /api/users/login`

**Request Body**:

```json
{
  "email_or_username": "zufan@example.com",
  "password": "securePassword123!"
}
```

**Response** (200 OK):

```json
{
  "message": "Login successful",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**Notes**:

- You can login with either email or username
- Access token expires in 15 minutes
- Store the refresh token securely for token renewal

**Error Responses**:

- `400 Bad Request`: Invalid input format
- `401 Unauthorized`: Invalid credentials

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/users/login`
- Body: raw JSON (use the example above)

---

### 3. User Logout

Logout a user and invalidate their tokens.

**Endpoint**: `POST /api/users/logout`

**Request Body**:

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response** (200 OK):

```json
{
  "message": "Logout successful"
}
```

**Notes**:

- This invalidates both access and refresh tokens
- User will need to login again to get new tokens

**Error Responses**:

- `400 Bad Request`: Invalid input format
- `401 Unauthorized`: Invalid or expired token

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/users/logout`
- Body: raw JSON `{ "access_token": "{{accessToken}}" }`

---

### 4. Forgot Password

Request a password reset link via email.

**Endpoint**: `POST /api/users/forgot-password`

**Request Body**:

```json
{
  "email": "zufan@example.com"
}
```

**Response** (200 OK):

```json
{
  "message": "If the email exists, a password reset link has been sent"
}
```

**Notes**:

- Always returns success to prevent email enumeration
- If email exists, a reset link will be sent
- Reset tokens expire after 1 hour

**Error Responses**:

- `400 Bad Request`: Invalid email format

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/users/forgot-password`
- Body: raw JSON (use the example above)

---

### 5. Reset Password

Reset password using the token from forgot password email.

**Endpoint**: `POST /api/users/reset-password`

**Request Body**:

```json
{
  "token": "reset_token_from_email",
  "new_password": "newSecurePassword123!"
}
```

**Response** (200 OK):

```json
{
  "message": "Password reset successfully"
}
```

**Validation Rules**:

- New password: Minimum 8 characters, must contain uppercase, lowercase, number, and special character
- Token must be valid and not expired

**Error Responses**:

- `400 Bad Request`: Invalid token or password validation failed
- `401 Unauthorized`: Expired or invalid reset token

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/users/reset-password`
- Body: raw JSON (use the example above)

---

### 6. Update User Profile

Update user profile information (requires authentication).

**Endpoint**: `PUT /api/users/profile`

**Headers**: `Authorization: Bearer <access_token>`

**Request Body**:

```json
{
  "first_name": "Zufan",
  "last_name": "Smith",
  "bio": "Software developer passionate about clean code",
  "profile_picture": "https://example.com/avatar.jpg",
  "contact_info": "zufan.smith@example.com"
}
```

**Response** (200 OK):

```json
{
  "id": "user_id_here",
  "username": "zufan_gebrehiwot",
  "first_name": "Zufan",
  "last_name": "Smith",
  "email": "zufan@example.com",
  "bio": "Software developer passionate about clean code",
  "profile_picture": "https://example.com/avatar.jpg",
  "contact_info": "zufan.smith@example.com",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T14:45:00Z"
}
```

**Notes**:

- All fields are optional
- Email and password cannot be updated through this endpoint
- Profile picture must be a valid URL

**Error Responses**:

- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Invalid or missing token

Postman:

- Method: PUT
- URL: `{{baseUrl}}/api/users/profile`
- Headers: `Authorization: Bearer {{accessToken}}`
- Body: raw JSON (use the example above)

---

## Token Management

### 1. Validate Access Token

Check if an access token is valid and get its claims.

**Endpoint**: `POST /api/auth/validate`

**Request Body**:

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response** (200 OK):

```json
{
  "valid": true,
  "claims": {
    "user_id": "user_id_here",
    "role": "user",
    "exp": 1705320000
  }
}
```

**Error Responses**:

- `400 Bad Request`: Invalid input format
- `401 Unauthorized`: Invalid or expired token

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/auth/validate`
- Body: raw JSON `{ "access_token": "{{accessToken}}" }`

---

### 2. Refresh Access Token

Get a new access token using a refresh token.

**Endpoint**: `POST /api/auth/refresh`

**Request Body**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response** (200 OK):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**Notes**:

- Refresh tokens are valid for 7 days
- Old access token remains valid until it expires
- Use this endpoint when access token is about to expire

**Error Responses**:

- `400 Bad Request`: Invalid input format
- `401 Unauthorized`: Invalid or expired refresh token

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/auth/refresh`
- Body: raw JSON `{ "refresh_token": "{{refreshToken}}" }`

---

## OAuth Integration

Our platform supports OAuth login with Google, GitHub, and Facebook.

### 1. Initiate OAuth Flow

Start the OAuth authentication process.

**Endpoint**: `GET /api/auth/{provider}/login`

**Providers**: `google`, `github`, `facebook`

**Response** (200 OK):

```json
{
  "auth_url": "https://accounts.google.com/oauth/authorize?...",
  "provider": "google"
}
```

**Usage**:

1. Call this endpoint to get the authorization URL
2. Redirect user to the `auth_url`
3. User will be redirected back to your callback URL with an authorization code

Postman:

- Method: GET
- URL: `{{baseUrl}}/api/auth/github/login` (or other provider)

---

### 2. Handle OAuth Callback

Process the OAuth callback and authenticate the user.

**Endpoint**: `GET /api/auth/{provider}/callback?code={authorization_code}&state={state}`

**Response** (200 OK):

```json
{
  "message": "OAuth login successful",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900,
  "is_new_user": true,
  "user": {
    "id": "user_id_here",
    "username": "john_doe",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com"
  }
}
```

**Notes**:

- `is_new_user` indicates if this is the first time the user logged in with OAuth
- User receives access and refresh tokens just like regular login
- If user doesn't exist, a new account is created automatically

**Error Responses**:

- `400 Bad Request`: Invalid authorization code or OAuth error

Postman:

- Method: GET
- URL: `{{baseUrl}}/api/auth/{provider}/callback?code=AUTH_CODE&state=STATE`

---

### 3. Link OAuth Account

Link an OAuth provider to an existing user account (requires authentication).

**Endpoint**: `POST /api/auth/{provider}/link`

**Headers**: `Authorization: Bearer <access_token>`

**Request Body**:

```json
{
  "code": "authorization_code_from_oauth_provider"
}
```

**Response** (200 OK):

```json
{
  "message": "OAuth account linked successfully",
  "provider": "google"
}
```

**Notes**:

- User must be authenticated
- This allows users to link multiple OAuth providers to their account
- Useful for account recovery and convenience

**Error Responses**:

- `400 Bad Request`: Invalid authorization code or OAuth error
- `401 Unauthorized`: Invalid or missing token

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/auth/{provider}/link`
- Headers: `Authorization: Bearer {{accessToken}}`
- Body: raw JSON `{ "code": "AUTH_CODE_FROM_PROVIDER" }`

---

## Admin Operations

These endpoints are only available to users with admin role.

### 1. Promote User to Admin

Promote a regular user to admin role.

**Endpoint**: `POST /api/admin/users/{userID}/promote`

**Headers**: `Authorization: Bearer <admin_access_token>`

**Path Parameters**:

- `userID`: ID of the user to promote

**Response** (200 OK):

```json
{
  "message": "User promoted to admin successfully",
  "userID": "target_user_id",
  "newRole": "admin"
}
```

**Business Rules**:

- Only admins can promote users
- Admin cannot promote themselves
- Target user must exist and be active
- User cannot be promoted if already an admin

**Error Responses**:

- `400 Bad Request`: Business rule violations
- `401 Unauthorized`: Invalid or missing token
- `403 Forbidden`: User is not an admin

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/admin/users/{userID}/promote`
- Headers: `Authorization: Bearer {{accessToken}}` (admin)

---

### 2. Demote Admin to User

Demote an admin back to regular user role.

**Endpoint**: `POST /api/admin/users/{userID}/demote`

**Headers**: `Authorization: Bearer <admin_access_token>`

**Path Parameters**:

- `userID`: ID of the admin to demote

**Response** (200 OK):

```json
{
  "message": "User demoted to user successfully",
  "userID": "target_user_id",
  "newRole": "user"
}
```

**Business Rules**:

- Only admins can demote other admins
- Admin cannot demote themselves
- Cannot demote the last admin in the system
- Target user must be an admin

**Error Responses**:

- `400 Bad Request`: Business rule violations
- `401 Unauthorized`: Invalid or missing token
- `403 Forbidden`: User is not an admin

Postman:

- Method: POST
- URL: `{{baseUrl}}/api/admin/users/{userID}/demote`
- Headers: `Authorization: Bearer {{accessToken}}` (admin)

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
