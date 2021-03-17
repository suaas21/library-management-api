package controller

import (
	"encoding/csv"
	"github.com/suaas21/library-management-api/database"
	"gopkg.in/macaron.v1"
	"io"
	"net/http"
	"os"
)

func ExportDataToCSV(ctx *macaron.Context) {
	maps, err := database.ExportDataFromDB()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}
	err = CSVExport(maps)
	if err != nil {
		ctx.JSON(http.StatusNotImplemented, err.Error())
		return
	}
	return
}


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

func CSVExport(mdata []map[string]string) error {
	data := convertMapToSlice(mdata)
	// Changed to csvExport, as it doesn't make much sense to export things from
	file, err := os.Create("loan-book-history.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(data); err != nil {
		return err // let's return errors if necessary, rather than having a one-size-fits-all error handler
	}
	return nil
}

func convertMapToSlice(maps []map[string]string) []string {
	totalMapSize := 0
	for _, mp := range maps {
		mpSize := len(mp)
		totalMapSize = mpSize + totalMapSize
	}

	v := make([]string, totalMapSize*len(maps))
	idx := 0
	for _, mp := range maps {
		for  _, m := range mp {
			v[idx] = m
			idx++
		}
	}
	return v
}