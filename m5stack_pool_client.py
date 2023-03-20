from m5stack import *
from m5ui import *
from uiflow import *
import urequests
import time
import json

setScreenColor(0x222222)

fqdn = "3da840e05d58dbccc1a52769e568811f.balena-devices.com"

# Describe this function...
def getPoolData():
  lcd.clear()
  
  label0 = M5TextBox(69, 67, "Status", lcd.FONT_DejaVu24, 0xFFFFFF, rotate=0)
  
  title = M5Title(title="Pool Party", x=3, fgcolor=0xFFFFFF, bgcolor=0x0000FF)
  status = M5TextBox(220, 2, "status", lcd.FONT_Default, 0xFFFFFF, rotate=0)
  pool_rct = M5Rect(205, 54, 101, 150, 0x222222, 0xFFFFFF)
  spa_rct = M5Rect(205, 54, 50, 50, 0x222222, 0xFFFFFF)
  spa_lbl = M5TextBox(216, 77, "", lcd.FONT_Default, 0x000000, rotate=0)
  btn_a_lbl = M5TextBox(35, 220, "Refresh", lcd.FONT_Default, 0xFFFFFF, rotate=0)  
  
  status.setText('Loading')
  try:
    req = urequests.request(method='GET', url='https://' + fqdn + '/pool', headers={})
    pool_data = json.loads(req.text)
    spa_on = pool_data["spa_heater_mode"]["state"]
    if spa_on:
      spa_rct = M5Rect(205, 54, 50, 50, 0xFF0000, 0x000000)
      spa_temp = pool_data["current_spa_temperature"]["state"].split()[0]
      spa_lbl = M5TextBox(216, 77, spa_temp + " F", lcd.FONT_Default, 0x000000, rotate=0)
    else:
      spa_rct = M5Rect(205, 54, 50, 50, 0x0000FF, 0x000000)
      spa_lbl = M5TextBox(216, 77, "Off", lcd.FONT_Default, 0x000000, rotate=0)
    status.setText('Connected')    
  except:
    label1.setText(str(req.text))
    status.setText('Error')

btnA.wasPressed(getPoolData)


spa_temp = 0
getPoolData()
while True:
  wait(300)
  getPoolData()

