//+build ignore

// Copyright 2021 Rubrik, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package main

import (
	"bytes"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Template used to output the generated Golang source file.
var tmpl = template.Must(template.New("").Parse(`// Code generated by queries_gen.go DO NOT EDIT

// MIT License
//
// Copyright (c) 2021 Rubrik
//	
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//	
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//	
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
	
package graphql

{{ range $name, $value := . }}
// {{ $name }} GraphQL query 
var {{ $name }}Query = {{ $value }}
{{ end }}
`))

// variableName creates a Golang variable name from the query file name by
// trimming the file suffix removing underscore characters and capitalizing
// the character after. The first character is not capitalized to prevent the
// name from being exported.
func variableName(fileName string) string {
	var sb strings.Builder

	name := strings.TrimSuffix(strings.ToLower(fileName), ".graphql")
	for i, part := range strings.Split(name, "_") {
		if i == 0 {
			sb.WriteString(part)
		} else {
			sb.WriteString(strings.Title(part))
		}
	}

	return sb.String()
}

func main() {
	queries := make(map[string]string)

	// Read all graphql files in the current directory and subdirectories.
	err := filepath.Walk("queries", func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".graphql") {
			return err
		}

		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		name := variableName(info.Name())
		query := strings.Replace(strings.TrimSpace(string(buf)), "RubrikPolarisSDKRequest", "SdkGolang_"+name, 1)
		queries[name] = "`" + query + "`"
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, queries)
	if err != nil {
		log.Fatal(err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("queries.go", src, 0666); err != nil {
		log.Fatal(err)
	}
}
