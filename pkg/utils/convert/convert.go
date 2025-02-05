package convert

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/a-h/templ"
)

func FormatLink(siteUrl string, relativePath string) string {
	return fmt.Sprintf("%s%s/", siteUrl, strings.Trim(relativePath, "/"))
}

func DefaultTo(value interface{}, defaults interface{}) interface{} {
	if value != nil {
		if reflect.TypeOf(value).Kind() != reflect.String || len(value.(string)) > 0 {
			return value
		}
	}

	return defaults
}

func AddTrailingSlash(str string) string {
	if str[len(str)-1:] != "/" {
		return str + "/"
	}

	return str
}

func ToSafeUrl(url string, trailingSlash bool) templ.SafeURL {
	if trailingSlash {
		return templ.URL(AddTrailingSlash(url))
	}

	return templ.URL(url)
}

func ToPath(dir string) string {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return path.Join(wd, dir)
}

func ToUnsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
