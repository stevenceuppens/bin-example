package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stevenceuppens/bin-example/openapi/gen/bin-api/client"
	"github.com/stevenceuppens/bin-example/openapi/gen/bin-api/client/group"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 1 {
		uploadFile(argsWithoutProg)
		return
	}
	downloadFile()
}

func uploadFile(argsWithoutProg []string) {
	file, err := os.Open(argsWithoutProg[0])

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	// read into buffer
	buffer := bufio.NewReader(file)
	readCloser := ioutil.NopCloser(buffer)

	// upload
	params := group.NewGroupAddPhotoParams().WithPhoto(readCloser)
	api := client.NewHTTPClientWithConfig(nil, client.DefaultTransportConfig().WithHost("127.0.0.1:3000"))
	res, err := api.Group.GroupAddPhoto(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}

func downloadFile() {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	params := group.NewGroupGetPhotoParams()
	api := client.NewHTTPClientWithConfig(nil, client.DefaultTransportConfig().WithHost("127.0.0.1:3000"))
	_, err := api.Group.GroupGetPhoto(params, writer)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(buffer)
}
