from gateway.constants import mapping

class slDevice:
    def __init__(self, slBridge, dataID, data):
        self._bridge = slBridge
        self._id = dataID
        self._name = data["name"]
        self._state = data["state"]
        self._hassType = data["hassType"]

    def update(self, data):
        self._state = data["state"]

    @property
    def id(self):
        return self._id

    @property
    def name(self):
        return self._name

    @property
    def state(self):
        return self._state

    @property
    def hassType(self):
        return self._hassType
