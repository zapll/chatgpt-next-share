package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//go:embed login.html
var loginHtml string

var offline = false

var offlineCode = ""

var accessTokenMap = sync.Map{}

func getLoginHtml(failed bool) []byte {
	data := map[string]interface{}{
		"Failed": failed,
	}
	tmpl, _ := template.New("example").Parse(loginHtml)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// 自定义代理逻辑
func handleProxy(res http.ResponseWriter, req *http.Request) {
	Debug(fmt.Sprintf("Request: %s %s", req.Method, req.URL.Path))
	if req.URL.Path == "/offline/"+offlineCode {
		offline = !offline
		res.Write([]byte(fmt.Sprintf("Service now is %s", func() string {
			if offline {
				return "offline"
			}
			return "online"
		}())))
		return
	}

	if offline {
		res.WriteHeader(http.StatusServiceUnavailable)
		res.Write([]byte("Service is offline"))
		return
	}

	if startWith(req.URL.Path, "/backend-api/accounts") && endWith(req.URL.Path, "/invites") {
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}

	if req.URL.Path == "/auth/login" && req.Method == http.MethodGet {
		res.Write(getLoginHtml(false))
		return
	}

	if req.URL.Path != "/auth/login" && req.Method == http.MethodGet {
		token, err := req.Cookie("session_token")

		if err != nil {
			http.Redirect(res, req, "/auth/login", http.StatusFound)
			return
		}
		_, err = getCookieByToken(token.Value)
		if err != nil {
			http.Redirect(res, req, "/auth/login", http.StatusFound)
			return
		}
	}

	if req.URL.Path == "/auth/login" && req.Method == http.MethodPost {
		body, err := readAndReplaceBody(req)
		if err != nil {
			Error("Failed to get body", err)
			res.Write(getLoginHtml(true))
			return
		}

		params, err := parseQueryString(string(body))
		if err != nil {
			Error("Failed to parse query string", err)
			res.Write(getLoginHtml(true))
			return
		}

		token := params["password"][0]

		cookie, err := getCookieByToken(token)
		if err != nil {
			Error("Failed to get cookie by token", err)
			res.Write(getLoginHtml(true))
			return
		}

		resp, err := resty.New().R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", cookie)).
			Post(openaiUrl("/auth/login/token"))

		if err != nil {
			Error("Failed to get body", err)
			res.Write(getLoginHtml(true))
			return
		}

		for _, c := range resp.Cookies() {
			if c.Name == "session_token" {
				c.Value = token
			}
			http.SetCookie(res, c)
		}

		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	if req.URL.Path == "/backend-api/conversations" && req.Method == http.MethodGet {
		rebuildReq(req)
		handleConversations(res, req)
		return
	}

	if req.Method == http.MethodGet && req.URL.Path == "/api/auth/session" {
		token, _ := req.Cookie("session_token")
		rebuildReq(req)
		client := newRequest(req)

		resp, err := client.Get(openaiUrl("/api/auth/session"))
		if err != nil {
			http.Error(res, "Failed to get boyd", http.StatusBadRequest)
			return
		}

		body := resp.Body()
		accessToken := gjson.GetBytes(body, "accessToken").String()
		accessTokenMap.Store(token.Value, accessToken)
		Debug("set accessToken", token.Value, accessToken)

		value, _ := sjson.Set(string(body), "accessToken", token.Value)
		value, _ = sjson.Set(value, "user.id", "user-xxxxxxxxxxxxxxxxxx")
		value, _ = sjson.Set(value, "user.name", "user@share.com")
		value, _ = sjson.Set(value, "user.email", "user@share.com")
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(value))
		return
	}

	if req.URL.Path == "/backend-api/me" && req.Method == http.MethodGet {
		rebuildReq(req)

		client := newRequest(req)
		resp, err := client.Get(openaiUrl("/backend-api/me"))
		if err != nil {
			http.Error(res, "Failed to get boyd", http.StatusBadRequest)
			return
		}

		body := resp.Body()
		value, _ := sjson.Set(string(body), "id", "user-xxxxxxxxxxxxxxxxxx")
		value, _ = sjson.Set(value, "name", "share")
		value, _ = sjson.Set(value, "email", "user@share.com")
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(value))
		return
	}

	if req.URL.Path == "/backend-api/conversation" && req.Method == http.MethodPost {
		user, err := getUser(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			Error("Failed to get user", err)
			return
		}
		expireAt := user.GetString("expire_at")
		if expireAt != "" {
			expireAt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", expireAt)
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				Error("Failed to parse expire_at", err)
				return
			}
			if expireAt.Before(time.Now()) {
				res.Header().Set("Content-Type", "application/json")
				res.WriteHeader(http.StatusGone)
				res.Write([]byte(`{"detail": "您的账号已过期，请联系管理员"}`))
				return
			}
		}
	}

	proxy, _url, err := openai(res, req)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		Error("Failed to get proxy", err)
		return
	}

	if proxy == nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("Service not found"))
		Error("Service not found")
		return
	}

	proxy.Director = func(request *http.Request) {
		targetQuery := _url.RawQuery
		request.URL.Scheme = _url.Scheme
		request.URL.Host = _url.Host
		request.Host = _url.Host
		request.URL.Path, request.URL.RawPath = joinURLPath(_url, request.URL)

		if targetQuery == "" || request.URL.RawQuery == "" {
			request.URL.RawQuery = targetQuery + request.URL.RawQuery
		} else {
			request.URL.RawQuery = targetQuery + "&" + request.URL.RawQuery
		}
		if _, ok := request.Header["User-Agent"]; !ok {
			request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
		}
	}

	proxy.ServeHTTP(res, req)
}

func main() {
	root := os.Getenv("CNS_DATA")
	ninja := os.Getenv("CNS_NINJA")
	if root == "" || ninja == "" {
		log.Fatal("CNS_DATA or CNS_NINJA is empty")
	}
	if os.Getenv("CNS_OFFLINE_CODE") != "" {
		offlineCode = os.Getenv("CNS_OFFLINE_CODE")
	} else {
		offlineCode = randStr(10)
		Warn("CNS_OFFLINE_CODE not set, will use", offlineCode)
	}
	http.HandleFunc("/", handleProxy)
	Debug("CNS_DATA", root)
	Debug("CNS_NINJA", ninja)
	Debug("Starting proxy server on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
