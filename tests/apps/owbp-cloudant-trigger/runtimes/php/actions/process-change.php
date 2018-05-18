// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

function main(array $args) : array
{
  if (array_key_exists("name", $args) && array_key_exists("color", $args)) {
    $name = $args["name"];
    $color = $args["color"];
    $message = "A $color cat named $name was added.";
    print($message);
    return ["change" => $message];
  } else {
    return ["error" => "Please make sure name and color parameters are set."];
  }
}
