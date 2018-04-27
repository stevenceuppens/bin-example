# go swagger - proxy example binary data

## Run Example

Generate openapi and build go applications
```bash
$ make build
```

Open 3 terminals in bin directory

Terminal 1 : Start Store
```bash
$ ./store
```

Terminal 2 : Start Store
```bash
$ ./proxy
```

Terminal 3 : launch CLI and pass an argument to a binary file (itselves)
```bash
$ ./cli cli
```

## Explenation

To pass binary in the proxy, you can forward the binary data
from the params (incomming request) to the params (outgoing)
without the need of reading into a buffer

```go
// Implement routes
openapi.GroupGroupAddPhotoHandler = group.GroupAddPhotoHandlerFunc(func(params group.GroupAddPhotoParams) middleware.Responder {
    // upload
    fmt.Println("Sync Start")

    api := c_client.NewHTTPClientWithConfig(nil, c_client.DefaultTransportConfig().WithHost("127.0.0.1:3001"))

    _, err := api.Group.GroupAddPhoto(c_group.NewGroupAddPhotoParams().WithPhoto(params.Photo))
    if err != nil {
        fmt.Println(err)
        return group.NewGroupAddPhotoInternalServerError()
    }

    fmt.Println("Sync OK")

    return group.NewGroupAddPhotoOK()
})
```