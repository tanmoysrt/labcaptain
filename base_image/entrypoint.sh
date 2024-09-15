# !/bin/sh -e

set -e
source $HOME/.bashrc
alias python=python3

# Set password for labuser
echo "labuser:$USER_PASSWORD" | sudo chpasswd

# Set capability to haproxy
sudo setcap 'cap_net_bind_service=+ep' /usr/sbin/haproxy

if [ "$ENABLE_WEB_TERMINAL" == "1" ]; then
    ttyd -p 8001 -W -w /home/labuser/Desktop/lab bash &
    echo "ttyd started"
else
    echo "ttyd disabled"
fi

if [ "$ENABLE_CODE_SERVER" == "1" ]; then
    code-server --bind-addr 0.0.0.0:8002 --auth none --disable-telemetry --disable-workspace-trust --disable-update-check  &
    echo "code-server started"
else
    echo "code-server disabled"
fi

if [ "$ENABLE_VNC" == "1" ]; then
    VNC_COL_DEPTH=24
    touch ~/.Xauthority
    sudo mkdir -p /tmp/.ICE-unix && \
    sudo chown root:root /tmp/.ICE-unix && \
    sudo chmod 1777 /tmp/.ICE-unix
    
    $NO_VNC_HOME/utils/novnc_proxy --vnc localhost:5901 --listen 8003 &
    echo "noVNC proxy started"
    
    vncserver :1 -depth 24 -geometry 1920x1080 -SecurityTypes None --I-KNOW-THIS-IS-INSECURE
    echo "VNC server started"
else
    echo "VNC + noVNC proxy disabled"
fi

if [ "$ENABLE_PORT_PROXY" == "1" ]; then
    port_proxy &
    echo "Port proxy started"
else
    echo "Port proxy disabled"
fi

# start haproxy as ingress
haproxy -f /etc/haproxy/haproxy.cfg &

# wait for all processes to finish
exit_all () {
    echo "Received exit signal, killing all processes..."
    exit 0
}

trap exit_all SIGINT SIGTERM SIGKILL SIGQUIT
echo "Waiting for processes to finish..."
wait