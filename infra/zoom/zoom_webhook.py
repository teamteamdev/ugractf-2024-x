import flask
import hmac
import os
import requests

app = flask.Flask(__name__)

SECRET_KEY = os.environ['SECRET_KEY']
BOT_TOKEN = os.environ['BOT_TOKEN']


@app.post('/ev2')
def endpoint():
    j = flask.request.json
    if j['event'] == 'endpoint.url_validation':
      return flask.jsonify({
           'plainToken': j['payload']['plainToken'],
           'encryptedToken': hmac.new(SECRET_KEY.encode(), j['payload']['plainToken'].encode(), 'sha256').hexdigest()
      })
    if j['event'] == 'meeting.participant_joined':
        user_name = j['payload']['object']['participant']['user_name']
        requests.post(f'https://api.telegram.org/bot{BOT_TOKEN}/sendMessage', json={"chat_id": -1001603104763, "text": f"User {user_name} joined Zoom Meeting!  Alarm! @nsychev"})
        with open("sound_trigger.txt", "w") as f:
            f.write("1")
    return ""


if __name__ == "__main__":
    app.run(port=3001)
