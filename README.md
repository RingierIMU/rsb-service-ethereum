# rsb-service-ethereum
A service that interacts with the Ethereum blockchain by deploying a smart contract for a raffle which is then used via a React frontend.

# Setup

This service mostly works via the Makefile:

Clean the workspace:

```
make clean
```

Compile the frontend:

```
make frontend
```

Compile the contracts, create ABI & Go code of them:

```
make contracts
```

# Run

To run the application, you must set two environment variables and then run the main.go file:

```
export ETH_NODE=""
export MNEMONIC=""
go run main.go
```

This starts a webserver which uses ./ui/build as the webroot to serve the React application (Note that you have to compile the contracts and the frontend first).

ETH_NODE is the URL of the Ethereum node to connect to (e.g. "http://localhost:8545").
MNEMONIC is the 12-word mnemonic for your wallet (We suggest to use Metamask to create one).

# Tooling

You need a working [Go](https://go.dev/) and [NodeJS](https://nodejs.org/) / [ReactJS](https://reactjs.org/) (npm & yarn) environment as well as the solc provider:

```
brew install go node yarn solidity ethereum
```

To make the application work, you need an Ethereum Blockchain. We recommend using [Ganache](https://github.com/trufflesuite/ganache). You can also use [Infura](https://infura.io/) and connect to the [Görli testnet](https://goerli.net/) or other networks. You can use [Metamask](https://metamask.io/) as your wallet.
