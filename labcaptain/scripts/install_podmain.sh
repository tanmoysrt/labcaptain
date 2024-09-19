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
            debian|ubuntu|linuxmint|kali)
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

# Main script execution
detect_debian
