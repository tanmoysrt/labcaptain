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
    # Custom installation as v4 is not available in Ubuntu 22.04 repositories
    ubuntu_version='22.04'
    key_url="https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/unstable/xUbuntu_${ubuntu_version}/Release.key"
    sources_url="https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/unstable/xUbuntu_${ubuntu_version}"

    echo "deb $sources_url/ /" | sudo tee /etc/apt/sources.list.d/devel:kubic:libcontainers:unstable.list
    curl -fsSL $key_url | sudo gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/devel_kubic_libcontainers_unstable.gpg > /dev/null
    sudo apt update -y
    sudo apt install podman
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
