package main

import (
	"bytes"
	"fmt"
	"github.com/klauspost/compress/zstd"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

// go sdk 源码
func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

// go sdk 源码
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// readAndReplaceBody 读取请求的Body内容，然后替换它，以便Body可以再次被读取。
func readAndReplaceBody(req *http.Request) ([]byte, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes, nil
}

func parseQueryString(query string) (map[string][]string, error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err // 返回错误
	}

	return values, nil
}

func startWith(s string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(s, p) {
			return true
		}

	}
	return false
}

func endWith(s string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasSuffix(s, p) {
			return true
		}

	}
	return false
}

func contains(str string, substr ...string) bool {
	for _, s := range substr {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}

func zstdCompare(bodyBytes []byte) ([]byte, error) {
	// 创建一个bytes.Buffer来保存解压后的数据
	decompressed := &bytes.Buffer{}

	// 创建一个Zstd解压器
	zr, err := zstd.NewReader(bytes.NewReader(bodyBytes))
	if err != nil {
		fmt.Println("Error creating Zstd reader:", err)
		return nil, err
	}
	defer zr.Close()

	// 将压缩数据解压到decompressed中
	_, err = decompressed.ReadFrom(zr)
	if err != nil {
		fmt.Println("Error decompressing data:", err)
		return nil, err
	}

	return decompressed.Bytes(), nil
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
