import os
from quart import Quart, render_template, redirect, request, session
import sqlite3

conn = sqlite3.connect('/nya/a.db', autocommit=True)

conn.cursor().execute("DROP TABLE IF EXISTS accounts")
conn.cursor().execute("CREATE TABLE accounts (login TEXT, password TEXT)")
conn.cursor().execute("CREATE TABLE IF NOT EXISTS donosy (who TEXT, whom TEXT)")

for user, password in [
    ('Калан', 'tn3895ny'),
    ('Сергей', '12345678'),
    ('Екаетрина', 'fjjjjfds'),
    ('Владимир', 'p4r11jjj'),
    ('Александр', 't1murrk4'),
    ('Марина', 'm4sha819'),
]:
    conn.cursor().execute("INSERT INTO accounts VALUES (?, ?)", (user, password))

import random
votable_in = 'Калан Сергей Владимир Александр Марина'.split()
votable_out = 'Калан Сергей Екаетрина Владимир Александр Марина'.split()
for uin in votable_in:
    for uout in votable_out:
        if random.choice(["yes", "no"]) == "yes":
            conn.cursor().execute("INSERT INTO donosy VALUES (?, ?)", (uin, uout))


app = Quart(__name__)
app.config.from_prefixed_env()
app.secret_key = 'tnsuyragmwpymkgtrmagsrntiroetnrsiotnrestulzfpyugmuyvnmrset'


@app.route("/login", methods=["POST"])
async def login():
    form = await request.form
    username = form.get("username")
    password = form.get("password")
    if conn.cursor().execute("SELECT * FROM accounts WHERE login = '%s' AND password = '%s'" % (username, password)).fetchone():
        session["authorized"] = username
        return redirect("/")
    else:
        return await render_template("login.html", wrong_password=True)


@app.route("/vote", methods=["POST"])
async def vote():
    user = session.get("authorized")
    if not user:
        return await render_template("login.html")

    conn.cursor().execute("INSERT INTO donosy VALUES (?, ?)", (user, (await request.form).get("whom")))

    return redirect("/")


@app.route("/")
async def status():
    user = session.get("authorized")
    if not user:
        return await render_template("login.html")

    (left,) = conn.cursor().execute("SELECT COUNT(*) FROM donosy WHERE who = ?", (user,)).fetchone()
    (right,) = conn.cursor().execute("SELECT COUNT(*) FROM donosy WHERE whom = ?", (user,)).fetchone()

    kek = []
    for (nya,) in conn.cursor().execute("SELECT login FROM accounts"):
        kek.append(nya)

    return await render_template(
        "console.html",
        nya=left - right,
        logins=kek,
    )


if __name__ == "__main__":
    app.run("0.0.0.0", 80)
