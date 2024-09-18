#!/bin/bash

# remove the quadlet file
rm -f /etc/containers/systemd/{{lab_id}}.container >/dev/null 2>&1 || true
# stop the service
systemctl stop {{lab_id}}.service >/dev/null 2>&1 || true
# do a daemon-reload
systemctl daemon-reload
