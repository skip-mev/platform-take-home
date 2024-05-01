package server

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"

	"math/big"

	"github.com/hashicorp/vault/api"
	"github.com/skip-mev/platform-take-home/logging"
	"go.uber.org/zap"
	"golang.org/x/crypto/ripemd160"

	"github.com/skip-mev/platform-take-home/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type APIServerImpl struct {
	types.UnimplementedAPIServer

	logger      *zap.Logger
	vaultClient *api.Client
	walletNames []string // In memory slice to keep track of wallet names
}

var _ types.APIServer = &APIServerImpl{}

func NewDefaultAPIServer(vaultAddr string) *APIServerImpl {
	logger, err := logging.DefaultLogger()
	if err != nil {
		logger.Fatal(types.ErrMsgFailedToInitializeLogger, zap.Error(err))
		panic(types.ErrMsgFailedToInitializeLogger)
	}

	config := api.DefaultConfig()
	config.Address = vaultAddr
	client, err := api.NewClient(config)
	if err != nil {
		logger.Fatal(types.ErrMsgFailedToCreateVaultClient, zap.Error(err))
		panic(types.ErrMsgFailedToCreateVaultClient)
	}

	return &APIServerImpl{
		logger:      logger,
		vaultClient: client,
	}
}

func publicKeyToAddress(pubKey *ecdsa.PublicKey) (string, []byte) {
	if pubKey == nil {
		return "", nil
	}

	pubKeyBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	if pubKeyBytes == nil {
		return "", nil
	}

	// Hash the public key to compute the address
	hash := sha256.New()
	hash.Write(pubKeyBytes[1:]) // Remove the 0x04 prefix if present
	publicSHA256 := hash.Sum(nil)

	ripemd160Hasher := ripemd160.New()
	ripemd160Hasher.Write(publicSHA256)
	publicRIPEMD160 := ripemd160Hasher.Sum(nil)

	addressBytes := publicRIPEMD160[len(publicRIPEMD160)-20:] // Take the last 20 bytes

	// Convert to hex, skipping the "0x" prefix
	address := hex.EncodeToString(addressBytes)

	return address, addressBytes
}

func (s *APIServerImpl) CreateWallet(ctx context.Context, request *types.CreateWalletRequest) (*types.CreateWalletResponse, error) {
	logging.FromContext(ctx).Info("CreateWallet", zap.String("name", request.Name))

	if request.Name == "" {
		s.logger.Error("wallet name cannot be empty")
		return nil, status.Error(codes.InvalidArgument, "wallet name cannot be empty")
	}

	// Check if the wallet name already exists
	for _, name := range s.walletNames {
		if name == request.Name {
			return &types.CreateWalletResponse{
				Error: &types.Error{
					Message: "wallet name already exists",
				},
			}, nil
		}
	}

	// Create a new transit key in Vault using the wallet name from the request
	saved, err := s.vaultClient.Logical().Write(fmt.Sprintf("transit/keys/%s", request.Name), map[string]interface{}{
		"type":       "ecdsa-p256",
		"exportable": true,
	})
	if err != nil {
		s.logger.Error("failed to create Vault transit key", zap.Error(err))
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "failed to create wallet",
			},
		}, err
	}

	if saved.Data == nil {
		s.logger.Error("saving data failed")
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "no data in saved response",
			},
		}, err
	}

	pubKeyResponse, err := s.vaultClient.Logical().Read(fmt.Sprintf("transit/keys/%s", request.Name))
	if err != nil {
		s.logger.Error("failed to retrieve public key", zap.Error(err))
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "failed to retrieve public key",
			},
		}, err
	}

	if pubKeyResponse.Data == nil {
		s.logger.Error("no data in response")
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "no data in response",
			},
		}, err
	}

	keysInterface, ok := pubKeyResponse.Data["keys"].(map[string]interface{})
	if !ok {
		s.logger.Error("keys map not found in response")
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "keys map not found in response",
			},
		}, err
	}

	keyData, ok := keysInterface["1"].(map[string]interface{})
	if !ok {
		s.logger.Error("key version data not found in response")
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "key version data not found in response",
			},
		}, err
	}

	pubKeyInterface, ok := keyData["public_key"]
	if !ok {
		s.logger.Error("public key not found in response")
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "public key not found in response",
			},
		}, err
	}

	pubKey, ok := pubKeyInterface.(string)
	if !ok {
		s.logger.Error("public key is not a string")
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "public key format error",
			},
		}, err
	}

	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		s.logger.Error("failed to parse PEM block containing the public key")
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "PEM block parsing failed",
			},
		}, fmt.Errorf("PEM block parsing failed")
	}

	pubKeyBytes := block.Bytes

	s.walletNames = append(s.walletNames, request.Name)

	pubKeyHash := sha256.Sum256(pubKeyBytes)
	bech32Addr, err := sdk.Bech32ifyAddressBytes("cosmos", pubKeyHash[:])
	if err != nil {
		s.logger.Error("failed to convert address to Bech32", zap.Error(err))
		return &types.CreateWalletResponse{
			Error: &types.Error{
				Message: "failed to convert address to Bech32",
			},
		}, err
	}

	newWallet := &types.Wallet{
		Name:         request.Name,
		Pubkey:       pubKeyBytes,
		AddressBytes: pubKeyHash[:],
		Address:      bech32Addr,
	}

	s.logger.Info("Wallet created successfully", zap.String("wallet_name", request.Name))

	return &types.CreateWalletResponse{
		Wallet: newWallet,
	}, nil
}

