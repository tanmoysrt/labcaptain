#!/bin/bash

set -e
# add alias python=python3 in .bashrc
echo "alias python=python3" >> "$HOME/.bashrc"

source "$HOME/.bashrc"

# Set password for labuser
echo "labuser:$USER_PASSWORD" | sudo chpasswd

# Set capability to haproxy
sudo setcap 'cap_net_bind_service=+ep' /usr/sbin/haproxy

# Function to start services
start_service() {
    "$@" &
    echo "$1 started"
}

# Function to terminate all background jobs with SIGKILL
cleanup() {
    echo "Received exit signal, killing all processes with SIGKILL..."
    kill -9 0  # This will send SIGKILL to all processes in the current process group
}

# Set up traps to handle signals
trap cleanup SIGINT SIGTERM

# Start services based on environment variables
if [ "$ENABLE_WEB_TERMINAL" == "1" ]; then
    start_service ttyd -p 8001 -W -w /home/labuser/Desktop/lab bash
else
    echo "ttyd disabled"
fi

if [ "$ENABLE_CODE_SERVER" == "1" ]; then
    start_service code-server --bind-addr 0.0.0.0:8002 --auth none --disable-telemetry --disable-workspace-trust --disable-update-check
else
    echo "code-server disabled"
fi

if [ "$ENABLE_VNC" == "1" ]; then
    VNC_COL_DEPTH=24
    touch ~/.Xauthority
    sudo mkdir -p /tmp/.ICE-unix && \
    sudo chown root:root /tmp/.ICE-unix && \
    sudo chmod 1777 /tmp/.ICE-unix

    start_service $NO_VNC_HOME/utils/novnc_proxy --vnc localhost:5901 --listen 8003

    vncserver :1 -depth 24 -geometry 1920x1080 -SecurityTypes None --I-KNOW-THIS-IS-INSECURE
    echo "VNC server started"
else
    echo "VNC + noVNC proxy disabled"
fi

if [ "$ENABLE_PORT_PROXY" == "1" ]; then
    start_service port_proxy
else
    echo "Port proxy disabled"
fi

# Start HAProxy as ingress
start_service haproxy -f /etc/haproxy/haproxy.cfg

# If `/autostart.sh` exists and is executable, run it
if [ -x /autostart.sh ]; then
    /autostart.sh &
fi

# Wait for all processes to finish
echo "Waiting for processes to finish..."
wait
