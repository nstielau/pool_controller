from bottle import route, run, template

@route('/')
def index(name):
    print request.body
    return template('<b>Hello {{name}}</b>!', name="there")

run(host='localhost', port=80)
