build:
	go build -o ./bin/goblockchain

run: build
	./bin/goblockchain
