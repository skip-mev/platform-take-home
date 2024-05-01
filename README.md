# Remote Signing Service

A basic remote data signing service utilizing the Vault Transit Engine as a secure key store.

## Features

1. Create a Secp256r1 (ECDSA-P256) public/private key pair.
2. Return public keys created by the service.
3. Sign bytes with the private key associated with the wallet.

Note: The service uses an in-memory data store and relies on a locally run Vault with transit enabled. It is not suitable for production use.

## Run

### Provision

Terraform provisioning:

```
make tf-init
```
The Terraform file 'main.tf' includes:
- Vault service configured for local development only.
- Resources for Vault including policy, backend auth, AppRole, etc.

### Build

Build the binary:
```
make build
```

### Run Vault

Start the Vault:
```
start-vault
```
Enable transit in Vault if it's not already enabled:
```
vault secrets enable -path=transit transit
``` 

### Setup Environment

Add the following to your `.env` in the root directory:
```
VAULT_ADDR=http://127.0.0.1:8200
VAULT_TOKEN=[token]
```
Obtain the `token` from the console output when you run the Vault start command.

### Run the Service

```
make start
```

### Graceful Shutdown

Stop the service:
```
make stop
```

## Test

Run the tests:
```
make test
```

Notes:
- Creating the same wallet twice will result in an error response because the wallet name is unique. You need to restart the service after every test run to flush the in-memory store since the test does not clean up data after completion, and the service does not allow deletion.
- Sometimes `make test` fails some tests although it often works. TODO: Debug.
- You can run the tests manually one by one inside the test file. They will pass.