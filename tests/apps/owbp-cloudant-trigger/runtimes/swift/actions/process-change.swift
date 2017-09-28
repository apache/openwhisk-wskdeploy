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
