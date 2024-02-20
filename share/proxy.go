package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/daodao97/goadmin/pkg/db"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

func rebuildReq(req *http.Request) {
	token, err := req.Cookie("session_token")
	if err != nil {
		Error("Failed to get session_token", err)
		return
	}

	cookie, _ := getCookieByToken(token.Value)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: cookie,
	})

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", getAccessToken(token.Value)))
}

func newRequest(req *http.Request) *resty.Request {
	client := resty.New().R()
	for _, c := range req.Cookies() {
		client.SetCookie(&http.Cookie{
			Name:  c.Name,
			Value: c.Value,
		})
	}

	if req.Header.Get("Authorization") != "" {
		client.SetHeader("Authorization", req.Header.Get("Authorization"))
	}
	return client
}

func openaiUrl(path string) string {
	return os.Getenv("CNS_NINJA") + path
}

func openai(res http.ResponseWriter, req *http.Request) (*httputil.ReverseProxy, *url.URL, error) {
	_url, _ := url.Parse(os.Getenv("CNS_NINJA"))
	proxy := httputil.NewSingleHostReverseProxy(_url)
	rebuildReq(req)

	if startWith(req.URL.Path, "/backend-api", "/api") && req.URL.Path != "/backend-api/register-websocket" {
		if req.Header.Get("Authorization") == "Bearer " {
			return nil, nil, errors.New("authorization is empty")
		}

		if startWith(req.URL.Path, "/backend-api/conversation/gen_title") {
			handelGenTitle(res, req, proxy)
		}

		if req.URL.Path == "/backend-api/conversation" && req.Method == http.MethodPost {
			handelConversation(res, req, proxy)
		}

		if startWith(req.URL.Path, "/backend-api/conversation/gen_title") {
			handelGenTitle(res, req, proxy)
		}

		if startWith(req.URL.Path, "/backend-api/conversation/") && req.Method == http.MethodPatch {
			cid := strings.ReplaceAll(req.URL.Path, "/backend-api/conversation/", "")
			deleteConv(cid)
		}
	}

	return proxy, _url, nil
}

func handelConversation(res http.ResponseWriter, req *http.Request, proxy *httputil.ReverseProxy) {
	msg, err1 := readAndReplaceBody(req)
	user, err2 := getUser(req)
	if err1 != nil || err2 != nil {
		return
	}
	cid := gjson.Get(string(msg), "conversation_id").String()
	if cid != "" {
		proxy.ModifyResponse = func(response *http.Response) error {
			if response.StatusCode == 200 {
				_ = updateConvMsgNum(cid)
				updateConvTitle(cid, req)
			}
			return nil
		}
		return
	}
	proxy.ModifyResponse = func(response *http.Response) error {
		if strings.Contains(response.Header.Get("Content-Type"), "text/event-stream") {
			originalBody := response.Body
			pipedReader, pipedWriter := io.Pipe()
			// 使用 io.Pipe() 来创建一个读写管道
			// 必须确保在读取数据的同时将数据写回到管道中，这样客户端才能接收到数据
			go func() {
				defer originalBody.Close()
				scanner := bufio.NewScanner(originalBody)
				created := false
				for scanner.Scan() {
					line := scanner.Text()
					if strings.HasPrefix(line, "data:") {
						_line := line[5:]
						if gjson.Get(_line, "conversation_id").String() != "" && !created {
							_ = createConv(db.Record{
								"cid":   gjson.Get(_line, "conversation_id").String(),
								"uid":   user.GetString("id"),
								"title": "New chat",
							})
							created = true
						}
						if gjson.Get(_line, "type").String() == "title_generation" {
							_ = updateConv(gjson.Get(_line, "conversation_id").String(), db.Record{
								"title": gjson.Get(_line, "title").String(),
							})
						}
					}
					// 将读取到的数据写回管道
					fmt.Fprintln(pipedWriter, line)
				}
				pipedWriter.Close()
			}()

			// 使用包装过的响应体，这样原始的响应体就不会在ModifyResponse函数结束时关闭
			response.Body = pipedReader
		}

		if strings.Contains(response.Header.Get("Content-Type"), "application/json") {
			bodyBytes, err := io.ReadAll(response.Body)
			if err != nil {
				Error("handelConversation Failed to read response body", err)
				return nil
			}

			ce := response.Header.Get("Content-Encoding")
			body, _ := bodyCompare(ce, bodyBytes)
			cid := gjson.GetBytes(body, "conversation_id").String()
			if cid != "" {
				_ = createConv(db.Record{
					"cid":   cid,
					"uid":   user.GetString("id"),
					"title": "New chat",
				})
				go func() {
					time.Sleep(time.Second * 1)
					updateConvTitle(cid, req)
				}()
			}

			response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		return nil
	}
}

func updateConvTitle(cid string, req *http.Request) {
	Debug("gen_title", cid)
	client := newRequest(req)
	resp, err := client.
		Get(openaiUrl("/backend-api/conversation/" + cid))
	if err != nil {
		Error("Failed to get conversation detail", err)
		return
	}
	_body := resp.Body()
	title := gjson.GetBytes(_body, "title").String()
	if title == "" {
		Warn("title is empty", cid, string(_body))
		return
	}
	_ = updateConv(cid, db.Record{
		"title": title,
	})
}

func handleConversations(res http.ResponseWriter, req *http.Request) {
	user, err := getUser(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		Error("Failed to get user", err)
		return
	}

	limit := cast.ToInt64(req.URL.Query().Get("limit"))
	offset := cast.ToInt64(req.URL.Query().Get("offset"))

	list := convs(user.GetString("id"), offset, limit)

	tag := false
	for _, item := range list.Items {
		if item.Title == "New chat" && item.Id != "" {
			updateConvTitle(item.Id, req)
			tag = true
		}
	}
	if tag {
		list = convs(user.GetString("id"), offset, limit)
	}

	res.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(res).Encode(list)
}

func handelGenTitle(res http.ResponseWriter, req *http.Request, proxy *httputil.ReverseProxy) {
	user, err := getUser(req)
	if err != nil {
		return
	}

	proxy.ModifyResponse = func(response *http.Response) error {
		title, err := genTitle(response)
		if err != nil {
			Error("gen title error", err)
			return nil
		}
		_ = createConv(db.Record{
			"cid":   strings.ReplaceAll(req.URL.Path, "/backend-api/conversation/gen_title/", ""),
			"uid":   user.GetString("id"),
			"title": title,
		})
		return nil
	}
}

func genTitle(response *http.Response) (string, error) {
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	ce := response.Header.Get("Content-Encoding")
	body, err := bodyCompare(ce, bodyBytes)
	if err != nil {
		return "", err
	}

	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	Debug("/gen_title response", ce, string(body))

	title := gjson.GetBytes(body, "title").String()
	if title != "" {
		return title, nil
	}

	message := gjson.GetBytes(body, "message").String()
	if !strings.Contains(message, "already has title") {
		return "", err
	}

	_title := strings.Trim(strings.Split(message, "already has title ")[1], "'")

	return _title, nil
}
