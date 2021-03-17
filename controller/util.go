package controller

import (
	"gopkg.in/macaron.v1"
	"io"
	"os"
)

func FileUpload(ctx *macaron.Context) (string, error) {
	//this function returns the filename(to save in local disk) of the saved file or an error if it occurs
	//ParseMultipartForm parses a request body as multipart/form-data
	err := ctx.Req.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", err
	}
	//retrieve the file from form data
	file, handler, err := ctx.Req.FormFile("file")
	if err != nil {
		return "", err
	}
	//close the file when we finish
	defer file.Close()
	//replace file with the key your sent your image with
	//this is path which  we want to store the file
	f, err := os.OpenFile("/home/sagor/Desktop/images"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "",err
	}
	defer f.Close()

	io.Copy(f, file)
	//here we save our file to our path
	return handler.Filename, nil
}