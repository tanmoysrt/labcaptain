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

# Search for a free port in the range 10000-30000
find_free_port() {
    for ((port=3000; port<=30000; port++)); do
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

# replace the template with the port
container_quadlet=$(echo "$container_quadlet_template" | sed "s/{{lab_published_port}}/$port/g")
# create the file
echo -e "$container_quadlet" > /etc/containers/systemd/{{lab_id}}.container
# daemon-reload without anything on stdout and stderr
systemctl daemon-reload > /dev/null 2>&1
# start the service
systemctl start {{lab_id}}.service >/dev/null 2>&1
# print the port no
echo "assigned_${port}_port"
