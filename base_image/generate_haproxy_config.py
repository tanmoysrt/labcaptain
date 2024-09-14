haproxy_cfg_template = ""
with open("haproxy.cfg.template", "r") as f:
    haproxy_cfg_template = f.read()

ports = range(1, 65535)

res = "\n".join([f"    server port_{i} 127.0.0.1:{i}" for i in ports])
res += "\n\n"
res += "\n".join([f"    use-server port_{i} if {{ var(txn.only_host) -m str port-{i} }}" for i in ports])

with open("haproxy.cfg", "w") as f:
    f.write(haproxy_cfg_template.replace("<insert_port_rules>", res))