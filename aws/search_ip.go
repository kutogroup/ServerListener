package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	ip := ""
	n := ""
	file := ""
	flag.StringVar(&ip, "p", "", "the ip you want search")
	flag.StringVar(&n, "n", "", "the name you want search")
	flag.StringVar(&file, "f", "", "the json file")
	flag.Parse()

	if (len(ip) == 0 && len(n) == 0) || len(file) == 0 {
		flag.PrintDefaults()
		return
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("cant open file, err=", err)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("file cant be opened")
		return
	}
	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println("json parsed failed, err=", err)
		return
	}

	instances := m["staticIps"].([]interface{})
	for _, instance := range instances {
		body := instance.(map[string]interface{})
		name := body["name"]
		ipBody := body["ipAddress"]

		if len(ip) > 0 {
			if ip == ipBody {
				fmt.Println("OK")
				fmt.Println(name)
				fmt.Println(ipBody)
				return
			}
		} else {
			if name == n {
				fmt.Println("OK")
				fmt.Println(name)
				fmt.Println(ipBody)
				return
			}
		}
	}
}
