repl:
	go run $(filter-out %_test.go,$(wildcard cmd/repl/*.go))


filecache-reset:
	rm -rf internal/filebasedcache/files/*
	rm internal/filebasedcache/ledger.txt
