#!/bin/bash

# Function to check if podman is already installed
check_podman() {
    if command -v podman &> /dev/null; then
        echo "Podman is already installed."
        return 0
    else
        return 1
    fi
}

# Function to install podman on Debian-based systems
install_debian() {
    echo "Installing Podman on a Debian-based system..."
    sudo apt update -y
    sudo apt install -y podman
}

# Detect if the system is Debian-based
detect_debian() {
    if [ -f /etc/os-release ]; then
        source /etc/os-release
        case "$ID" in
            debian|ubuntu|linuxmint)
                if ! check_podman; then
                    install_debian
                fi
                ;;
            *)
                echo "This script supports only Debian-based systems."
                exit 1
                ;;
        esac
    else
        echo "Cannot detect operating system."
        exit 1
    fi
}

# Function to create a systemd service for podman system service
create_podman_service() {
    SERVICE_FILE_PATH="/etc/systemd/system/podman-service.service"
    PODMAN_BINARY=$(which podman)

    echo "Creating systemd service file at $SERVICE_FILE_PATH..."

    sudo bash -c "cat > $SERVICE_FILE_PATH" <<EOF
[Unit]
Description=Podman System Service
After=network.target

[Service]
ExecStart=$PODMAN_BINARY system service --time=0 tcp://0.0.0.0:1001
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

    # Reload systemd to recognize the new service
    echo "Reloading systemd..."
    sudo systemctl daemon-reload

    # Enable the service to start on boot
    echo "Enabling podman-service..."
    sudo systemctl enable podman-service

    # Start the service
    echo "Starting podman-service..."
    sudo systemctl start podman-service

    # Check the service status
    sudo systemctl status podman-service
}
# Main script execution
detect_debian
create_podman_service
