package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/stevenceuppens/bin-example/openapi/gen/bin-store/server"
	"github.com/stevenceuppens/bin-example/openapi/gen/bin-store/server/operations"
	"github.com/stevenceuppens/bin-example/openapi/gen/bin-store/server/operations/group"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

var servicePort = 3001

var tmp []byte

func main() {
	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(server.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// create new service API
	openapi := operations.NewPhotoStoreAPI(swaggerSpec)
	server := server.NewServer(openapi)
	defer server.Shutdown()

	// set the port this service will be run on
	server.Port = servicePort

	// Implement routes
	openapi.GroupGroupAddPhotoHandler = group.GroupAddPhotoHandlerFunc(func(params group.GroupAddPhotoParams) middleware.Responder {
		buf, err := ioutil.ReadAll(params.Photo)
		if err != nil {
			fmt.Println("got err: ", err)
			return group.NewGroupAddPhotoInternalServerError()
		}

		fmt.Println("got data: ", buf)

		tmp = buf

		return group.NewGroupAddPhotoOK()
	})

	openapi.GroupGroupGetPhotoHandler = group.GroupGetPhotoHandlerFunc(func(params group.GroupGetPhotoParams) middleware.Responder {

		reader := bytes.NewReader(tmp)
		return group.NewGroupGetPhotoOK().WithPayload(ioutil.NopCloser(reader))
	})

	fmt.Println("Start Store @ ", servicePort)

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
