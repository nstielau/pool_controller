import RPi.GPIO as GPIO
import time

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

# Blinker Loop
while True:
  GPIO.output(BLINKER_PIN, True)
  time.sleep(blink_interval)
  GPIO.output(BLINKER_PIN, False)
  time.sleep(blink_interval)

  input_state = GPIO.input(BUTTON_PIN)
  if input_state == False:
    print('Button Pressed')
    GPIO.output(LED_PIN, True)
    time.sleep(0.2)
  else:
    GPIO.output(LED_PIN, False)


# Release Resources
GPIO.cleanup()

