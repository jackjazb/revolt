be:
	(cd revolt-server && go run .)
fe:
	(cd revolt-frontend && pnpm dev)

test_be:
	(cd revolt-server && go test ./...)
