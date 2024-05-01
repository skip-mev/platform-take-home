# Remote Singing Service

A basic remote data signing service that uses the Vault Transit Engine as a secure key-store.

## Features

1. Create a Secp256r1 (ECDSA-P256) public/private key pair
2. Returning public keys created by the service
3. Signing bytes with the private key associated to the wallet

Note: The service is using in-memory data store. Also it's relying on locally ran Vault with transit enabled. Not suitable for production.

## Run

### Provision

Terraform provisioning:

```
make tf-init
```

### Build

Build the binary:
```
make build
```

### Run Vault

Install Vault first, then
```
start-vault
```
And enable transit:
```
vault secrets enable -path=transit transit
``` 

### Setup environment

Add these to your .env in root:
```
VAULT_ADDR=http://127.0.0.1:8200
VAULT_TOKEN=[token]
```
Get the `token` from the console when you run the Vault start command.


### Rum the service

```
make start
```

### Graceful shutdown

```
make stop
```

## Test

```
make test
```

Notes: 
- Creating the same wallet twice will get an error response because wallet name is unique. You need to restart the service after every test run to flush the in-memory store since the test doesn't do any data clean up after it's done, and the service does't allow a delete policy.
- Sometimes `make test` fails some tests although it often works. TODO: debug.
- You can run the tests manually one by one inside the test file. They will pass.s