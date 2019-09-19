package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type ArrayFlags []string

func (a *ArrayFlags) String() string {
	return "ArrayFlags"
}

func (a *ArrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func (a *ArrayFlags) Map() map[string]string {
	result := make(map[string]string)
	for _, v := range *a {
		data := strings.Split(v, "=")
		if len(data) != 2 {
			continue
		}
		result[data[0]] = data[1]
	}
	return result
}

func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func GroupsFromNamespace(ns string) []string {
	parts := strings.Split(ns, "/")
	groups := []string{}
	for index := 0; index < len(parts); index++ {
		part := strings.Join(parts[0:index], "/")
		if len(part) == 0 {
			continue
		}
		groups = append(groups, part)
	}
	return groups
}
