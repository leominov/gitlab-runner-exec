package main

import (
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
