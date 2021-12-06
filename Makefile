.PHONY: all contracts bin abi clean-contracts clean clean-frontend frontend run go-update

all: contracts frontend run

contracts: clean-contracts bin abi
	@abigen --bin=./contracts/raffle.bin --abi=./contracts/raffle.abi --pkg=contracts --type=raffle --out=./contracts/raffle.go

bin:
	@solc --bin -o ./contracts/ ./contracts/raffle.sol

abi:
	@solc --abi -o ./contracts/ ./contracts/raffle.sol

clean-contracts:
	@rm -f ./contracts/raffle.bin ./contracts/raffle.abi ./contracts/raffle.go

clean: clean-contracts clean-frontend
	@rm -f ./rsb-service-example

clean-frontend:
	@rm -rf ./ui/build

frontend: clean-frontend
	@cd ./ui && npm install
	@cd ./ui && yarn build

run:
	@go run main.go

go-update:
	@gofmt -w .
	@go get -u ./...
	@go mod tidy
