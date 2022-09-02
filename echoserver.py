import logging

from screenlogic.screenlogic import slBridge    

from bottle import route, run, template, request

from ask_sdk_core.skill_builder import SkillBuilder
from ask_sdk_webservice_support.webservice_handler import WebserviceSkillHandler
from ask_sdk_core.utils import is_request_type, is_intent_name
from ask_sdk_core.handler_input import HandlerInput
from ask_sdk_model import Response
from ask_sdk_model.ui import SimpleCard

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

#####################################################################
#####################################################################
# Handlers
#####################################################################

skill_builder = SkillBuilder()

@skill_builder.request_handler(can_handle_func=is_request_type("LaunchRequest"))
def launch_request_handler(handler_input):
    speech_text = "Pool party time"

    handler_input.response_builder.speak(speech_text).set_card(
        SimpleCard(speech_text, speech_text)).set_should_end_session(False)
    return handler_input.response_builder.response

@skill_builder.request_handler(can_handle_func=is_intent_name("StartSwimJetIntent"))
def start_swimjet_intent_handler(handler_input):
    speech_text = "Pool jet started"

    slBridge(True).setCircuit(502, 1)

    handler_input.response_builder.speak(speech_text).set_card(
        SimpleCard(speech_text, speech_text)).set_should_end_session(True)
    return handler_input.response_builder.response

@skill_builder.request_handler(can_handle_func=is_intent_name("StopSwimJetIntent"))
def stop_swimjet_intent_handler(handler_input):
    speech_text = "Pool jet stopped"

    slBridge(True).setCircuit(502, 0)

    handler_input.response_builder.speak(speech_text).set_card(
        SimpleCard(speech_text, speech_text)).set_should_end_session(True)
    return handler_input.response_builder.response

webservice_handler = WebserviceSkillHandler(skill=skill_builder.create())

#####################################################################
#####################################################################
# Webserver
#####################################################################
@route('/', method=['GET'])
def index():
    return "hello"

@route('/', method=['POST'])
def index():
    body = request.body.read().decode()
    headers = request.headers
    print(headers)
    print(body)
    print()
    return webservice_handler.verify_request_and_dispatch(headers, body)

run(host='0.0.0.0', port=80)
