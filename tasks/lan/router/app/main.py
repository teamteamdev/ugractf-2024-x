import aiohttp
import json
import re
import subprocess
import time
from quart import Quart, render_template, redirect, request, session, url_for

USERNAME = "admin"
PASSWORD = "xbYQsHUjj6hztTyzttkk"
MAX_RULES = 1000


app = Quart(__name__)
app.config.from_prefixed_env()

clients = {}
n_rules = 1

time_started = time.time()


@app.route("/")
async def index():
    return await render_template("login.html")


@app.route("/login", methods=["POST"])
async def login():
    form = await request.form
    username = form.get("username")
    password = form.get("password")
    if username == USERNAME and password == PASSWORD:
        session["authorized"] = True
        return redirect(url_for("status"))
    else:
        return await render_template("login.html", wrong_password=True)


@app.route("/console")
async def status():
    if not session.get("authorized"):
        return redirect(url_for("index"))

    async with aiohttp.ClientSession() as aio_session:
        async with aio_session.get("https://api.ipify.org") as response:
            external_ip = await response.text()

        start = time.time()
        async with aio_session.get("http://1.1.1.1") as response:
            end = time.time()

    return await render_template(
        "console.html",
        page="status",
        external_ip=external_ip,
        ping_ms=int((end - start) * 1000),
        clients=clients,
        uptime_min=int(time.time() - time_started) // 60
    )


@app.route("/console/configuration")
async def configuration():
    if not session.get("authorized"):
        return redirect(url_for("index"))

    proc = subprocess.run(["iptables", "-t", "nat", "-S", "PREROUTING"], capture_output=True, check=True)
    forwarded_ports = []
    for line in proc.stdout.decode().splitlines():
        match = re.fullmatch(r"-A PREROUTING -i tun0 -p tcp -m tcp --dport (\d+) -j DNAT --to-destination ([0-9.]+):(\d+)", line)
        if match:
            external_port = match[1]
            device_ip = match[2]
            internal_port = match[3]
            forwarded_ports.append((external_port, device_ip, internal_port))

    return await render_template(
        "console.html",
        page="configuration",
        clients=clients,
        forwarded_ports=forwarded_ports
    )


@app.route("/console/add-tcp-forwarding", methods=["POST"])
async def add_tcp_forwarding():
    if not session.get("authorized"):
        return redirect(url_for("index"))

    form = await request.form
    external_port = form["external_port"]
    device_ip = form["device_ip"]
    internal_port = form["internal_port"]

    assert 1 <= int(external_port) <= 65535
    assert re.fullmatch(r"[0-9.]+", device_ip)
    assert 1 <= int(internal_port) <= 65535

    global n_rules
    if n_rules < MAX_RULES:
        subprocess.run(["iptables", "-t", "nat", "-A", "PREROUTING", "-i", "tun0", "-p", "tcp", "-m", "tcp", "--dport", external_port, "-j", "DNAT", "--to-destination", f"{device_ip}:{internal_port}"], check=True)
        n_rules += 1
    else:
        return "Too many rules!", 429

    return redirect(url_for("configuration"))


@app.route("/console/delete-tcp-forwarding", methods=["POST"])
async def delete_tcp_forwarding():
    if not session.get("authorized"):
        return redirect(url_for("index"))

    form = await request.form
    external_port = form["external_port"]
    device_ip = form["device_ip"]
    internal_port = form["internal_port"]

    assert 1 <= int(external_port) <= 65535
    assert re.fullmatch(r"[0-9.]+", device_ip)
    assert 1 <= int(internal_port) <= 65535

    global n_rules
    subprocess.run(["iptables", "-t", "nat", "-D", "PREROUTING", "-i", "tun0", "-p", "tcp", "-m", "tcp", "--dport", external_port, "-j", "DNAT", "--to-destination", f"{device_ip}:{internal_port}"], check=True)
    n_rules -= 1

    return redirect(url_for("configuration"))


@app.route("/new-client", methods=["POST"])
async def new_client():
    form = await request.form
    ip = request.remote_addr
    name = form["name"]
    clients[ip] = name
    return ""



@app.route("/reset-tables")
async def reset_tables():
    global n_rules
    n_rules = 1
    return ""


if __name__ == "__main__":
    app.run("0.0.0.0")
