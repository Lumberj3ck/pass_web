# pass_web

A web interface for the Unix [pass](https://www.passwordstore.org/) password manager. It lets you browse, search, decrypt, insert, and delete passwords from your password store through the browser, with all cryptographic operations happening client-side.

## Why

The standard `pass` CLI works well on a single machine, but accessing your passwords from another device or through a more visual interface requires extra tooling. Existing solutions often involve sending your private key or decrypted passwords to a server, which defeats the purpose of PGP-based encryption.

pass_web solves this by keeping all sensitive operations in the browser. Your private key and passphrase never leave the client. The server only stores encrypted password files and verifies your identity through PGP challenge-response authentication -- it never sees a plaintext password.

## Features

- **PGP challenge-response authentication** -- no server-side passwords. You prove your identity by signing a random challenge with your PGP private key.
- **Client-side decryption** -- passwords are decrypted entirely in the browser using OpenPGP.js. The server only serves encrypted blobs.
- **Client-side encryption on insert** -- new passwords are encrypted in the browser before being sent to the server.
- **Encrypted key storage** -- your private key and passphrase are encrypted with a master password (PBKDF2 + AES-GCM) and stored in localStorage.
- **Fuzzy search** -- quickly find passwords across the entire store with real-time fuzzy matching (Ctrl+K).
- **Folder navigation** -- browse the hierarchical directory structure of your password store.
- **Insert and delete** -- create new password entries and remove existing ones from the web interface.
- **Dark and light themes** -- toggle between themes, preference saved in localStorage.
- **JWT session management** -- authenticated sessions last 24 hours.

## Tech Stack

- **Backend:** Go, Gorilla Mux, gopenpgp, golang-jwt
- **Frontend:** HTMX, Alpine.js, OpenPGP.js, Tailwind CSS
- **Templates:** Go html/template (embedded in the binary)
- **Static assets:** Embedded via Go's `embed` package

## Prerequisites

- Go 1.23+
- Node.js and npm (only needed if modifying Tailwind CSS)
- A [pass](https://www.passwordstore.org/) password store (`~/.password-store/`)
- A PGP key pair (the public key must be accessible to the server)

## Installation

1. **Clone the repository**

   ```
   git clone <repo-url>
   cd pass_web
   ```

2. **Configure environment variables**

   ```
   cp .sample.env .env
   ```

   Edit `.env` with your values:

   | Variable              | Description                                     | Required |
   |-----------------------|-------------------------------------------------|----------|
   | `jwt_secret`          | Secret key used to sign JWT tokens               | Yes      |
   | `PASSWEB_PUB_KEY_PATH`| Absolute path to your PGP public key file        | Yes      |
   | `PREFIX`              | Path to your password store directory             | No (defaults to `~/.password-store/`) |

3. **Install Go dependencies**

   ```
   go mod download
   ```

4. **Build and run**

   ```
   go build -o pass_web
   ./pass_web -port 8080
   ```

   Or run directly without building:

   ```
   go run main.go -port 8080
   ```

5. **Open** `http://localhost:8080` in your browser.

### Development

For live reloading during development, install [Air](https://github.com/air-verse/air) and run:

```
air
```

If you need to modify styles, install the npm dependencies and run Tailwind:

```
npm install
npx @tailwindcss/cli -i static/css/input.css -o static/css/output.css --watch
```

A development Dockerfile is also available:

```
docker build -f Dockerfile.dev -t pass_web:dev .
docker run -p 8080:8080 -v $(pwd):/usr/local/app pass_web:dev
```

## How Authentication Works

1. You visit `/auth` and the server generates a random 20-character challenge.
2. You provide your PGP private key and passphrase in the browser.
3. The browser signs the challenge using OpenPGP.js and sends the signature to the server.
4. The server verifies the signature against your public key using gopenpgp.
5. On success, a JWT token is issued as an HttpOnly cookie (valid for 24 hours).
6. Challenges expire automatically after 2 minutes.