func (s *APIServerImpl) GetWallet(ctx context.Context, request *types.GetWalletRequest) (*types.GetWalletResponse, error) {
	logging.FromContext(ctx).Info("GetWallet", zap.String("name", request.Name))

	if request.Name == "" {
		s.logger.Error("wallet name cannot be empty")
		return nil, status.Error(codes.NotFound, "wallet name cannot be empty")
	}

	// Check if the wallet name exists in the wallet names list
	found := false
	for _, name := range s.walletNames {
		if name == request.Name {
			found = true
			break
		}
	}
	if !found {
		s.logger.Error("wallet not found", zap.String("wallet_name", request.Name))
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	// Read the key information from Vault
	pubKeyResponse, err := s.vaultClient.Logical().Read(fmt.Sprintf("transit/keys/%s", request.Name))
	if err != nil {
		s.logger.Error("failed to retrieve public key from Vault", zap.Error(err))
		return &types.GetWalletResponse{
			Error: &types.Error{
				Message: "failed to retrieve public key",
			},
		}, err
	}

	// Handle missing data in response
	if pubKeyResponse == nil || pubKeyResponse.Data == nil {
		s.logger.Error("no data in response when retrieving public key")
		return &types.GetWalletResponse{
			Error: &types.Error{
				Message: "no data in response",
			},
		}, fmt.Errorf("no data in response")
	}

	keysInterface, ok := pubKeyResponse.Data["keys"].(map[string]interface{})
	if !ok || keysInterface == nil {
		s.logger.Error("keys map not found in response")
		return &types.GetWalletResponse{
			Error: &types.Error{
				Message: "keys map not found in response",
			},
		}, fmt.Errorf("keys map not found in response")
	}

	keyData, ok := keysInterface["1"].(map[string]interface{})
	if !ok || keyData == nil {
		s.logger.Error("key version data not found in response")
		return &types.GetWalletResponse{
			Error: &types.Error{
				Message: "key version data not found in response",
			},
		}, fmt.Errorf("key version data not found in response")
	}

	pubKeyInterface, ok := keyData["public_key"]
	if !ok {
		s.logger.Error("public key not found in response")
		return &types.GetWalletResponse{
			Error: &types.Error{
				Message: "public key not found in response",
			},
		}, fmt.Errorf("public key not found in response")
	}

	pubKey, ok := pubKeyInterface.(string)
	if !ok {
		s.logger.Error("public key format error")
		return &types.GetWalletResponse{
			Error: &types.Error{
				Message: "public key format error",
			},
		}, fmt.Errorf("public key format error")
	}

	// Decode the PEM encoded public key
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		s.logger.Error("failed to parse PEM block containing the public key")
		return &types.GetWalletResponse{
			Error: &types.Error{
				Message: "PEM block parsing failed",
			},
		}, fmt.Errorf("PEM block parsing failed")
	}

	pubKeyBytes := block.Bytes

	pubKeyHash := sha256.Sum256(pubKeyBytes)
	bech32Addr, err := sdk.Bech32ifyAddressBytes("cosmos", pubKeyHash[:])
	if err != nil {
		s.logger.Error("failed to convert address to Bech32", zap.Error(err))
		return &types.GetWalletResponse{
			Error: &types.Error{
				Message: "failed to convert address to Bech32",
			},
		}, err
	}

	// Construct the Wallet object to return
	wallet := &types.Wallet{
		Name:         request.Name,
		Pubkey:       block.Bytes,
		AddressBytes: pubKeyHash[:],
		Address:      bech32Addr,
	}

	s.logger.Info("Wallet retrieved successfully", zap.String("wallet_name", request.Name))

	return &types.GetWalletResponse{
		Wallet: wallet,
	}, nil
}

