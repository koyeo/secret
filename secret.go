package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"plugin"
	"strings"
)

const secretSoEnv = "SECRET_SO"

func main() {
	_key := ""
	checkKeyLength := func(key string) (err error) {
		if len(key)%8 != 0 {
			err = fmt.Errorf("key length should be a multiple of 8")
			return
		}
		return
	}
	getKey := func(c *cli.Context) (key string, err error) {
		key = _key
		if c.String("key") != "" {
			key = c.String("key")
		}
		key = strings.TrimSpace(key)
		if key == "" {
			err = fmt.Errorf("key is empty")
			return
		}
		err = checkKeyLength(key)
		if err != nil {
			return
		}

		return
	}
	app := &cli.App{
		Name:  "secret",
		Usage: "use secret.so encrypt or decrypt data",
		Commands: []*cli.Command{
			{
				Name:    "encrypt",
				Aliases: []string{"e"},
				Usage:   "encrypt data,    e.g: secret encrypt, -key 32lengthStr plainString",
				Action: func(c *cli.Context) (err error) {
					key, err := getKey(c)
					if err != nil {
						return
					}
					data := strings.TrimSpace(c.Args().First())
					if data == "" {
						err = fmt.Errorf("data is empty")
						return
					}
					res, err := encryptText(data, key)
					if err != nil {
						return
					}
					fmt.Println(res)
					return
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "key",
						Usage: "as encrypt key",
					},
				},
			},
			{
				Name:    "decrypt",
				Aliases: []string{"d"},
				Usage:   "decrypt data,    e.g: secret decrypt -key 32lengthStr cipherString",
				Action: func(c *cli.Context) (err error) {
					key, err := getKey(c)
					if err != nil {
						return
					}
					data := strings.TrimSpace(c.Args().First())
					if data == "" {
						err = fmt.Errorf("data is empty")
						return
					}
					res, err := decryptText(data, key)
					if err != nil {
						return
					}
					fmt.Println(res)
					return
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "key",
						Usage: "as encrypt key",
					},
				},
			},
			{
				Name:  "hash",
				Usage: "hash data,    e.g: secret hash -key 32lengthStr cipherString",
				Action: func(c *cli.Context) (err error) {
					key, err := getKey(c)
					if err != nil {
						return
					}
					data := strings.TrimSpace(c.Args().First())
					if data == "" {
						err = fmt.Errorf("data is empty")
						return
					}
					res := hashText(data, key)
					fmt.Println(res)
					return
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "key",
						Usage: "as hash key",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func getSoPath() string {
	path := os.Getenv(secretSoEnv)
	if path == "" {
		path = "./secret.so"
	}
	return path
}

func getFunc(name string) (symbol plugin.Symbol, err error) {

	path := getSoPath()

	if !pathExist(path) {
		err = fmt.Errorf("so file not exist, path: \"%s\"", path)
		return
	}

	so, err := plugin.Open(path)
	if err != nil {
		fmt.Printf("Env var: \"%s\"\n", secretSoEnv)
		fmt.Printf("Use so file path: \"%s\"\n", path)
		fmt.Printf("Open error : \"%s\"\n", err)
	}
	symbol, err = so.Lookup(name)
	if err != nil {
		fmt.Printf("Lookup Func %s error: %s\n", name, err)
		return
	}
	return
}

func pathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return true
}

func encryptText(text, salt string) (cipher string, err error) {
	f, err := getFunc("EncryptText")
	if err != nil {
		return
	}
	return f.(func(string, string) (string, error))(text, salt)
}

func decryptText(cipher, salt string) (text string, err error) {

	f, err := getFunc("DecryptText")
	if err != nil {
		return
	}
	return f.(func(string, string) (string, error))(cipher, salt)
}

func hashText(data, salt string) (hash string) {

	f, err := getFunc("Hash")
	if err != nil {
		return
	}
	return f.(func(string, string) string)(data, salt)
}
