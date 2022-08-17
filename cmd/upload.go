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
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadEndpoint = "/v1/upload"

// https://github.com/mimoo/eureka/blob/master/folders.go
func tarDir(src string) (string, error) {

	// prepare dest tarfile and tarWriter
	tempDir, err := os.MkdirTemp(".", "tmp-")
	if err != nil {
		return "", err
	}
	destFilePath := filepath.Join(tempDir, filepath.Base(src)+".tar")
	log.Debug().Msg(destFilePath)
	tarfile, err := os.OpenFile(destFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tw := tar.NewWriter(fileWriter)
	defer tw.Close()

	// move to the base dir
	srcFi, err := os.Stat(src)
	if err != nil {
		return "", err
	}
	if srcFi.Mode().IsDir() {
		if err := os.Chdir(src); err != nil {
			return "", err
		}
		if err := os.Chdir("../"); err != nil {
			return "", err
		}
	}

	// start tar
	// is file a folder?
	baseDir := filepath.Base(src)
	fi, err := os.Stat(baseDir)
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
		data, err := os.Open(baseDir)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(tw, data); err != nil {
			return "", err
		}
	} else if mode.IsDir() { // folder

		// walk through every file in the folder
		filepath.Walk(baseDir, func(file string, fi os.FileInfo, err error) error {
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
	cwd, err := os.Getwd() // remember current directory
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	if uploadDirPath != "" {
		// do compress
		var err error
		uploadFilePath, err = tarDir(uploadDirPath)
		fmt.Printf("tar completed. time elapsed: %.3f seconds\n", time.Since(startTime).Seconds())
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		// back to current dir after tarDir
		log.Debug().Str("back to current directory", cwd).Send()
		if err := os.Chdir(cwd); err != nil {
			log.Fatal().Err(err).Send()
		}

		defer func() {
			tmpDir := filepath.Dir(uploadFilePath)
			log.Debug().Str("remove dir", tmpDir).Send()
			if err := os.RemoveAll(tmpDir); err != nil {
				log.Fatal().Err(err).Send()
			}
			log.Debug().Str("remove tar", uploadFilePath)
		}()
	}

	log.Debug().Str("uploadFilePath", uploadFilePath)
	file, err := os.Open(uploadFilePath)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer file.Close()

	// file size
	fi, err := file.Stat()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	fmt.Printf("tar size: %d MB\n", fi.Size()/1024/1024)

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, _ := writer.CreateFormFile("file", file.Name())
	io.Copy(part, file)
	writer.Close()

	uploadURL := viper.GetString(viperServerURL) + uploadEndpoint
	// TODO:verify url format

	req, err := http.NewRequest(http.MethodPost, uploadURL, requestBody)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	// req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viper.GetString(viperToken)))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer resp.Body.Close()

	fmt.Printf("--upload elapsed time: %.3f seconds\n", time.Since(startTime).Seconds())
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	fmt.Println("--response body:", string(respBody))
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload file or directory",
	Long: `upload file or directory
example:

$ annocli upload --dir=upload/coco1

$ annocli upload --file=upload/coco1/COCO_train2014_000000000009.jpg
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
	uploadCmd.Flags().StringVar(&uploadDirPath, "dir", "", "directory path to upload")
	// uploadCmd.MarkFlagRequired("file")
}
