import time
import sys

from screenlogic.screenlogic import slBridge

LED_PIN=10
BUTTON_PIN=12

SWIM_JET_CIRCUIT=502

bridge = slBridge(True)

if(len(sys.argv) > 1): # look for any any, e.g. "run"
  import RPi.GPIO as GPIO

  # Configure the board
  GPIO.setmode(GPIO.BOARD) # Use Board numbers https://pinout.xyz/
  GPIO.setup(LED_PIN, GPIO.OUT)
  GPIO.setup(BUTTON_PIN, GPIO.IN, pull_up_down=GPIO.PUD_UP)
  GPIO.setwarnings(False)

  # Loop!
  while True:
    input_state = GPIO.input(BUTTON_PIN)
    if input_state == False: # pressed!
      print('Button Pressed at ' + time.ctime())
      GPIO.output(LED_PIN, True)
      time.sleep(0.2)
      current_value = bridge.getCircuit(SWIM_JET_CIRCUIT)
      print('current value is ' + current_value)
      new_value = 0 if current_value == "On" else 1
      print('New value is ' + new_value)    
      bridge.setCircuit(502, new_value)
    else:
      GPIO.output(LED_PIN, False)
