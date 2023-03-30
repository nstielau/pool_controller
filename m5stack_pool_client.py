import json
import time

import urequests
import wifiCfg

from m5stack import *
from m5ui import *
from uiflow import *

wifiCfg.doConnect('SSID', 'password')

# Nice dark grey
setScreenColor(0x222222)

# Set initial vars
# nvs.write(str('TOKEN'), 'MYTOKEN')
# nvs.write(str('FQDN'), 'MYFQDN')
# str(nvs.read_str('TOKEN'))

# FQDN for pool interface
fqdn = "fq.dn"

spa_on = False
spa_temp = 0

# Query pool dta from server
def getPoolData():
  global spa_on
  global spa_temp
  title = M5Title(title="Hot Tub Time Machine", x=3, fgcolor=0xFFFFFF, bgcolor=0x0000FF)
  status = M5TextBox(220, 4, "Loading", lcd.FONT_Default, 0xFFFFFF, rotate=0)
  try:
    req = urequests.request(
      method='GET',
      url='https://' + fqdn + '/pool')
    pool_data = json.loads(req.text)
    lcd.clear()
    title = M5Title(title="Hot Tub Time Machine", x=3, fgcolor=0xFFFFFF, bgcolor=0x0000FF)
    spa_on = pool_data["spa_heater_mode"]["state"]
    if spa_on:
      spa_temp = pool_data["current_spa_temperature"]["state"].split()[0]
      btn_a_lbl = M5TextBox(35, 220, "Refresh", lcd.FONT_Default, 0xFFFFFF, rotate=0)
      label0 = M5TextBox(100, 70, spa_temp + " F", lcd.FONT_DejaVu24, 0xFFFFFF, rotate=0)
      thermTop = M5Rect(180, 40, 40, 150, 0xFF0000, 0xFF0000)
      thermBall = M5Circle(200, 190, 30, 0xFF0000, 0xFF0000)
      height = int(150/104*int(spa_temp))
      thermTop2 = M5Rect(180, 40, 40, 150-height, 0x0000FF, 0x0000FF)
    else:
      title = M5Title(title="Hot Tub Time Machine", x=3, fgcolor=0xFFFFFF, bgcolor=0x0000FF)
      btn_a_lbl = M5TextBox(35, 220, "Refresh", lcd.FONT_Default, 0xFFFFFF, rotate=0)
      spa_lbl = M5TextBox(100, 70, "Off", lcd.FONT_DejaVu24, 0xFFFFFF, rotate=0)
      thermTop = M5Rect(180, 40, 40, 150, 0x0000FF, 0x0000FF)
      thermBall = M5Circle(200, 190, 30, 0x0000FF, 0x0000FF)
      thermTop2 = M5Rect(180, 40, 40, 30, 0x0000FF, 0x0000FF)
  except Exception as e:
    status.setText('Error')
    M5TextBox(20, 24, "FQDN: " + fqdn, lcd.FONT_Default, 0xFFFFFF, rotate=0)
    M5TextBox(20, 64, "Error: " + str(e), lcd.FONT_Default, 0xFFFFFF, rotate=0)
    M5TextBox(20, 104, "Type: " + str(e.__class__), lcd.FONT_Default, 0xFFFFFF, rotate=0)
    raise e

# Button listeners
btnA.wasPressed(getPoolData)

# Main loop
wait(3)

while True:
  getPoolData()
  if spa_on and int(spa_temp) > 97:
    for s in range(150):
      R = 255
      G = 0
      B = 0
      for i in range(256):
        rgb.setColorFrom(1, 10, (R << 16) | (G << 8) | B)
        rgb.setBrightness(i)
        wait_ms(2)
      for i in range(255, -1, -1):
        rgb.setColorFrom(1, 10, (R << 16) | (G << 8) | B)
        rgb.setBrightness(i)
        wait_ms(2)
  else:
    wait(600)


