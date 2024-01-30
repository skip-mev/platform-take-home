# Skip Platform Engineer - take home

## Background

> 🎯 **Goal:** Create a basic remote data signing service that uses the Vault Transit Engine as a secure key-store

You can fork this repo to use as a template for the take-home.

## Outline

Our products interact with dozens of blockchains, including signing and submitting transactions to them. We want to be assured that the backing material for these wallets (secp256k1/r1 private keys) do not get exposed publicly or even to Skip employees. 

This is possible to accomplish with a remote signing service. The remote signing service for the sake of this exercise needs to support three basic operations: 

1.  Create a private/public key pair and address on the Cosmos Hub
2. Query a public key
3. Sign transaction bytes

The service should use Vault Transit Engine to store the keys securely at rest. 

## Requirements

### Features

1. Create a Secp256r1 public/private key pair
2. Returning public keys created by the service
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

## Specific things we’re looking for

1. Familiarity with (or ability to quickly learn) infrastructure-as-code platform such as a Terraform or similar
2. Familiarity with (or ability to quickly learn) core components of our platform stack, including Vault, Make, etc…
3. Ability to handle complex problems that stretch across traditional backend services, blockchain transactions/queries, and cloud infrastructure
4. Capacity to design systems that are readable, extensible, and functional & adhere to reasonable software design principles

> 🆘 Please reach out to us on Telegram (@bpiv400, @magelinskaas) or email (barry@skip.money, zygis@skip.money) if you have any questions.

