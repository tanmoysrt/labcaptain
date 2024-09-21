#!/bin/bash

# Function to check if a port is free
is_port_free() {
    local port=$1
    if ! lsof -i:$port >/dev/null 2>&1; then
        return 0  # Port is free
    else
        return 1  # Port is in use
    fi
}

# Function to get assigned ports from existing .container files
get_assigned_ports() {
    local assigned_ports=()
    for file in /etc/containers/systemd/*.container; do
        if [[ -f $file ]]; then
            # Extract the published port from the file
            port=$(grep -oP 'PublishPort=\K\d+' "$file")
            if [[ -n $port ]]; then
                assigned_ports+=("$port")
            fi
        fi
    done
    echo "${assigned_ports[@]}"
}

# Search for a free port in the range 10000-30000
find_free_port() {
    local assigned_ports
    assigned_ports=($(get_assigned_ports))

    for ((port=10000; port<=30000; port++)); do
        # First check if the port is in the assigned list
        if [[ " ${assigned_ports[@]} " =~ " ${port} " ]]; then
            continue  # Skip to the next port
        fi

        # Then check if the port is free
        if is_port_free $port; then
            echo $port
            return 0
        fi
    done
    echo "No free ports found in the range 10000-30000"
    exit 1
}

# Template for container
container_quadlet_template="
[Container]\n
Environment={{lab_environment_variables}}\n
Image={{lab_image}}\n
PublishPort={{lab_published_port}}:80
"

# Run the function to find a free port
port=$(find_free_port)

# Replace the template with the port
container_quadlet=$(echo -e "$container_quadlet_template" | sed "s/{{lab_published_port}}/$port/g")
# Create the file
echo -e "$container_quadlet" > /etc/containers/systemd/{{lab_id}}.container
# Daemon-reload without anything on stdout and stderr
systemctl daemon-reload > /dev/null 2>&1
# Start the service
systemctl start {{lab_id}}.service >/dev/null 2>&1
# Print the port number
echo "assigned_${port}_port"
