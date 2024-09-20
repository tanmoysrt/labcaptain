# Build proxy
FROM --platform=$BUILDPLATFORM golang:1.21-rc-bookworm AS proxy-build-env
ENV CGO_ENABLED=1
WORKDIR /src
COPY ./proxy/* .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o port_proxy .

# Actual image
FROM ubuntu:22.04

# set environment variables
ENV DEBIAN_FRONTEND=noninteractive \
    DISPLAY=:1 \
    TERM=xterm \
    NO_VNC_HOME=/home/labuser/noVNC \
    USER_PASSWORD=12345 \
    ENABLE_WEB_TERMINAL=1 \
    ENABLE_CODE_SERVER=1 \
    ENABLE_VNC=1 \
    ENABLE_PORT_PROXY=1


# base packages
RUN apt update -y && apt install -y \
    build-essential \
    cmake \
    curl \
    wget \
    python3 \
    python3-pip \
    python3-dev \
    python3-venv \
    nodejs \
    npm \
    git \
    sudo \
    vim \
    htop \
    nano \
    tini

# Add a new user
RUN useradd -m -s /bin/bash labuser && \
    usermod -aG sudo labuser && \
    echo "labuser ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# install ttyd
COPY ./ttyd.i686 /usr/local/bin/ttyd
RUN chmod +x /usr/local/bin/ttyd

# Install VNC Server + noVNC + xfce UI
ADD ./xfce/ /home/labuser
ADD  ./setup-vnc.sh /setup-vnc.sh
RUN /setup-vnc.sh && rm -rf /setup-vnc.sh
ADD ./xstartup /home/labuser/.vnc/xstartup
RUN chown -R labuser:labuser /home/labuser/.vnc

# Install HAProxy and configure
RUN curl https://haproxy.debian.net/bernat.debian.org.gpg | gpg --dearmor > /usr/share/keyrings/haproxy.debian.net.gpg && \
    echo deb "[signed-by=/usr/share/keyrings/haproxy.debian.net.gpg]" http://haproxy.debian.net bookworm-backports-3.0 main > /etc/apt/sources.list.d/haproxy.list && \
    apt update -y && apt install -y haproxy=3.0.\*
ADD ./haproxy.cfg /etc/haproxy/haproxy.cfg

# Switch to new user
USER labuser

# Create a working directory
RUN mkdir -p /home/labuser/Desktop && mkdir -p /home/labuser/Desktop/lab

# Install code server
RUN curl -fsSL https://code-server.dev/install.sh | sh

# Add port proxy
COPY --from=proxy-build-env /src/port_proxy /usr/local/bin/port_proxy

# Add Xstartup
RUN chmod +x /home/labuser/.vnc/xstartup

# EXPOSE PORT
EXPOSE 80

# Set entrypoint
COPY ./entrypoint.sh /

ENTRYPOINT ["tini", "--", "/entrypoint.sh"]
