import struct
#import numpy as np

def getSome(want, buff, offset):
  fmt = "<" + want
  newoffset = offset + struct.calcsize(fmt)
  return struct.unpack_from(fmt, buff, offset)[0], newoffset

def decodeButtonPressAnswer(data):

  state = {}
  
  ok, offset = getSome("I", data, 0)
  print("ok: {}".format(ok))

  remainder, offset = getSome("I", data, offset)
  print("remainder: {}".format(remainder))
  #state['chemistry']['alarms'] = alarms

  return state
