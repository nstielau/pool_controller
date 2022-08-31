import socket
import struct
from gateway.messageHelper import makeMessageString, getMessageString, makeMessage, decodeMessage
from gateway.constants import *

# returns a formatted login message
def createLoginMessage():
  # these constants are only for this message. keep them here.
  schema = 348
  connectionType = 0
  clientVersion = makeMessageString('Android')
  pid  = 2
  password = "mypassword" # passwd must be <= 16 chars. empty is not OK.
  passwd = makeMessageString(password)
  fmt = "<II" + str(len(clientVersion)) + "s" + str(len(passwd)) + "sxI"
  return struct.pack(fmt, schema, connectionType, clientVersion, passwd, pid)

def gatewayLogin(gatewayIP, gatewayPort):
  # cut/paste from python manual
  tcpSock = None
  for res in socket.getaddrinfo(gatewayIP, gatewayPort, socket.AF_UNSPEC, socket.SOCK_STREAM):
    af, socktype, proto, canonname, sa = res
    try:
      tcpSock = socket.socket(af, socktype, proto)
    except OSError as msg:
      tcpSock = None
      continue
    try:
      tcpSock.connect(sa)
    except OSError as msg:
      tcpSock.close()
      tcpSock = None
      continue
    break

  if tcpSock is None:
    sys.stderr.write("ERROR: {}: Could not open socket to gateway host.\n".format(me))
    sys.exit(10)

  #with tcpSock:
  # get the gateway's attention. The Protocol_Document.pdf explains how, not why.
  connectString = b'CONNECTSERVERHOST\r\n\r\n'  # not a string...
  tcpSock.sendall(connectString)
  # the gateway does not respond to the connect message. don't wait for something here because you aren't going to get it

  # tx/rx challenge  (?)  (gateway returns its mac address in the form 01-23-45-AB-CD-EF)
  # why? dunno.
  tcpSock.sendall(makeMessage(code.CHALLENGE_QUERY))
  data = tcpSock.recv(48)
  if not data:
    sys.stderr.write("WARNING: {}: no {} data received.\n".format(me, "CHALLENGE_ANSWER"))
  rcvcode, data = decodeMessage(data)
  if(rcvcode != code.CHALLENGE_ANSWER):
    sys.stderr.write("WARNING: {}: rcvCode2({}) != {}.\n".format(me, CHALLENGE_ANSWER))
    sys.exit(10)


  # now that we've "connected" and "challenged," we can "login." None of these things
  # actually do anything, but they are required.
  msg = createLoginMessage()
  tcpSock.sendall(makeMessage(code.LOCALLOGIN_QUERY, msg))
  data = tcpSock.recv(48)
  if not data:
    sys.stderr.write("WARNING: {}: no {} data received.\n".format(me, "LOCALLOGIN_ANSWER"))
  rcvCode, data = decodeMessage(data)
  if(rcvCode != code.LOCALLOGIN_ANSWER):
    sys.stderr.write("WARNING: {}: rcvCode({}) != {}.\n".format(me, rcvCode, code.LOCALLOGIN_ANSWER))
    sys.exit(10)
  # response should be empty
  return tcpSock
