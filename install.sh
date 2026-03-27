#!/bin/bash

echo "Installing mytunnel..."

OS=$(uname)

if [ "$OS" = "Linux" ]; then
    URL="https://yourdomain.com/mytunnel-linux"
elif [ "$OS" = "Darwin" ]; then
    URL="https://yourdomain.com/mytunnel-mac"
else
    echo "Unsupported OS"
    exit 1
fi

curl -L $URL -o mytunnel
chmod +x mytunnel
sudo mv mytunnel /usr/local/bin/

echo "✅ Installed! Run: mytunnel http 3000"