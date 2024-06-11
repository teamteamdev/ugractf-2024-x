import json
import time
from quart import Quart, render_template, request

app = Quart(__name__)


FLAG = "we_re_not_strangers_to_love"

NOTES = [
    "B3", "C#4", "D4", "D4", "E4", "C#4", "B3", "A3",
    "B3", "B3", "C#4", "D4", "B3", "A3", "A4", "A4", "E4",
    "B3", "B3", "C#4", "D4", "B3", "D4", "E4", "C#4", "B3", "C#4", "B3", "A3",
    "B3", "B3", "C#4", "D4", "B3", "A3", "E4", "E4", "E4", "F#4", "E4",

    "D4", "E4", "F#4", "D4", "E4", "E4", "E4", "F#4", "E4", "A3",
    "B3", "C#4", "D4", "B3", "E4", "F#4", "E4",

    "A3", "B3", "D4", "B3", "F#4", "F#4", "E4",
    "A3", "B3", "D4", "B3", "E4", "E4", "D4", "C#4", "B3",
    "A3", "B3", "D4", "B3", "D4", "E4", "C#4", "A3", "A3", "E4", "D4",

    "A3", "B3", "D4", "B3", "F#4", "F#4", "E4",
    "A3", "B3", "D4", "B3", "A4", "C#4", "D4", "C#4", "B3",
    "A3", "B3", "D4", "B3", "D4", "E4", "C#4", "A3", "A3", "E4", "D4"
]


@app.route("/")
async def index():
    return await render_template("index.html")


progression_by_flow_id = {}
flag_time_by_flow_id = {}
next_flow_id = 0

@app.route("/api/start", methods=["POST"])
async def start():
    global next_flow_id
    flow_id = next_flow_id
    next_flow_id += 1
    progression_by_flow_id[flow_id] = 0
    return {"flowId": flow_id}

@app.route("/api/enter/<flow_id>", methods=["POST"])
async def enter(flow_id):
    flow_id = int(flow_id)
    note = (await request.data).decode()

    progression = progression_by_flow_id[flow_id]
    ok = NOTES[progression] == note
    code = None

    if ok:
        progression_by_flow_id[flow_id] = progression + 1
        if progression_by_flow_id[flow_id] == len(NOTES):
            code = "document.write(" + json.dumps(open("templates/congratulations.html").read()) + ")"
            flag_time_by_flow_id[flow_id] = time.time() + 3 * 60 + 33

    return {"ok": ok, "code": code}

@app.route("/api/flag/<flow_id>")
async def flag(flow_id):
    flow_id = int(flow_id)

    if time.time() < flag_time_by_flow_id[flow_id]:
        return "Еще рано, ждите"
    else:
        return f"Флаг: {FLAG}"


if __name__ == "__main__":
    app.run("0.0.0.0")
