package hamgo

import "strings"

type Filter interface {
	AddAnnoURL(url string) Filter
	DoFilter(ctx Context) bool
	IsAnnoURL(url string) bool
}

type filter struct {
	Filter
	handler func(ctx Context) bool
	annoURL []string
}

func (f *filter) AddAnnoURL(url string) Filter {
	f.annoURL = append(f.annoURL, url)
	return f
}

func (f *filter) DoFilter(ctx Context) bool {
	return f.handler(ctx)
}

func (f *filter) IsAnnoURL(path string) bool {
	for _, url := range f.annoURL {
		if strings.HasPrefix(path, url) {
			//hase anno , pass it ,return true
			return true
		}
	}
	return false
}

func newFilter(handler func(ctx Context) bool) Filter {
	return &filter{handler: handler}
}
