package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	config     Config
	folderData DirData
)

type Config struct {
	Port       int
	GameFolder string
	Force      bool
}

type DirData struct {
	ChecksumHeader string
	ChecksumsBody  []byte
}

func loadConfig(path string) {
	var err error
	filedata, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(filedata, &config)
	if err != nil {
		log.Fatal(err)
	}
	config.GameFolder, err = filepath.Abs(config.GameFolder)
	if err != nil {
		log.Fatal(err)
	}
}

func loadFolderData() {
	var err error
	dirHasher := sha256.New()
	err = filepath.WalkDir(config.GameFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if d == nil {
				err = fmt.Errorf("invalid root directory")
			}
			log.Fatal(err)
		}
		if d.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			log.Fatal(err)
		}
		checksum := hex.EncodeToString(hasher.Sum(nil))
		path = strings.ReplaceAll(strings.TrimPrefix(path, config.GameFolder), "\\", "/")
		line := []byte(fmt.Sprintf("%s\t%s\n", checksum, path))
		folderData.ChecksumsBody = append(folderData.ChecksumsBody, line...)
		dirHasher.Write(line)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	folderData.ChecksumHeader = fmt.Sprintf("\"%s\"", hex.EncodeToString(dirHasher.Sum(nil)))
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	etag := r.Header.Get("If-None-Match")
	if !config.Force && etag == folderData.ChecksumHeader {
		w.WriteHeader(304)
		return
	}
	w.Header().Add("ETag", folderData.ChecksumHeader)
	w.WriteHeader(200)
	w.Write(folderData.ChecksumsBody)
}

func main() {
	configPath := flag.String("config", "./patch_config.json", "config file")
	flag.Parse()

	loadConfig(*configPath)
	loadFolderData()
	http.HandleFunc("/check", checkHandler)
	http.Handle("/", http.FileServer(http.Dir(config.GameFolder)))
	log.Printf("Starting server on port %d\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
