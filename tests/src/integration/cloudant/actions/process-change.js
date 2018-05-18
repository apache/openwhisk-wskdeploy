/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor
 * license agreements; and to You under the Apache License, Version 2.0.
 */

/**
 * Print the document to console which has changes
 *
 * @param document which has changes
 */

function main(params) {
      var message = "The changed document ID is:"+params._id;
      console.log(message);
      return {change: message}
}
