/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var uploadURL = "http://localhost:1323/api/upload"

// https://github.com/mimoo/eureka/blob/master/folders.go
func tarDir(src string) (string, error) {

	// prepare dest tarfile and tarWriter
	tempDir, err := os.MkdirTemp(".", "tmp-")
	if err != nil {
		return "", err
	}
	destFilePath := filepath.Join(tempDir, filepath.Base(src)+".tar")
	log.Println(destFilePath)
	tarfile, err := os.OpenFile(destFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	// tarfile, err := os.Create(destFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tw := tar.NewWriter(fileWriter)
	defer tw.Close()

	// start tar
	// is file a folder?
	fi, err := os.Stat(src)
	if err != nil {
		return "", err
	}
	mode := fi.Mode()
	if mode.IsRegular() {
		// get header
		header, err := tar.FileInfoHeader(fi, src)
		if err != nil {
			return "", err
		}
		// write header
		if err := tw.WriteHeader(header); err != nil {
			return "", err
		}
		// get content
		data, err := os.Open(src)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(tw, data); err != nil {
			return "", err
		}
	} else if mode.IsDir() { // folder

		// walk through every file in the folder
		filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
			// generate tar header
			header, err := tar.FileInfoHeader(fi, file)
			if err != nil {
				return err
			}

			// must provide real name
			// (see https://golang.org/src/archive/tar/common.go?#L626)
			header.Name = filepath.ToSlash(file)

			// write header
			if err := tw.WriteHeader(header); err != nil {
				return err
			}
			// if not a dir, write file content
			if !fi.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					return err
				}
				if _, err := io.Copy(tw, data); err != nil {
					return err
				}
			}
			return nil
		})
	} else {
		return "", fmt.Errorf("error: file type not supported")
	}

	// produce tar
	if err := tw.Close(); err != nil {
		return "", err
	}

	return destFilePath, nil
}

func upload(cmd *cobra.Command, args []string) {
	startTime := time.Now()
	if uploadDirPath != "" {
		// do compress
		var err error
		uploadFilePath, err = tarDir(uploadDirPath)
		log.Printf("made tar time: %v\n", time.Since(startTime))
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			tmpDir := filepath.Dir(uploadFilePath)
			log.Println("remove dir ", tmpDir)
			if err := os.RemoveAll(tmpDir); err != nil {
				log.Println(err)
				return
			}
			log.Println("remove tar: ", uploadFilePath)
		}()
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

	log.Printf("upload time: %v\n", time.Since(startTime))
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
