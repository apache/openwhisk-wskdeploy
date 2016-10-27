package main

import (
  "fmt"
  "net/url"
  "net/http"
  "github.com/openwhisk/go-whisk/whisk"
)

func testCreateAction() {
  createAction("helloWskTool", "func main(args: [String:Any]) -> [String:Any] { return [\"msg\":\"wsktool says hello\"] }", "swift:3")
}

func createAction(name string, code string, kind string) {

  baseUrl, err := getURLBase("openwhisk.ng.bluemix.net")
  if err != nil {
    fmt.Println("Got error making baseUrl ", err)
  }

  clientConfig := &whisk.Config{
      AuthToken:  "MyAuthToken",
      Namespace:  "MyNameSpace",
      BaseURL:    baseUrl,
      Version:    "v1",
      Insecure:   false, // true if you want to ignore certificate signing
  }

  // Setup network client
  client, err := whisk.NewClient(http.DefaultClient, clientConfig)
  if err != nil {
    fmt.Println("Got error making whisk client ", err)
  }

  // create action struct
  action := new(whisk.Action)
  action.Exec = new(whisk.Exec)
  action.Exec.Code = code
  action.Exec.Kind = kind
  action.Name = name
  action.Namespace = "castrop@us.ibm.com"
  action.Publish = false
  // action.Parameters =
  // action.Annotations =
  // action.Limits =

  // call ActionService Thru Client
  _, _, err = client.Actions.Insert(action, false, false)
  if err != nil {
    fmt.Println("Got error inserting action ", err)
  }
}

// Utility to convert hostname to URL object
func getURLBase(host string) (*url.URL, error)  {

    urlBase := fmt.Sprintf("%s/api/", host)
    url, err := url.Parse(urlBase)

    if len(url.Scheme) == 0 || len(url.Host) == 0 {
        urlBase = fmt.Sprintf("https://%s/api/", host)
        url, err = url.Parse(urlBase)
    }

    return url, err
}
