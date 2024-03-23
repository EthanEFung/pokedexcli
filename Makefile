run:
	go build && ./pokedexcli

filecache-reset:
	rm -rf internal/filebasedcache/files/*
	rm internal/filebasedcache/ledger.txt
