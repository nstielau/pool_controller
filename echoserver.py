from bottle import route, run, template, request

@route('/')
def index():
    print(request.body.read())
    print()
    return "hello"

run(host='localhost', port=80)
