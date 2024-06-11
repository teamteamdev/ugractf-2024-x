import os
import time
import requests

# Prepare:
# pactl load-module module-null-sink sink_name="virtual_speaker" sink_properties=device.description="virtual_speaker"
# pactl load-module module-remap-source master="virtual_speaker.monitor" source_name="virtual_mic" source_properties=device.description="virtual_mic"
# Then set mic in Zoom

while True:
  print("...", flush=True)
  r = requests.get(f"https://{os.environ['DOMAIN']}/sound_trigger.txt")
  if r.status_code == 200:
    os.execve(
      "/usr/bin/ffmpeg",
      ["/usr/bin/ffmpeg", "-i", "~/sound.wav", "-f", "pulse", "stream name"],
      {"PULSE_SINK": "virtual_speaker"}
    )
  time.sleep(15)
