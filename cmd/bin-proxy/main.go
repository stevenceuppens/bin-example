package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/stevenceuppens/bin-example/openapi/gen/bin-api/server"
	"github.com/stevenceuppens/bin-example/openapi/gen/bin-api/server/operations"
	"github.com/stevenceuppens/bin-example/openapi/gen/bin-api/server/operations/group"

	c_client "github.com/stevenceuppens/bin-example/openapi/gen/bin-store/client"
	c_group "github.com/stevenceuppens/bin-example/openapi/gen/bin-store/client/group"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

var servicePort = 3000

func main() {
	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(server.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// create new service API
	openapi := operations.NewPhotoAPI(swaggerSpec)
	server := server.NewServer(openapi)
	defer server.Shutdown()

	// set the port this service will be run on
	server.Port = servicePort

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

	openapi.GroupGroupGetPhotoHandler = group.GroupGetPhotoHandlerFunc(func(params group.GroupGetPhotoParams) middleware.Responder {
		fmt.Println("Sync Start")

		api := c_client.NewHTTPClientWithConfig(nil, c_client.DefaultTransportConfig().WithHost("127.0.0.1:3001"))

		// create buffer to temporarly store the data
		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)

		// call upstream
		_, err := api.Group.GroupGetPhoto(c_group.NewGroupGetPhotoParams(), writer)
		if err != nil {
			fmt.Println(err)
			return group.NewGroupGetPhotoInternalServerError()
		}

		fmt.Println("Sync OK")

		// convert buffer into reader
		reader := bufio.NewReader(&buffer)
		// create closer
		closer := ioutil.NopCloser(reader)

		return group.NewGroupGetPhotoOK().WithPayload(closer)
	})

	fmt.Println("Start Api @ ", servicePort)

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
