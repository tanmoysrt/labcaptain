#!/usr/bin/env bash

# Credits goes to https://github.com/ConSol/docker-headless-vnc-container (Apache License 2.0)
# The scripts have helped a lot to curate this vnc setup script

set -e

# Install common tools
echo "Install some common tools for further installation"
apt install -y vim wget net-tools locales bzip2 procps apt-utils python3-numpy #used for websockify/novnc
echo "generate locales für en_US.UTF-8"
echo "en_US.UTF-8 UTF-8" > /etc/locale.gen
locale-gen

# Install custom fonts
export LANG=en_US.UTF-8
export LANGUAGE=en_US:en
export LC_ALL=en_US.UTF-8
echo "Installing ttf-wqy-zenhei"
apt install -y ttf-wqy-zenhei

# Install tigervnc
echo "Install TigerVNC server"
apt install -y tigervnc-standalone-server
printf '\n# docker-headless-vnc-container:\n$localhost = "no";\n1;\n' >>/etc/tigervnc/vncserver-config-defaults

# Install noVNC
echo "Install noVNC - HTML5 based VNC viewer"
mkdir -p $NO_VNC_HOME/utils/websockify
wget -qO- https://github.com/novnc/noVNC/archive/refs/tags/v1.3.0.tar.gz | tar xz --strip 1 -C $NO_VNC_HOME
# use older version of websockify to prevent hanging connections on offline containers, see https://github.com/ConSol/docker-headless-vnc-container/issues/50
wget -qO- https://github.com/novnc/websockify/archive/refs/tags/v0.10.0.tar.gz | tar xz --strip 1 -C $NO_VNC_HOME/utils/websockify
# create index.html to forward automatically to `vnc_lite.html`
ln -s $NO_VNC_HOME/vnc.html $NO_VNC_HOME/index.html

# Install xfce ui
echo "Install Xfce4 UI components"
apt update
apt install -y supervisor xfce4  xfce4-terminal xterm dbus-x11 libdbus-glib-1-2

# Change permission of /home/labuser/.config files
sudo chown -R labuser:labuser /home/labuser/.config