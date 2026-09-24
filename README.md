# frosk

Frosk is a lightweight local password manager written in Go.

It provides a simple desktop interface for storing and managing credentials while keeping sensitive data encrypted in a local SQLite database. The application is built with [Gio](https://gioui.org/) and does not require an external service or remote database.

## Features

- Local credential storage using SQLite
- Master-password protected vault
- AES-256-GCM encryption
- Argon2id password-based key derivation
- Randomly generated encryption keys and nonces
- Encrypted usernames and passwords
- Searchable list of stored services
- Add, open and delete credential entries
- Cross-platform desktop application written in Go
- No cloud backend or external account required

## Encryption design

Frosk separates the master password from the key used to encrypt stored credentials.

When a vault is initialized:

1. A random salt is generated.
2. The master password is processed using **Argon2id**.
3. The Argon2id output is split into:
   - a value used to verify the master password,
   - a 256-bit key used to protect the vault encryption key.
4. A separate random 256-bit **user secret key** is generated.
5. The user secret key is encrypted using **AES-256-GCM** and stored in the local database.
6. Stored credentials are encrypted using the user secret key.

When the vault is unlocked, the encryption key is reconstructed from the supplied master password, the user secret key is decrypted, and that key is then used to decrypt individual credential entries.

This means the master password itself is not stored in plaintext and the key directly protecting credentials is randomly generated rather than being the password itself.

## Technology

- **Go**
- **Gio** — native desktop UI
- **SQLite**
- **Argon2id**
- **AES-256-GCM**

## Running

First, install Gio:

```bash
go install gioui.org/cmd/gogio@latest
```

Then run:

```bash
make run
```