def main(params):
    if 'name' not in params or 'color' not in params:
        return { 'error': 'Please make sure name and color are passed in params.' }
    name = params['name']
    color = params['color']
    message = 'A ' + color + ' cat named ' + name + ' was added.';
    print(message)
    return { 'change': message }
