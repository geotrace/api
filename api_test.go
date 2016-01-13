package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Param struct {
	Name        string
	Type        string
	Description string
}

func Description(descriptions ...string) {
	if len(descriptions) > 0 {
		for _, description := range descriptions {
			fmt.Printf("%s\n\n", description)
		}
	}
}

func Resource(name, url string, params []Param, descriptions ...string) {
	if len(params) > 0 {
		names := make([]string, len(params))
		for i, param := range params {
			names[i] = fmt.Sprintf("{%s}", param.Name)
		}
		Description(descriptions...)
		fmt.Printf("## %s [%s/%s]\n\n+ Parameters\n", name, strings.Join(names, "/"), url)
		for _, param := range params {
			fmt.Printf("\t+ %s (%s) - %s\n", param.Name, param.Type, param.Description)
		}
		fmt.Print("\n")
	} else {
		fmt.Printf("## %s [%s]\n\n", name, url)
		Description(descriptions...)
	}
}

func Action(name string, req *http.Request, descriptions ...string) (*http.Response, error) {
	fmt.Printf("### %s [%s]\n\n", name, req.Method)
	Description(descriptions...)

	if req.Method != "GET" {
		fmt.Printf("+ Request (%s)\n\n", req.Header.Get("Content-Type"))
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
		var buf bytes.Buffer
		if err := json.Indent(&buf, data, "\t\t", "\t"); err != nil {
			return nil, err
		}
		fmt.Print("\t\t")
		if _, err := buf.WriteTo(os.Stdout); err != nil {
			return nil, err
		}
		fmt.Print("\n\n")
		req.Body = ioutil.NopCloser(bytes.NewReader(data))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("+ Response %d", resp.StatusCode)
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		fmt.Printf(" (%s)", resp.Header.Get("Content-Type"))
	}
	fmt.Print("\n\n")

	var headers bytes.Buffer
	for name, value := range resp.Header {
		switch name {
		case "Content-Length", "Content-Type", "Date":
		default:
			fmt.Fprintf(&headers, "\t\t%s: %s\n", name, strings.Join(value, ", "))
		}
	}
	if headers.Len() > 0 {
		fmt.Print("\t+ Headers\n\n\t\t")
		headers.WriteTo(os.Stdout)
		fmt.Println()
	}

	if resp.ContentLength > 0 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		resp.Body.Close()
		if len(data) > 0 {
			var buf bytes.Buffer
			if err := json.Indent(&buf, data, "\t\t\t", "\t"); err != nil {
				return resp, err
			}
			fmt.Print("\t+ Body\n\n\t\t\t")
			buf.WriteTo(os.Stdout)
			fmt.Print("\n\n")
		}
		resp.Body = ioutil.NopCloser(bytes.NewReader(data))
	}
	return resp, nil
}
