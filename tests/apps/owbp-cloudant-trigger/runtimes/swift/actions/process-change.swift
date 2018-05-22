// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

func main(args: [String:Any]) -> [String:Any] {
    if let color = args["color"] as? String,
        let name = args["name"] as? String
    {
      let message = "A \(color) cat named \(name) was added."
      print(message)
      return [ "change": message ]
    } else {
      return [ "error": "Please make sure to pass color and name into params" ]
    }
}
