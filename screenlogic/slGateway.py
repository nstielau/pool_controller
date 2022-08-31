import socket
from gateway.gatewayLogin import gatewayLogin
from gateway.gatewayQuery import queryGateway, queryConfig, queryStatus, queryButtonPress

class slGateway:
    def __init__(self, ip, port):
        self.__ip = ip
        self.__port = port
        self.__connected = False

    def connect(self):
        self.__socket = gatewayLogin(self.__ip, self.__port)
        if(self.__socket):
            self.__version = ""
            self.__version = queryGateway(self.__socket)
            if(self.__version):
                self.__connected = True
                return True
        return False

    def disconnect(self):
        self.__socket.close()
        self.__connected = False

    def connected(self):
        return self.__connected

    def getConfig(self, data):
        if(self.__connected or self.connect()):
            queryConfig(self.__socket, data)

    def getStatus(self, data):
        if(self.__connected or self.connect()):
            queryStatus(self.__socket, data)

    def setCircuit(self, circuitID, circuitState):
        if(self.__connected or self.connect()):
            return queryButtonPress(self.__socket, circuitID, circuitState)
