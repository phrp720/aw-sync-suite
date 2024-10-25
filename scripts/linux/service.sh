#!/bin/bash

APP_NAME="aw-sync-agent"
APP_PATH="/bin/aw/agent"
SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"
GO_VERSION="go1.23.2"

## Check if Go 1.23 is installed
#if ! go version | grep -q "$GO_VERSION"; then
#  echo "Go $GO_VERSION is not installed. Installing..."
#  wget https://golang.org/dl/$GO_VERSION.linux-amd64.tar.gz
#  sudo tar -C /usr/local -xzf $GO_VERSION.linux-amd64.tar.gz
#  echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
#  source ~/.profile
#else
#  echo "Go $GO_VERSION is already installed."
#fi

# Build the Go application
echo "Building ActivityWatch Sync Agent..."
go build -o $APP_PATH main.go
cp -r .env /bin/aw/.env

# Create systemd service file
echo "Creating systemd service file..."
sudo bash -c "cat > $SERVICE_FILE <<EOF
[Unit]
Description=ActivityWatch Sync Agent
After=network.target

[Service]
ExecStart=$APP_PATH -service
Restart=always
User=$(whoami)
Group=$(whoami)
WorkingDirectory=$(dirname $APP_PATH)
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF"

# Reload systemd, enable and start service
echo "Reloading systemd and starting service..."
sudo systemctl daemon-reload
sudo systemctl enable $APP_NAME
sudo systemctl start $APP_NAME

