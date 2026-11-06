# Secrets in Vault

The API keeps its whole configuration (secrets and settings) in **one Vault
entry**, `secret/adfluence/prod`, whose keys are the usual variable names
(`JWT_SECRET`, `DATABASE_URL`, `SMTP_PASSWORD`, `FRONTEND_URL`, …).

The server keeps only:

| Where | What |
|---|---|
| `.env` (`chmod 600`) | `VAULT_TOKEN`: a **read-only** token for that one entry |
| `docker-compose.override.yml` | `VAULT_ADDR=http://vault:8200`, `VAULT_PATH=secret/adfluence/prod` |

At startup the API reads the entry and applies every key before anything else.

- If Vault is unreachable or sealed, or the token is refused, the API **refuses to start** and logs why (names only, never values). Docker keeps restarting it, so it comes up by itself once Vault is back.
- The API renews its token every day.
- Local runs and CI don't set `VAULT_ADDR`, so they keep using `.env` exactly as before.

## Files

| File | Purpose |
|---|---|
| `vault.hcl` | Server config: file storage (encrypted), UI on, no TLS (traffic never leaves the server) |
| `adfluence-api-policy.hcl` | The API's token may **read** `secret/data/adfluence/prod`, nothing else |
| `env-to-vault.py` | Turns a `.env` file into JSON, piped straight into `vault kv put` (values never touch disk or shell history) |

## Safety rules

- **Unseal keys and root token are kept offline** (password manager), never on the server or in chat.
- **Vault is never published to the internet.** It listens on the Docker network, plus `127.0.0.1:8200` on the server for the UI over SSH:
  `ssh -i adfluence-key.pem -L 8200:127.0.0.1:8200 ubuntu@<server>`, then open http://localhost:8200.
- **Use the root token only for setup**, then revoke it; generate a new one with the unseal keys when needed (`vault operator generate-root`).
- **The audit log** (`/vault/logs/audit.log`) records every access.

## After a server reboot (manual unseal)

Vault starts **sealed** and the API waits. Unseal with 3 of the 5 keys:

```bash
docker exec -it adfluence_vault vault operator unseal   # paste key 1 (hidden)
docker exec -it adfluence_vault vault operator unseal   # key 2
docker exec -it adfluence_vault vault operator unseal   # key 3
```

The API starts by itself within a few seconds.

## Changing a value

Use the UI (over the SSH tunnel), or the command line, which keeps the value out of history:

```bash
read -rs -p "Root token: " VAULT_TOKEN && export VAULT_TOKEN && echo
read -rs -p "New SMTP_PASSWORD: " V && echo
printf '%s' "$V" | docker exec -i -e VAULT_TOKEN adfluence_vault vault kv patch -mount=secret adfluence/prod SMTP_PASSWORD=-; unset V
docker compose restart backend
```

Use `printf '%s'`, not `<<<`: `<<<` adds a newline to the stored value. `kv patch` changes only that key and keeps the others.

The API reads Vault only at startup, so **restart the backend after a change**. Every change is a new version: `vault kv rollback` undoes one.
