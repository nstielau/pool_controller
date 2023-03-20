import time
import sys

import RPi.GPIO as GPIO

from screenlogic.slBridge import slBridge

LED_PIN=10
BUTTON_PIN=12

SWIM_JET_CIRCUIT=502 # 502 is swim jet, 503 is lights
# 500 - Spa
# 501 - Cleaner
# 502 - Swim Jets
# 503 - Pool Light
# 504 - Spa Light
# 505 - Pool
# 506 - Aux 5
# 507 - Aux 6
# 508 - Aux 7

# Configure the board
GPIO.setmode(GPIO.BOARD) # Use Board numbers https://pinout.xyz/
GPIO.setup(LED_PIN, GPIO.OUT)
GPIO.setup(BUTTON_PIN, GPIO.IN, pull_up_down=GPIO.PUD_UP)
GPIO.setwarnings(False)

# Delay to avoid intitial false-positive button press
time.sleep(2)

# Loop!
while True:
  input_state = GPIO.input(BUTTON_PIN)
  if input_state == False: # pressed!
    print('Button Pressed at ' + time.ctime())
    GPIO.output(LED_PIN, True)
    time.sleep(0.2)

    bridge = slBridge(True) # bridge update logic seems broke
    current_value = bridge.getCircuit(SWIM_JET_CIRCUIT)
    print('current value is ' + current_value)
    new_value = 0 if current_value == "On" else 1
    print('New value is ' + str(new_value))    
    bridge.setCircuit(SWIM_JET_CIRCUIT, new_value)
  else:
    GPIO.output(LED_PIN, False)
