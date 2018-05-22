// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

func main(args: [String:Any]) -> [String:Any] {
    var msg = ["greeting": "Hello stranger!"]
    if let name = args["name"] as? String {
        if !name.isEmpty {
            msg["greeting"] = "Hello \(name)!"
        }
    }
    print (msg)
    return msg
}


