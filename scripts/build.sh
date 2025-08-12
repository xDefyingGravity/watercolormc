#!/bin/zsh
# Build frontend to dist/static and backend to dist/<platform>/<arch>/app

PLATFORM=$1
ARCH=$2

if [ -z "$PLATFORM" ]; then
  echo "usage: $0 <platform> <arch>"
  echo "example: $0 linux amd64"
  echo "supported platforms: linux, windows, darwin"
  echo "supported architectures: amd64, arm64"
  echo "special: use 'current' for both to build for your current system"
  exit 1
fi

if [ "$PLATFORM" = "current" ] || [ "$ARCH" = "current" ]; then
  PLATFORM=$(uname | tr '[:upper:]' '[:lower:]')
  case $PLATFORM in
    darwin) PLATFORM=darwin ;;
    linux) PLATFORM=linux ;;
    msys*|mingw*|cygwin*) PLATFORM=windows ;;
  esac
  ARCH=$(uname -m)
  case $ARCH in
    x86_64) ARCH=amd64 ;;
    arm64|aarch64) ARCH=arm64 ;;
    *) echo "[error] unsupported architecture: $ARCH" && exit 1 ;;
  esac
fi

echo "[build] target: $PLATFORM/$ARCH"

# Frontend output directory
STATIC_OUT="dist/static"
rm -rf "$STATIC_OUT"
mkdir -p "$STATIC_OUT"
STATIC_OUT=$(realpath "$STATIC_OUT")

# Backend output directory
BACKEND_OUT="dist/$PLATFORM/$ARCH"
mkdir -p "$BACKEND_OUT"
BACKEND_OUT=$(realpath "$BACKEND_OUT")

# --- Build Frontend ---
echo "[build] building frontend..."
cd frontend || exit 1

if [ ! -d "node_modules" ]; then
  npm install
fi

npm run build
if [ $? -ne 0 ]; then
  echo "[error] frontend build failed"
  exit 1
fi

cp -r build/* "$STATIC_OUT"

# --- Build Backend ---
echo "[build] building backend..."
cd .. || exit 1

BIN_NAME="app"
if [[ "$PLATFORM" == "windows" ]]; then
  BIN_NAME="app.exe"
fi

# Use mingw-w64 cross compiler if building for windows on macOS
if [[ "$PLATFORM" == "windows" ]] && [[ "$(uname)" == "Darwin" ]]; then
  if [[ "$ARCH" == "amd64" ]]; then
    export CC=x86_64-w64-mingw32-gcc
    export CXX=x86_64-w64-mingw32-g++
  else
    echo "[error] unsupported windows architecture for mingw: $ARCH, only amd64 is supported when compiling on MacOS."
    exit 1
  fi
  export CGO_ENABLED=1
  export CGO_CFLAGS=""
  export CGO_CPPFLAGS=""
  export CGO_LDFLAGS=""
  GOOS=windows GOARCH=$ARCH go build -o "$BACKEND_OUT/$BIN_NAME"
else
  CGO_ENABLED=1 GOOS=$PLATFORM GOARCH=$ARCH go build -o "$BACKEND_OUT/$BIN_NAME"
fi

cp "$BACKEND_OUT/$BIN_NAME" ./frontend/src-tauri/resources/

cd ./frontend || exit 1

npx tauri build

cp -r ./src-tauri/target/release/bundle/* "$BACKEND_OUT/"


echo "[done] built frontend to $STATIC_OUT and backend to $BACKEND_OUT/$BIN_NAME"