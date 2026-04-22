.PHONY: dev test backend frontend clean

dev:
    ./dev.sh

test:
    cd server && ./test.sh
    cd frontend && ./test.sh

backend:
    cd server && air

frontend:
    cd frontend && pnpm dev

clean:
    rm -rf server/tmp
    rm -rf frontend/node_modules
