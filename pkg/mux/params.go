package mux

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// PathPart ...
type PathPart struct {
	name        string
	placeholder bool
}

var parameter = make(map[string]string)

// MuxParams ...
type MuxParams struct {
	pattern string
	params  []PathPart
	weight  int
	handler http.Handler
}

// GetParamsFromURL ...
func GetParamsFromURL(path string) []string {
	if path[0] != 47 {
		log.Println("not found slash")
	}

	params := make([]string, 0)
	start := 0
	end := 0
	for i, p := range path {

		if p == '{' {
			start = i + 1
		}

		if p == '}' {
			end = i
		}

		if start != 0 && end != 0 {
			params = append(params, path[start:end])
			start = 0
			end = 0
		}

	}

	return params
}

// CompareURLs ...
func CompareURLs(path, pattern string) map[string]string {
	params := make([]string, 0)
	indexs := make([]int, 0)
	start := 0

	for i, p := range path {
		if p == '/' {
			start = i
		}
		if start != 0 {
			indexs = append(indexs, start)
			start = 0
		}
	}

	indexs = append(indexs, len(path))
	for i, p := range indexs {
		if i != 0 && i%2 == 1 {
			params = append(params, path[indexs[i-1]+1:p])
		}
	}

	values := GetParamsFromURL(pattern)
	for i := range params {
		parameter[values[i]] = params[i]
	}

	return parameter
}

func parsePathParams(pattern string) MuxParams {
	parts := strings.Split(pattern, "/")
	params := MuxParams{
		pattern: pattern,
		params:  make([]PathPart, 0, len(parts)),
		weight:  calculateWeight(pattern),
	}
	for _, part := range parts {
		params.params = append(params.params, ParsePathPart(part))
	}
	return params
}

// ParsePathPart ...
func ParsePathPart(part string) PathPart {
	if part == "" {
		pathPart := PathPart{
			name:        part,
			placeholder: false,
		}
		return pathPart
	}

	if part[0] == '{' {
		if part[len(part)-1] != '}' {
			panic(fmt.Errorf("invalid path part: %s", part))
		}

		pathPart := PathPart{
			name:        part[1 : len(part)-1],
			placeholder: true,
		}

		return pathPart
	}

	pathPart := PathPart{
		name:        part,
		placeholder: false,
	}

	return pathPart
}

func (p *MuxParams) placeholders() []string {
	result := make([]string, 0)

	for _, param := range p.params {
		if !param.placeholder {
			continue
		}
		result = append(result, param.name)
	}
	return result
}

// Match ...
func (p *MuxParams) Match(path string) (map[string]string, bool) {
	parts := strings.Split(path, "/")
	if len(parts) != len(p.params) {
		return nil, false
	}

	params := make(map[string]string)

	for index, param := range p.params {
		if !param.placeholder {
			if param.name != parts[index] {
				return nil, false
			}
			continue
		}

		if parts[index] == "" {
			return nil, false
		}

		params[param.name] = parts[index]
	}

	return params, true
}

// calculateWeight ...
func calculateWeight(pattern string) int {
	if pattern == "/" {
		return 0
	}

	count := (strings.Count(pattern, "/") - 1) * 2
	if !strings.HasSuffix(pattern, "/") {
		return count + 1
	}
	return count
}
