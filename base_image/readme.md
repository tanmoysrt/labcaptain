!!! This is an experimental docker image and still under development. Use at your own risk.

**setup**
- ubuntu 22.04
- root and non-root user : labuser / 12345
- noVNC server + xfce ui
- code server
- ttyd (web terminal)

**folder structure**
- /home/labuser/
  - .config/
    - code-server/
    - xfce4/
  - Desktop/
    - lab/ # put your lab related files here

**current support**
- web terminal at port 8001 : http://localhost:8001/?folder=/home/labuser/Desktop/lab
- code server at port 8002 : http://localhost:8002/?folder=/home/labuser/Desktop/lab
- vnc at port 8003 : http://localhost:8003/?resize=remote&autoconnect=1

**issues**
- dont auto lock vnc screen on idle, disable power management completely
- add a proxy inside the container, domain mapping will be there like this
  - http://terminal -> localhost:8001
  - http://editor -> localhost:8002
  - http://vnc -> localhost:8003
  - http://port-<port_no> -> localhost:<port_no>

**development**
```bash
docker build -t test_img .
docker run --rm -it -p 8001:8001 -p 8002:8002 -p 8003:8003 test_img
```