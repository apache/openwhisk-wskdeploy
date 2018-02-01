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

// qualifiedname.go
package utils

import (
	"errors"
	"strings"
)

type QualifiedName struct {
	Namespace   string
	PackageName string
	EntityName  string
}

// from go whisk cli
/*
Parse a (possibly fully qualified) resource name into namespace and name components. If the given qualified name isNone,
then this is a default qualified name and it is resolved from properties. If the namespace is missing from the qualified
name, the namespace is also resolved from the property file.

Return a qualifiedName struct

Examples:
      foo => qName {namespace: "_", entityName: foo}
      pkg/foo => qName {namespace: "_", entityName: pkg/foo}
      /ns/foo => qName {namespace: ns, entityName: foo}
      /ns/pkg/foo => qName {namespace: ns, entityName: pkg/foo}
*/
func ParseQualifiedName(name string, defaultNamespace string) (QualifiedName, error) {
	var qualifiedName QualifiedName

	// If name has a preceding delimiter (/), it contains a namespace. Otherwise the name does not specify a namespace,
	// so default the namespace to the namespace value set in the properties file; if that is not set, use "_"
	if strings.HasPrefix(name, "/") {
		parts := strings.Split(name, "/")
		qualifiedName.Namespace = parts[1]

		if len(parts) < 2 || len(parts) > 4 {
			err := errors.New("A valid qualified name was not detected")
			return qualifiedName, err
		}

		for i := 1; i < len(parts); i++ {
			if len(parts[i]) == 0 || parts[i] == "." {
				err := errors.New("A valid qualified name was not detected")
				return qualifiedName, err
			}
		}

		qualifiedName.EntityName = strings.Join(parts[2:], "/")
	} else {
		if len(name) == 0 || name == "." {
			err := errors.New("A valid qualified name was not detected")
			return qualifiedName, err
		}

		qualifiedName.EntityName = name
		if defaultNamespace == "" {
			defaultNamespace = "_"
		}
		qualifiedName.Namespace = defaultNamespace
	}

	return qualifiedName, nil
}
