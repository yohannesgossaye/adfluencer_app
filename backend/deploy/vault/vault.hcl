# Adfluence's Vault server (one container next to the API).
# Reachable only inside Docker's network, plus 127.0.0.1 on the server for the
# web UI through an SSH tunnel. Never publish it to the internet.

ui            = true
api_addr      = "http://127.0.0.1:8200"
disable_mlock = false # the container gets IPC_LOCK, so secrets never go to swap

storage "file" {
  path = "/vault/file" # encrypted by Vault; useless without the unseal keys
}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = true # traffic never leaves this server (Docker network / SSH tunnel)
}
