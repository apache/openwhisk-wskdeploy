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
