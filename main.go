package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/nirandas/dbmigrate/dbmigrate"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "init" {
		createConfig()
		return
	}
	if loadConfig() {
		dbmigrate.Process(os.Args[1:])
	}
}

func loadConfig() bool {
	var cp = os.Getenv("dbmigrateCONFIGPATH")
	if cp == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		cp = path.Join(dir, "dbmigrate.json")
	}
	fs, err := os.Open(cp)
	if err != nil {
		fmt.Println("Error opening config file " + cp + err.Error())
		fmt.Println("Use dbmigrate init to create an empty config file in this directory")
		return false
	}
	defer fs.Close()

	e := json.NewDecoder(fs)
	err = e.Decode(&dbmigrate.Config)
	if err != nil {
		panic("Error parsing config file " + err.Error())
	}
	dbmigrate.Config.DSN = os.ExpandEnv(dbmigrate.Config.DSN)
	return true
}

func createConfig() {
	var cp = os.Getenv("dbmigrateCONFIGPATH")
	if cp == "" {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		cp = path.Join(dir, "dbmigrate.json")
	}

	if _, err := os.Stat(cp); os.IsNotExist(err) {
		fs, err := os.Create(cp)
		if err != nil {
			fmt.Println("Failed to initialize ", err)
			return
		}
		defer fs.Close()
		fs.Write([]byte(`{
"type":"postgres",
"dsn":"",
"path":"."
}`))
		fmt.Println("Configuration file created at ", cp)
	} else {
		fmt.Println("dbmigrate.json configuration file already exists")
	}
}
