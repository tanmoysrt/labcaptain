<img src="./assets/logo.png" width="200">

### What is LabCaptain ?
LabCaptain is a tiny daemon + CLI tool that helps to deploy lab environment in cluster.

### Why LabCaptain ?
I have planned to build a lab environment or dev platform first.
Then thought about this abstraction layer which will help to build and deploy lab environment easily.
So, I decided to build this lightweight tool to manage it.

The name `LabCaptain` is due to the fact that it takes the responsibility of lab environment management. You will just say it you need to deploy a lab environment for certain period of time and it will do the task.

### Tech Stack
- Golang
- SQLite3 (for database)
- Podman
- Prometheus
- HAProxy
- NoVNC + noVNC proxy
- ttyd

### Installation guide (Ubuntu 22.04)
1. Install golang (https://go.dev/doc/install)
2. Clone the repo
3. Go inside `labcaptain` folder
4. Run `go build` to build the binary
5. Move the `labcaptain` binary to `/usr/local/bin`
6. Run `labcaptain local-setup` to setup labcaptain on the local machine
7. Run `labcaptain server add <ip>` to add a new server
8. Run `labcaptain server list` to list all servers
9. Run `labcaptain server setup-podman <ip>` to setup podman on the server
10. Run `labcaptain server setup-prometheus <ip>` to setup prometheus exporter on the server
11. Run `labcaptain server enable <ip>` to enable a server
12. Create a systemd service file `/etc/systemd/system/labcaptain.service`
```bash
[Unit]
Description=LabCaptain
After=network.target

[Service]
User=root
Type=simple
Environment="LAB_CAPTAIN_BASE_DOMAIN=example.com"
Environment="LABCAPTAIN_API_TOKEN=random_secret"
ExecStart=/usr/local/bin/labcaptain
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```
13. Run `systemctl daemon-reload` and `systemctl enable labcaptain`
14. Run `systemctl start labcaptain` to start labcaptain

**Note:** Your server may have not ssh-agent installed. Run `echo $SSH_AUTH_SOCK` to check if it's installed.

- Generate a ssh private key for all of your other servers (if you haven't already)
  ```bash
  ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
  ```
  Put the private key in `/root/.ssh/id_rsa` and the public key in `/root/.ssh/id_rsa.pub`

- Add the public key to all of your other servers
- Create/edit `/etc/rc.local` file and add the following line
  ```bash
  eval $(ssh-agent -s)
  ssh-add /root/.ssh/id_rsa
  ```
  and make that file executable with `chmod +x /etc/rc.local` and reboot the server.

### Future Work
-[] Implement SSH connection pool for faster communication
-[] Implement support for resource limits
-[] Configurable option for lab proxy at port 443 ssl (P.S : currently it's also possible by editing `labcaptain/nginx.conf.template`)

### License
MIT License

### Credits
Special thanks to these projects
- https://github.com/ConSol/docker-headless-vnc-container (Apache License 2.0). It made it possible to have proper configuration for Xfce4 + noVNC setup for headsup.
- https://github.com/tsl0922/ttyd (MIT License). It made it possible to have a web terminal.
- https://github.com/coder/code-server (MIT License). It made it possible to have a code server.
