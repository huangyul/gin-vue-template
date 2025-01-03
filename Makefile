WEB_DIR=./web

.PHONY: web
web:
	@nvm use 22.11.0 && node -v && cd $(WEB_DIR) && pnpm dev