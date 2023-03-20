#!/usr/bin/env python3

import json
import sys
import time

from multiprocessing import Lock
from gateway.gatewayDiscovery import discoverGateway
from screenlogic.slGateway import slGateway
from screenlogic.slSwitch import slSwitch
from screenlogic.slSensor import slSensor

class slBridge:
    def __init__(self, verbose=False, updateInterval=30, gatewayIP=None, gatewayPort=None):
        self.__lastUpdate = 0
        self.__updateInterval = updateInterval
        self.__lock = Lock()
        self.__data = {}
        self.__devices = {}
        
        if(not gatewayIP):
              gatewayIP, gatewayPort, gatewayType, gatewaySubtype, \
              gatewayName, okchk = discoverGateway(verbose)

        if(gatewayIP):
            self.__gateway = slGateway(gatewayIP, gatewayPort)
            if(self.__gateway.connect()):
                if(verbose):
                    print("Connection success!")
                self.__gateway.getConfig(self.__data)
                self.__gateway.getStatus(self.__data)
                self.__gateway.disconnect()
                self._updateDevices()
            else:
                if(verbose):
                    print("Connection failed!")
        else:
            print("Discovery Failed!")

    def update(self):
        curTime = time.time()
        with self.__lock:
            if ((curTime - self.__lastUpdate) > self.__updateInterval):
                self.__lastUpdate = curTime
                self._updateData()
                self._updateDevices()

    def _updateData(self):
        self.__gateway.getStatus(self.__data)

    def _updateDevices(self):
        self._updateSwitches()
        self._updateSensors()

    def _updateSwitches(self):
        for k, v in self.__data['circuits'].items():
            if('hassType' in v):
                if(k in self.__devices):
                    self.__devices[k].update(v)
                else:
                    self.__devices[k] = slSwitch(self, k, v)

    def _updateSensors(self):
        for k, v in self.__data['sensors'].items():
            if('hassType' in v):
                if(k in self.__devices):
                    self.__devices[k].update(v)
                else:
                    self.__devices[k] = slSensor(self, k, v)
        for i in self.__data['bodies']:
            for k, v in self.__data['bodies'][i].items():
                if('hassType' in v):
                    kI = "{}_{}".format(k, i)
                    if(kI in self.__devices):
                        self.__devices[kI].update(v)
                    else:
                        self.__devices[kI] = slSensor(self, kI, v)
        for k, v in self.__data['chemistry'].items():
            if('hassType' in v):
                if(k in self.__devices):
                    self.__devices[k].update(v)
                else:
                    self.__devices[k] = slSensor(self, k, v)


    def getDevices(self):
        return self.__devices


    def getJson(self):
        dictOut = {}
        for k, d in self.__devices.items():
            if(d.hassType == 'sensor'):
                dictData = dict(name=d.name,state=d.friendlyState)#state,unit=d.unit,friendly_state=d.friendlyState)
            else:
                dictData = {}
                dictData['id'] = k
                dictData['name'] = d.name
                dictData['state'] = self._jsonName(d.friendlyState)
            dictOut[self._jsonName(d.name)] = dictData 
        return json.dumps(dictOut)

    def _jsonName(self, name):
        #s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
        #return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()
        return name.replace(" ","_").lower()


    def getFriendly(self):
        self._updateDevices()
        for k, d in self.__devices.items():
            print("{} - {}: {}".format(k, d.name, d.friendlyState))

    def getConfig(self):
        return self.__data['config']

    def getChemistry(self):
        return self.__data['states']['chemistry']

    def getCircuit(self, circuitID):
        if circuitID in self.__devices:
            return self.__devices[circuitID].friendlyState
        else:
            return "error"
    
    def setCircuit(self, circuitID, circuitState):
        if(circuitID in self.__devices and self.__gateway.setCircuit(circuitID, circuitState)):
            self._updateData()
            return True
        
    def dumpDict(self):
        return json.dumps(self.__data)

    def getKeys(self):
        for k, d in self.__devices.items():
            print("      - {}".format(self._jsonName(d.name)))



        
if __name__ == "__main__":
    bridge = slBridge(True)
    if(len(sys.argv) > 1):
        if(sys.argv[1] == 'get'):
            if(len(sys.argv) == 3):
                print(bridge.getCircuit(int(sys.argv[2])))
        elif(sys.argv[1] == 'set'):
            if(len(sys.argv) == 4):
                print(bridge.setCircuit(int(sys.argv[2]), int(sys.argv[3])))
        elif(sys.argv[1] == 'json'):
            print(bridge.getJson())
        else:
            print("Unknown option!")
    else:
        print(bridge.getFriendly())#getJson())
