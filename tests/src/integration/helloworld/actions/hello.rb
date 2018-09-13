# Licensed to the Apache Software Foundation (ASF) under one or more contributor
# license agreements; and to You under the Apache License, Version 2.0.

def main(args)
  name = args["name"] || "stranger"
  greeting = "Hello #{name}!"
  puts greeting
  { "greeting" => greeting }
end
