#!/usr/bin/env python3

import pathlib
import sys

# Hack to include modules from project root
sys.path.append(str(pathlib.Path(__file__).parent.parent.resolve()))

from screenlogic.slBridge import slBridge
        
if __name__ == "__main__":
    bridge = slBridge(True)
    if(len(sys.argv) > 1):
        if(sys.argv[1] == 'get'):
            if(len(sys.argv) == 3):
                print(bridge.getCircuit(int(sys.argv[2])))
        elif(sys.argv[1] == 'set'):
            if(len(sys.argv) == 4):
                print(bridge.setCircuit(int(sys.argv[2]), int(sys.argv[3])))
        elif(sys.argv[1] == 'data'):
            print(bridge.dumpDict())
        elif(sys.argv[1] == 'json'):
            print(bridge.getJson())
        else:
            print("Unknown option!")
    else:
        print(bridge.getFriendly())#getJson())
