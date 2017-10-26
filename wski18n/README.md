# How to generate the file i18n_resources.go for internationalization

As a contributor to wskdeploy, the file of i18n_resources.go needs to regenerated,
when you add or change any localized message. In order to generate i18n_resources.go,
you need to install go-bindata first:

```
$ go get -u github.com/jteeuwen/go-bindata/...
```

Then, go the HOME directory of wskdeploy and run the following command:

```
$ $GOPATH/bin/go-bindata -pkg wski18n -o wski18n/i18n_resources.go wski18n/resources;
```

Finally, add the default ASF license header to i18n_resources.go. Since each file of
source code starts with the ASF license header, you need to add it to i18n_resources.go
each time it is regenerated. You can find this license header in any other file of source
code, e.g. i18n.go.
