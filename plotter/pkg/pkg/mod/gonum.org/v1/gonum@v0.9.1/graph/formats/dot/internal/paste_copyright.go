// This file is dual licensed under CC0 and The Gonum License.
//
// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright ©2017 Robin Eklind.
// This file is made available under a Creative Commons CC0 1.0
// Universal Public Domain Dedication.`)

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var location = []byte(`// Code generated by gocc; DO NOT EDIT.`)
var copyright = []byte(`// Code generated by gocc; DO NOT EDIT.

// This file is dual licensed under CC0 and The Gonum License.
//
// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright ©2017 Robin Eklind.
// This file is made available under a Creative Commons CC0 1.0
// Universal Public Domain Dedication.`)

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Dir(path) == "." || filepath.Ext(path) != ".go" {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		content = bytes.Replace(content, location, copyright, 1)
		return ioutil.WriteFile(path, content, info.Mode())
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
	}
}
