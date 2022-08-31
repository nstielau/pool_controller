import time

import RPi.GPIO as GPIO

from screenlogic import screenlogic

BLINKER_PIN=8
LED_PIN=10
BUTTON_PIN=12

# Configure the PIN # 8
GPIO.setmode(GPIO.BOARD)
GPIO.setup(LED_PIN, GPIO.OUT)
GPIO.setup(BLINKER_PIN, GPIO.OUT)
GPIO.setup(BUTTON_PIN, GPIO.IN, pull_up_down=GPIO.PUD_UP)
GPIO.setwarnings(False)

# Blink Interval 
blink_interval = .5 #Time interval in Seconds

bridge = screenlogic.slBridge(True)

# Blinker Loop
while True:
  #GPIO.output(BLINKER_PIN, True)
  #time.sleep(blink_interval)
  #GPIO.output(BLINKER_PIN, False)
  #time.sleep(blink_interval)

  input_state = GPIO.input(BUTTON_PIN)
  if input_state == False:
    print('Button Pressed at ' + time.ctime())
    GPIO.output(LED_PIN, True)
    time.sleep(0.2)
    current_value = bridge.getCircuit(502)
    print('current value is ' + current_value)
    new_value =  254-(~int(current_value) & 255)
    print('New value is ' + new_value)    
    bridge.setCircuit(502, new_value)
  else:
    GPIO.output(LED_PIN, False)

# Release Resources
GPIO.cleanup()
