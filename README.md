# Skip Platform Engineer - take home exercise

## Background

🎯 **Goal:** Create a basic remote signing service that uses the Vault Transit Engine as a secure key-store

## Outline

Our services depend on interacting with the chains, not only querying them, but also submitting transactions from Skip wallets. We want to be assured that the backing material for these wallets (private keys) do not get exposed publicly or even Skip employees. This is possible to accomplish with a remote signing service. The remote signing service for the sake of this exercise should be pretty basic - allowing to create wallets on the Cosmos Hub, return their public keys and sign transaction bytes.

## Requirements

### Features

1. Creating a wallet associated with a secp256r1 key
2. Returning wallets created by the service
3. Signing bytes with the private key associated to the wallet

### API structure

The API structure is defined in Protobuf in the given code template. You should not modify the API structure as the given tests depend on the API requests and response being constant.

### Vault security

There are requirements to how the local Vault instance should be configured. You can use whichever infrastructure-as-code tool you feel most comfortable with to do this. The requirements for the Vault configuration are as such:

- There should be an app role that can only interact with the transit engine
- The app role should be the only one (aside from the root user) that can create, read, update the private keys
- The private keys should not be exportable or deletable by the app role.
- The remote signing service should communicate with Vault authenticated as the app role.

### Running tests

To ensure your service passes basic functionality tests, you can run `make test` in the root directory of the service template. This will test your API for basic functionality such as creating a wallet, getting all wallets from the API and validating a signature signed by the service.

## Helpers

To help you with the boilerplate, there's already an API server with the correct typing set up in `api/server/api_server.go`. None of the methods are implemented correctly and you should do so according to the requirements.

To start the server you can run `make start` and in another window you can run `make test` to test the correctness of your API.

## Specific things we’re looking for


<aside>
🆘 Please reach out to us on Telegram (@bpiv400, @magelinskaas) or email (barry@skip.money, zygis@skip.money) if you have any questions.

</aside>

###
