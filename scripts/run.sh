#!/bin/zsh
mkdir -p "dist/static"

concurrently \
  --names "backend,frontend" \
  --prefix "[{name}]" \
  "APP_ENV=development go run ." \
  "cd frontend && npm run dev"