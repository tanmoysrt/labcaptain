**labcaptain_base** docker image is a base image to build your own lab environment on top of.

### Why should you use this image?
This image has many things pre-configured to make it easy to get started with a lab environment.

- Built on top of ubuntu 22.04
- Non-root user **labuser** configured
- Xfce4 with TigerVNC and noVNC support
- Web based code editor (VSCode like code-server)
- Web based terminal (ttyd)
- Integrated ingress handling with HAProxy
- Port proxy to handle port-based routing

### How to use this image?
```curl
docker run -it --rm -p 8080:80 tanmoysrt/labcaptain_base:1
```

Or, In your docker image, you can use this image as a base image
```dockerfile
FROM tanmoysrt/labcaptain_base:1
```

### Where to put your lab files?
You can put your lab files in the following folders:
- /home/labuser/Desktop/lab

### How to access lab software?
- Access web terminal at http://terminal:8080
- Access code editor at http://editor:8080
- Access VNC at http://vnc:8080
- Access port-based routing to any service inside the container by http://port-<port_no>:8080

> In the above examples, terminal, editor, vnc and port-<port_no> - all the dummy domains are pointed to 127.0.0.1.

If you want to use any proxy in front of the container, you should set the `Host` header correctly. It should be out of these specific formats:
- terminal
- editor
- vnc
- port-<port_no>
- terminal:any_port_or_string
- editor:any_port_or_string
- vnc:any_port_or_string
- port-<port_no>:any_port_or_string

### Configuration via environment variables
| Environment variable | Default value | Description          |
| -------------------- | ------------- | -------------------- |
| USER_PASSWORD        | 12345         | Password for labuser |
| ENABLE_WEB_TERMINAL  | 1             | Enable web terminal  |
| ENABLE_CODE_SERVER   | 1             | Enable code server   |
| ENABLE_VNC           | 1             | Enable VNC           |
| ENABLE_PORT_PROXY    | 1             | Enable port proxy    |

### Restrictions
- In the lab image, do not run any service on these ports - 80, 8001, 8002, 8003, 8004

### License
Apache License 2.0

### Credits
Special thanks to this project - https://github.com/ConSol/docker-headless-vnc-container (Apache License 2.0). It made it possible to have proper configuration for Xfce4 + noVNC setup.

### Author
Tanmoy Sarkar [@tanmoysrt](https://github.com/tanmoysrt)
Contact at tanmoy@swiftwave.org
