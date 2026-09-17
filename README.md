# chirpy-boot-dev

A small Twitter-like REST API built in Go as part of the [boot.dev](https://boot.dev) backend course. Users can register, log in, post short "chirps", and browse/delete them.

## Setup

1. Copy `.env` with the following variables:
   ```
   DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
   SECRET=your-jwt-secret
   POLKA_KEY=your-polka-api-key
   ```
2. Run database migrations (see `sql/`).
3. Start the server:
   ```
   go run .
   ```

The server listens on `:8080`.

## API

### Auth
| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/users` | Create a user (`email`, `password`) |
| PUT | `/api/users` | Update the authenticated user's email/password |
| POST | `/api/login` | Log in, returns a JWT access token and a refresh token |
| POST | `/api/refresh` | Exchange a valid refresh token for a new JWT |
| POST | `/api/revoke` | Revoke a refresh token |

### Chirps
| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/chirps` | Create a chirp (requires auth) |
| GET | `/api/chirps` | List chirps, optional `?author_id=` filter and `?sort=asc\|desc` (by creation date) |
| GET | `/api/chirps/{chirpID}` | Get a single chirp |
| DELETE | `/api/chirps/{chirpID}` | Delete a chirp you authored (requires auth) |

Chirps must be 140 characters or fewer. The words `kerfuffle`, `sharbert`, and `fornax` are censored.

### Webhooks
| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/polka/webhooks` | Upgrades a user to Chirpy Red (requires `Authorization: ApiKey <POLKA_KEY>`) |

### Misc
| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/healthz` | Health check |
| GET | `/admin/metrics` | Number of times the app has been visited |
| POST | `/admin/reset` | Reset metrics and delete all users (dev only) |

## Auth

Authenticated endpoints require a JWT in the `Authorization: Bearer <token>` header, obtained from `/api/login`.
