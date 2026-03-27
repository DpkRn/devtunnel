#!/bin/bash

echo "Installing mytunnel..."

OS=$(uname)

if [ "$OS" = "Linux" ]; then
    URL="https://github.com/DpkRn/gotunnel/releases/download/v0.1.0/mytunnel-linux"
elif [ "$OS" = "Darwin" ]; then
    URL="https://github.com/DpkRn/gotunnel/releases/download/v0.1.0/mytunnel-mac"
else
    echo "Unsupported OS"
    exit 1
fi

curl -L $URL -o mytunnel
chmod +x mytunnel
sudo mv mytunnel /usr/local/bin/

echo "✅ Installed! Run: mytunnel http 3000"