func (s *APIServerImpl) GetWallets(ctx context.Context, request *types.EmptyRequest) (*types.GetWalletsResponse, error) {
	logging.FromContext(ctx).Info("GetWallets")

	var wallets []*types.Wallet
	var errorsEncountered bool

	// Iterate through all wallet names stored in memory
	for _, name := range s.walletNames {
		// Use the GetWallet method to fetch each wallet by its name
		response, err := s.GetWallet(ctx, &types.GetWalletRequest{Name: name})
		if err != nil {
			s.logger.Error("error retrieving wallet", zap.String("wallet_name", name), zap.Error(err))
			errorsEncountered = true
			continue // Skip to the next wallet if there's an error
		}
		if response.Error != nil {
			s.logger.Error("error in GetWallet response", zap.String("wallet_name", name), zap.String("error_message", response.Error.Message))
			errorsEncountered = true
			continue
		}

		// Add the wallet to the list if successfully retrieved
		wallets = append(wallets, response.Wallet)
	}

	if errorsEncountered && len(wallets) == 0 {
		return &types.GetWalletsResponse{
			Error: &types.Error{
				Message: "no wallets could be retrieved",
			},
		}, nil
	}

	return &types.GetWalletsResponse{
		Wallets: wallets,
	}, nil
}

func HexToBech32(hexAddr string, prefix string) (string, error) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(prefix, prefix+"pub")

	addrBytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex address: %w", err)
	}

	address := sdk.AccAddress(addrBytes)

	return address.String(), nil
}

func (s *APIServerImpl) Sign(ctx context.Context, request *types.WalletSignatureRequest) (*types.WalletSignatureResponse, error) {
	logging.FromContext(ctx).Info("Signing transaction", zap.String("wallet_name", request.WalletName), zap.ByteString("tx_bytes", request.TxBytes))

	if request.WalletName == "" {
		s.logger.Error("wallet name cannot be empty")
		return nil, status.Error(codes.InvalidArgument, "wallet name cannot be empty")
	}

	path := fmt.Sprintf("transit/sign/%s/sha2-256", request.WalletName)
	payload := map[string]interface{}{
		"input": base64.StdEncoding.EncodeToString(request.TxBytes),
	}

	response, err := s.vaultClient.Logical().Write(path, payload)
	if err != nil {
		s.logger.Error("failed to sign transaction with Vault", zap.Error(err))
		return &types.WalletSignatureResponse{
			Error: &types.Error{
				Message: "failed to sign transaction",
			},
		}, err
	}

	s.logger.Info("Vault response for signature", zap.Any("response", response))

	signatureData, found := response.Data["signature"]
	if !found || signatureData == nil {
		s.logger.Error("signature not found in Vault response")
		return &types.WalletSignatureResponse{
			Error: &types.Error{
				Message: "signature not found in Vault response",
			},
		}, fmt.Errorf("signature not found in Vault response")
	}

	signatureStr, ok := signatureData.(string)
	if !ok {
		s.logger.Error("signature format invalid", zap.Any("signatureData", signatureData))
		return &types.WalletSignatureResponse{
			Error: &types.Error{
				Message: "signature format invalid",
			},
		}, fmt.Errorf("signature format invalid")
	}

	trimmedSignature := strings.TrimPrefix(signatureStr, "vault:v1:")

	s.logger.Info("Trimmed raw signature data", zap.String("trimmedSignature", trimmedSignature))

	signatureBytes, err := base64.StdEncoding.DecodeString(trimmedSignature)
	if err != nil {
		s.logger.Error("failed to decode signature", zap.Error(err))
		return &types.WalletSignatureResponse{
			Error: &types.Error{
				Message: "failed to decode signature",
			},
		}, err
	}

	// Decode DER signature to get R and S
	ecdsaSig, err := decodeDERSignature(signatureBytes)
	if err != nil {
		s.logger.Error("failed to decode DER signature", zap.Error(err))
		return &types.WalletSignatureResponse{
			Error: &types.Error{
				Message: "DER signature parsing failed",
			},
		}, err
	}

	// Ensure each value is 32 bytes
	rBytes := ecdsaSig.R.Bytes()
	sBytes := ecdsaSig.S.Bytes()
	if len(rBytes) < 32 {
		rBytes = append(make([]byte, 32-len(rBytes)), rBytes...)
	}
	if len(sBytes) < 32 {
		sBytes = append(make([]byte, 32-len(sBytes)), sBytes...)
	}

	// Convert R and S to a raw signature format expected by your Verify function
	rawSig := append(rBytes, sBytes...)

	s.logger.Info("Signature decoded successfully", zap.ByteString("signature", signatureBytes))

	return &types.WalletSignatureResponse{
		Signature: rawSig,
	}, nil
}

type ECDSASignature struct {
	R, S *big.Int
}

func decodeDERSignature(der []byte) (*ECDSASignature, error) {
	var sig ECDSASignature
	_, err := asn1.Unmarshal(der, &sig)
	if err != nil {
		return nil, err
	}
	return &sig, nil
}
