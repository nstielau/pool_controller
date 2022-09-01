from bottle import route, run, template, request

@route('/', method=['GET', 'POST'])
def index():
    print(request.body.read())
    print()
    return "hello"

run(host='0.0.0.0', port=80)
