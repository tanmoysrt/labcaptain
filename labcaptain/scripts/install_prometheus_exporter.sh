#!/bin/bash

# Script to install and configure Prometheus Node Exporter on Debian-based systems

# Variables
VERSION="1.8.2"  # Replace with the latest stable version if needed
DOWNLOAD_URL="https://github.com/prometheus/node_exporter/releases/download/v$VERSION/node_exporter-$VERSION.linux-amd64.tar.gz"
INSTALL_DIR="/usr/local/bin"
USER="node_exporter"
SERVICE_FILE="/etc/systemd/system/node_exporter.service"

# Update system packages
echo "Updating system packages..."
sudo apt update

# Install required packages
echo "Installing necessary packages..."
sudo apt install -y wget tar

# Create node_exporter user if not exists
if ! id -u "$USER" > /dev/null 2>&1; then
    echo "Creating $USER user..."
    sudo useradd --no-create-home --shell /bin/false $USER
fi

# Download and install node_exporter
echo "Downloading Node Exporter..."
wget -q $DOWNLOAD_URL -O /tmp/node_exporter.tar.gz

echo "Extracting Node Exporter..."
tar -xzf /tmp/node_exporter.tar.gz -C /tmp

echo "Moving Node Exporter binary to $INSTALL_DIR..."
sudo mv /tmp/node_exporter-$VERSION.linux-amd64/node_exporter $INSTALL_DIR/

# Set correct permissions
echo "Setting permissions for node_exporter..."
sudo chown $USER:$USER $INSTALL_DIR/node_exporter
sudo chmod +x $INSTALL_DIR/node_exporter

# Create systemd service file
echo "Creating systemd service file..."
sudo bash -c "cat > $SERVICE_FILE" << EOF
[Unit]
Description=Prometheus Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=$USER
Group=$USER
Type=simple
ExecStart=$INSTALL_DIR/node_exporter

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and start node_exporter service
echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

echo "Enabling Node Exporter service to start on boot..."
sudo systemctl enable node_exporter

echo "Starting Node Exporter service..."
sudo systemctl start node_exporter

# Check if the service is running
if systemctl is-active --quiet node_exporter; then
    echo "Node Exporter is running successfully!"
    echo "You can verify it by accessing http://<server-ip>:9100/metrics"
else
    echo "Failed to start Node Exporter. Check the logs for more details."
fi

# Cleanup
echo "Cleaning up temporary files..."
rm -rf /tmp/node_exporter*

echo "Node Exporter setup is complete."
