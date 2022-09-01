import logging

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

# Implement request handlers, exception handlers, etc.
# Register the handlers to the skill builder instance.

@sb.request_handler(can_handle_func=is_request_type("LaunchRequest"))
def launch_request_handler(handler_input):
    # type: (HandlerInput) -> Response
    speech_text = "Welcome to the Alexa Skills Kit, you can say hello!"

    handler_input.response_builder.speak(speech_text).set_card(
        SimpleCard("Hello World", speech_text)).set_should_end_session(
        False)
    return handler_input.response_builder.response

@sb.request_handler(can_handle_func=is_intent_name("HelloWorldIntent"))
def hello_world_intent_handler(handler_input):
    # type: (HandlerInput) -> Response
    speech_text = "Hello World!"

    handler_input.response_builder.speak(speech_text).set_card(
        SimpleCard("Hello World", speech_text)).set_should_end_session(
        True)
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
    body = request.body.read()
    headers = request.headers
    print(headers)
    print(body)
    print()
    # Convert the response str into web service format and return.
    return webservice_handler.verify_request_and_dispatch(headers, body)

run(host='0.0.0.0', port=80)
