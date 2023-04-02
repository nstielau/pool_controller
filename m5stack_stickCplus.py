from m5stack import *
from m5ui import *
from uiflow import *

import urequests

import json
import wifiCfg

wifiCfg.doConnect('SSID', 'passsword')

setScreenColor(0x111111)

status = M5TextBox(10, 220, "", lcd.FONT_Default, 0xFFFFFF, rotate=0)
temp_lbl = M5TextBox(30, 20, "??F", lcd.FONT_DejaVu40, 0xFFFFFF, rotate=0)
lbl = M5TextBox(75, 80, "HOTTUB", lcd.FONT_DejaVu24, 0xFFFFFF, rotate=90)

from numbers import Number

while True:
  if btnA.isPressed():
    status.setText("Loading...")
    req = urequests.request(
      method='GET',
      url='https://myFQDN.com/pool')
    pool_data = json.loads(req.text)
    spa_on = pool_data["spa_heater_mode"]["state"]
    if spa_on:
      spa_temp = pool_data["current_spa_temperature"]["state"].split()[0]
      temp_lbl.setText(str(spa_temp) + "F")
    status.setText("")
  wait_ms(2)
