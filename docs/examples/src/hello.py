# Licensed to the Apache Software Foundation (ASF) under one or more contributor
# license agreements; and to You under the Apache License, Version 2.0.

"""OpenWhisk "Hello world" in Python.

// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.
"""


def main(dict):
    """Hello world."""
    if 'name' in dict:
        name = dict['name']
    else:
        name = "stranger"
    if 'place' in dict:
        place = dict['place']
    else:
        place = "unknown"
    msg = "Hello, " + name + " from " + place
    return {"greeting": msg}
