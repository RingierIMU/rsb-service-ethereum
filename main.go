package main

import (
	"context"
	"crypto/ecdsa"
	"embed"
	"github.com/RingierIMU/rsb-service-ethereum/contracts"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"io"
	"io/fs"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"
)

var (
	ethNode         = os.Getenv("ETH_NODE")
	mnemonic        = os.Getenv("MNEMONIC")
	contractAddress = os.Getenv("CONTRACT_ADDRESS")
	gasLimit        = uint64(900000)

	//go:embed ui/build/*
	content embed.FS
)

func main() {
	go deploy("raffle")
	runWebServer()
}

func runWebServer() {
	mux := http.NewServeMux()

	webRoot, err := fs.Sub(content, "ui/build")
	if err != nil {
		log.Fatal("unable to find web root: ", err)
	}

	mux.Handle("/", http.FileServer(http.FS(webRoot)))
	mux.HandleFunc("/contract-address", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, contractAddress)
		if err != nil {
			log.Printf("unable to find write contract address to HTTP client: %s\n", err)
		}
	})
	mux.HandleFunc("/contract-abi", func(w http.ResponseWriter, r *http.Request) {
		abi, err := os.ReadFile("./contracts/Raffle.abi")
		_, err = io.WriteString(w, string(abi))
		if err != nil {
			log.Printf("unable to find write ABI to HTTP client: %s\n", err)
		}
	})

	log.Printf("Starting webserver on http://localhost:8080/,"+
		" using %q as the ethereum node and %q as the mnemonic.\n", ethNode, mnemonic)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("unable to start webserver: ", err)
	}
}

// deploy deploys the given contract to an ethereum blockchain
func deploy(contractName string) {
	client, privateKey, nonce, gasPrice, err := prepDeploy()
	if err != nil {
		log.Fatal("error preparing deploy: ", err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	switch contractName {
	case "raffle":
		address, tx, instance, err := contracts.DeployRaffle(auth, client)
		if err != nil {
			log.Fatal("unable to deploy contract: ", err)
		}
		log.Printf("Contract address: %q, Transaction ID: %q\n", address.Hex(), tx.Hash().Hex())
		contractAddress = address.Hex()
		time.Sleep(time.Minute)
		_ = instance
	default:
		log.Fatal("Wrong contract name")
	}
}

// openWallet opens a wallet from the given mnemonic and returns it alongside an account
func openWallet(mnemonic string) (*hdwallet.Wallet, *accounts.Account, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, nil, err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, nil, err
	}

	return wallet, &account, nil
}

// prepDeploy connects to the Ethereum node and prepares the deployment transaction
func prepDeploy() (*ethclient.Client, *ecdsa.PrivateKey, uint64, *big.Int, error) {
	client, err := ethclient.Dial(ethNode)
	if err != nil {
		log.Fatal("error connecting to network: ", err)
	}

	wallet, account, err := openWallet(mnemonic)
	if err != nil {
		log.Fatal("error opening wallet: ", err)
	}

	p, err := wallet.PrivateKeyHex(*account)
	if err != nil {
		log.Fatal("error extracting private key: ", err)
	}

	privateKey, err := crypto.HexToECDSA(p)
	if err != nil {
		log.Fatal("error converting private key: ", err)
	}

	publicKey, err := wallet.PublicKey(*account)
	if err != nil {
		log.Fatal("error converting public key to hex: ", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("error extracting address from public key: ", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("error getting gas price prediction: ", err)
	}
	return client, privateKey, nonce, gasPrice, err
}
