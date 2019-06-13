/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * Return success saying a book was updated into the book store.
 */
function main(params) {
    return new Promise(function(resolve, reject) {
        console.log(params.name);
        console.log(params.isbn);

        if (!params.name) {
            console.error('name parameter not set.');
            reject({
                'error': 'name parameter not set.'
            });
            return;
        } else if (!params.isbn) {
            console.error('isbn parameter not set.');
            reject({
                'error': 'isbn parameter not set.'
            });
            return;
        } else {
            var message = 'A book ' + params.name + ' was updated to a new ISBN ' + params.isbn;
            console.log(message);
            resolve({
                result: message
            });
            return;
        }
    });
}
