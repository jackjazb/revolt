be:
	(cd revolt-server && go run .)
fe:
	(cd revolt-frontend && pnpm dev)

test:
	(cd revolt-server && go test ./...)
	(cd revolt-frontend && pnpm test run)
