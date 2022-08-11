/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var uploadURL = "http://localhost:1323/api/upload"

func tarDir(sourceDir string) (string, error) {
	destFilePath := filepath.Join(sourceDir, filepath.Base(sourceDir)+".tar")

	dir, err := os.Open(sourceDir)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	// get list of files
	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	// create tar file
	tarfile, err := os.Create(destFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tw := tar.NewWriter(fileWriter)
	defer tw.Close()

	for _, fileInfo := range files {
		// if fi÷leInfo.IsDir() {
		if !fileInfo.Mode().IsRegular() {
			continue
		}

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// prepare the tar header
		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tw.WriteHeader(header)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(tw, file)
		if err != nil {
			log.Fatal(err)
		}
	}

	return destFilePath, nil
}

func upload(cmd *cobra.Command, args []string) {

	if uploadDirPath != "" {
		// do compress
		var err error
		uploadFilePath, err = tarDir(uploadDirPath)
		if err != nil {
			log.Fatal(err)
		}
		defer func(tarfile string) {
			err := os.Remove(tarfile)
			log.Println("remove tar: ", tarfile, err)
		}(uploadFilePath)
	}
	log.Println("uploadFilePath:", uploadFilePath)
	file, err := os.Open(uploadFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", file.Name())
	io.Copy(part, file)
	writer.Close()

	r, _ := http.NewRequest("POST", uploadURL, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err:", err.Error())
		return
	}
	log.Println("respbody:", string(respbody))
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload file",
	Long: `upload file
Usage:

wip
`,
	Run: upload,
}

// flags
var (
	uploadFilePath string
	uploadDirPath  string
)

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	uploadCmd.Flags().StringVarP(&uploadFilePath, "file", "f", "", "file path to upload")
	uploadCmd.Flags().StringVarP(&uploadDirPath, "dir", "d", "", "directory path to upload")
	// uploadCmd.MarkFlagRequired("file")
}
