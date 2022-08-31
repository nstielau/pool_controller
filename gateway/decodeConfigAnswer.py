import struct
from gateway.decodeData import getSome, getString

def decodeConfigAnswer(buff, data):

  #{ name="", state= }
  if('config' not in data):
    data['config'] = {}
    
  controlerID, offset = getSome("I", buff, 0)
  data['config']['controler_id'] = dict(name="Controler ID", state=controlerID)

  minSetPoint0, offset = getSome("B", buff, offset)
  maxSetPoint0, offset = getSome("B", buff, offset)
  minSetPoint1, offset = getSome("B", buff, offset)
  maxSetPoint1, offset = getSome("B", buff, offset)

  #if('bodies' not in data):
  #  data['bodies'] = {}
    
  data['config']['min_set_point'] = dict(name="Minimum Temperature", state=[minSetPoint0, minSetPoint1])
  data['config']['max_set_point'] = dict(name="Maximum Temperature", state=[maxSetPoint0, maxSetPoint1])

  degC, offset = getSome("B", buff, offset)
  data['config']['is_celcius'] = dict(name="Is Celcius", state=degC)
  
  controllerType, offset = getSome("B", buff, offset)
  data['config']['controler_type'] = controllerType
  
  hwType, offset = getSome("B", buff, offset)
  data['config']['hardware_type'] = hwType

  controllerbuff, offset = getSome("B", buff, offset)
  data['config']['controler_buffer'] = controllerbuff

  equipFlags, offset = getSome("i", buff, offset)
  data['config']['equipment_flags'] = equipFlags

  paddedGenName, offset = getString(buff, offset)
  genCircuitName = paddedGenName.decode("utf-8").strip('\0')
  data['config']['generic_circuit_name'] = dict(name="Default Circuit Name", state=genCircuitName)
         
  circuitCount , offset = getSome("I", buff, offset)
  data['config']['circuit_count'] = dict(name="Number of Circuits", state=circuitCount)
         
  if('circuits' not in data):
    data['circuits'] = {}
     
  for i in range(circuitCount):

    circuitID, offset = getSome("i", buff, offset)

    if(circuitID not in data['circuits']):
      data['circuits'][circuitID] = {}
      
    data['circuits'][circuitID]['id'] = circuitID
    
    paddedName, offset = getString(buff, offset)
    circuitName = paddedName.decode("utf-8").strip('\0')
    data['circuits'][circuitID]['name'] = circuitName
    
    cNameIndex, offset = getSome("B", buff, offset)
    data['circuits'][circuitID]['name_index'] = cNameIndex
    
    cFunction, offset = getSome("B", buff, offset)
    data['circuits'][circuitID]['function'] = cFunction
    
    cInterface, offset = getSome("B", buff, offset)
    data['circuits'][circuitID]['interface'] = cInterface
    
    cFlags, offset = getSome("B", buff, offset)
    data['circuits'][circuitID]['flags'] = cFlags
    
    cColorSet, offset = getSome("B", buff, offset)
    data['circuits'][circuitID]['color_set'] = cColorSet

    cColorPos, offset = getSome("B", buff, offset)
    data['circuits'][circuitID]['color_position'] = cColorPos

    cColorStagger, offset = getSome("B", buff, offset)
    data['circuits'][circuitID]['color_stagger'] = cColorStagger

    cDeviceID, offset = getSome("B", buff, offset)
    data['circuits'][circuitID]['device_id'] = cDeviceID

    cDefaultRT, offset = getSome("H", buff, offset)
    data['circuits'][circuitID]['default_rt'] = cDefaultRT

    offset = offset + struct.calcsize("2B")

  colorCount , offset = getSome("I", buff, offset)
  data['config']['color_count'] = dict(name="Number of Colors", state=colorCount)
  
  if('colors' not in data['config'] or len(data['config']['colors']) != colorCount):
    data['config']['colors'] = [{} for x in range(colorCount)]
  
  for i in range(colorCount):
    paddedColorName, offset = getString(buff, offset)
    colorName = paddedColorName.decode("utf-8").strip('\0')
    rgbR, offset = getSome("I", buff, offset)
    rgbG, offset = getSome("I", buff, offset)
    rgbB, offset = getSome("I", buff, offset)
    data['config']['colors'][i] = dict(name=colorName, state=[ rgbR, rgbG, rgbB])
    
  pumpCircuitCount = 8
  if('pumps' not in data['config']):
    data['config']['pumps'] = {}
    
  for i in range(pumpCircuitCount):
    if(i not in data['config']['pumps']):
      data['config']['pumps'][i] = {}
      
    pumpData , offset = getSome("B", buff, offset)
    data['config']['pumps'][i]['data'] = pumpData
    
  interfaceTabFlags , offset = getSome("I", buff, offset)
  data['config']['interface_tab_flags'] = interfaceTabFlags
  
  showAlarms , offset = getSome("I", buff, offset)
  data['config']['show_alarms'] = showAlarms
  
