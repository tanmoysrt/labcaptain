#!/bin/bash

PROMETHEUS_VERSION="2.53.2"
PROMETHEUS_URL="https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz"
PROMETHEUS_DIR="/usr/local/bin/prometheus-${PROMETHEUS_VERSION}"

# Function to check if a package is installed
check_installed() {
    dpkg -l | grep -q "$1"
}

# Function to install a package if not installed
install_package() {
    if ! check_installed "$1"; then
        echo "$1 is not installed. Installing..."
        apt update && apt install -y "$1"
    else
        echo "$1 is already installed."
    fi
}

# Function to enable and start the service if it's not running
enable_and_start_service() {
    systemctl is-enabled --quiet "$1"
    if [ $? -ne 0 ]; then
        echo "Enabling $1 service..."
        systemctl enable "$1"
    fi

    systemctl is-active --quiet "$1"
    if [ $? -ne 0 ]; then
        echo "Starting $1 service..."
        systemctl start "$1"
    else
        echo "$1 service is already running."
    fi
}

# Install Nginx using apt
install_package "nginx"
enable_and_start_service "nginx"

# Check if Prometheus is already installed
if [ ! -d "$PROMETHEUS_DIR" ]; then
    echo "Downloading and installing Prometheus v${PROMETHEUS_VERSION}..."

    # Download and extract Prometheus
    wget $PROMETHEUS_URL -O /tmp/prometheus.tar.gz
    mkdir -p $PROMETHEUS_DIR
    tar -xzf /tmp/prometheus.tar.gz --strip-components=1 -C $PROMETHEUS_DIR

    # Create symlinks for prometheus and promtool
    ln -sf ${PROMETHEUS_DIR}/prometheus /usr/local/bin/prometheus
    ln -sf ${PROMETHEUS_DIR}/promtool /usr/local/bin/promtool

    # Clean up
    rm /tmp/prometheus.tar.gz

    echo "Prometheus installed successfully."
else
    echo "Prometheus is already installed."
fi

# Create a systemd service for Prometheus
echo "Setting up Prometheus as a service..."
tee /etc/systemd/system/prometheus.service > /dev/null <<EOF
[Unit]
Description=Prometheus Monitoring System
Wants=network-online.target
After=network-online.target

[Service]
User=root
ExecStart=/usr/local/bin/prometheus \
  --config.file=/etc/prometheus/prometheus.yml

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd to recognize the new service, enable, and start Prometheus
systemctl daemon-reload
systemctl enable prometheus
systemctl start prometheus

# Check Prometheus status
if systemctl is-active --quiet prometheus; then
    echo "Prometheus service is running."
else
    echo "Failed to start Prometheus service."
fi
