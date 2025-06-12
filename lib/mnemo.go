package lib

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gtank/cryptopasta"
	"github.com/kkrishguptaa/mnemo/util"
)

type Store struct {
	Name string    `json:"name"`
	Data []Snippet `json:"data"`
}

type Snippet struct {
	Id        string `json:"id"`
	Value     string `json:"value"`
	Encrypted bool   `json:"encrypted"`
}

var home = util.ErrorHandler(os.UserHomeDir())
var stores = path.Join(home, ".mnemo", "stores")

func Encrypt(value string, password string) string {
	password = strings.TrimSpace(password)
	hash := sha256.Sum256([]byte(password))
	cipher := util.ErrorHandler(cryptopasta.Encrypt([]byte(value), &hash))

	return base64.StdEncoding.EncodeToString(cipher)
}

func Decrypt(value string, password string) string {
	password = strings.TrimSpace(password)
	hash := sha256.Sum256([]byte(password))

	decoded, err := base64.StdEncoding.DecodeString(value)
	util.ErrorPrinter(err) // or panic(err), or your own handler

	decrypted := util.ErrorHandler(cryptopasta.Decrypt(decoded, &hash))
	return string(decrypted)
}

func FetchStore(name string) Store {
	file := path.Join(stores, name+".json")

	if _, err := os.Stat(stores); os.IsNotExist(err) {
		util.ErrorPrinter(os.MkdirAll(stores, 0755))
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		if name != "default" {
			util.ErrorPrinter(fmt.Errorf("store %s does not exist", name))
			return Store{}
		}

		CreateStore(name)
	}

	bytes := util.ErrorHandler(os.ReadFile(file))

	store := Store{Name: name}
	util.ErrorPrinter(json.Unmarshal(bytes, &store))

	return store
}

func CreateStore(name string) Store {
	if name == "" {
		util.ErrorPrinter(fmt.Errorf("store name cannot be empty"))
		return Store{}
	}

	if _, err := os.Stat(stores); os.IsNotExist(err) {
		util.ErrorPrinter(os.MkdirAll(stores, 0755))
	}

	file := path.Join(stores, name+".json")

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		util.ErrorPrinter(fmt.Errorf("store %s already exists", name))
		return Store{}
	}

	store := Store{Name: name, Data: []Snippet{}}
	bytes, err := json.Marshal(store.Data)
	util.ErrorPrinter(err)

	util.ErrorPrinter(os.WriteFile(file, bytes, 0644))

	return store
}

func ListStores() []string {
	if _, err := os.Stat(stores); os.IsNotExist(err) {
		util.ErrorPrinter(os.MkdirAll(stores, 0755))
	}

	files, err := os.ReadDir(stores)
	util.ErrorPrinter(err)

	var storeNames []string
	for _, file := range files {
		if !file.IsDir() && path.Ext(file.Name()) == ".json" {
			storeNames = append(storeNames, file.Name()[:len(strings.Split(file.Name(), ".")[0])])
		}
	}

	return storeNames
}

func WriteStore(name string, snippets []Snippet) {
	file := path.Join(stores, name+".json")

	store := Store{Name: name, Data: snippets}

	bytes := util.ErrorHandler(json.Marshal(store))

	util.ErrorPrinter(os.WriteFile(file, bytes, 0644))
}

func CreateSnippet(store Store, snipper Snippet) Snippet {
	if snipper.Id == "" {
		util.ErrorPrinter(fmt.Errorf("snippet id cannot be empty"))
		return Snippet{}
	}

	for _, snippet := range store.Data {
		if snippet.Id == snipper.Id {
			util.ErrorPrinter(fmt.Errorf("snippet with id %s already exists", snipper.Id))
			return Snippet{}
		}
	}

	store.Data = append(store.Data, snipper)
	WriteStore(store.Name, store.Data)

	return snipper
}

func DeleteStore(name string) {
	file := path.Join(stores, name+".json")

	if _, err := os.Stat(file); os.IsNotExist(err) {
		util.ErrorPrinter(fmt.Errorf("store %s does not exist", name))
		return
	}

	util.ErrorPrinter(os.Remove(file))
}
