import json
import math
import time
from quart import Quart, render_template, request

app = Quart(__name__)


FLAG = "this_is_not_an_exit"


@app.route("/")
async def index():
    return await render_template("index.html")


@app.route("/enter", methods=["POST"])
async def enter():
    story = await request.json
    handle_touch_id = None
    handle_radius = None
    angle = 0
    for event, *args in story:
        if event == "start":
            touch_id, node_id, dx, dy = args
            if node_id == "handle":
                handle_touch_id = touch_id
                handle_radius = (dx * dx + dy * dy) ** 0.5
                angle = 0
        elif event == "move":
            touch_id, dx, dy = args
            if touch_id == handle_touch_id:
                radius = (dx * dx + dy * dy) ** 0.5
                if handle_radius is None or dx < 0 or abs(radius - handle_radius) > 8:
                    angle = 0
                else:
                    angle = min(max(-0.1, math.atan2(dy, dx)), 1.2)
        elif event == "stop":
            touch_id, = args
            if touch_id == handle_touch_id:
                angle = 0
                handle_touch_id = None
        elif event == "door":
            if angle == 1.2:
                return f"Флаг: {FLAG}"
    return "Дверь закрыта."


if __name__ == "__main__":
    app.run("0.0.0.0")
