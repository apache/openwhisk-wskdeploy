package main

import "encoding/json"
import "fmt"
import "os"

func main() {
	//program receives one argument: the JSON object as a string
	arg := os.Args[1]

	// unmarshal the string to a JSON object
	var obj map[string]interface{}
	json.Unmarshal([]byte(arg), &obj)

	// can optionally log to stdout (or stderr)
	fmt.Println("hello Go action")

	name, ok := obj["name"].(string)
	if !ok {
		name = "Stranger"
	}

	// last line of stdout is the result JSON object as a string
	msg := map[string]string{"msg": ("Hello, " + name + "!")}
	res, _ := json.Marshal(msg)
	fmt.Println(string(res))
}
