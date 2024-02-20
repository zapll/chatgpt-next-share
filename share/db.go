package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/daodao97/goadmin/pkg/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cast"
)

func init() {
	file := os.Getenv("CNS_DATA") + "/db.sqlite"
	err := db.Init(map[string]*db.Config{
		"default": {DSN: file, Driver: "sqlite3"},
	})
	if err != nil {
		panic(err)
	}
}

func getUser(req *http.Request) (*db.Row, error) {
	sessionToken, err := req.Cookie("session_token")
	if err != nil {
		return nil, err
	}

	user := db.New("user").SelectOne(db.WhereEq("token", sessionToken.Value))
	if user.Err != nil {
		return nil, user.Err
	}
	return user, nil
}

func createConv(conv db.Record) error {
	row := db.New("conversation").SelectOne(db.WhereEq("cid", conv["cid"]))
	if row.Data != nil {
		return nil
	}
	conv["msg_num"] = 1
	_, err := db.New("conversation").Insert(conv)
	return err
}

func updateConv(cid string, conv db.Record) error {
	_, err := db.New("conversation").Update(conv, db.WhereEq("cid", cid))
	return err
}

func deleteConv(cid string) error {
	_, err := db.New("conversation").Delete(db.WhereEq("cid", cid))
	return err
}

func updateConvMsgNum(cid string) error {
	db, _ := db.DB("default")
	_, err := db.Exec("update conversation set msg_num = msg_num + 1 where cid = ?", cid)
	return err
}

type Item struct {
	Id                     string      `json:"id"`
	Title                  string      `json:"title"`
	CreateTime             string      `json:"create_time"`
	UpdateTime             string      `json:"update_time"`
	Mapping                interface{} `json:"mapping"`
	CurrentNode            interface{} `json:"current_node"`
	ConversationTemplateId interface{} `json:"conversation_template_id"`
	GizmoId                interface{} `json:"gizmo_id"`
	IsArchived             bool        `json:"is_archived"`
	WorkspaceId            interface{} `json:"workspace_id"`
}

type Convs struct {
	Items                   []*Item `json:"items"`
	Total                   int     `json:"total"`
	Limit                   int     `json:"limit"`
	Offset                  int     `json:"offset"`
	HasMissingConversations bool    `json:"has_missing_conversations"`
}

var emptyConvs = &Convs{
	Items:                   []*Item{},
	Total:                   0,
	Limit:                   0,
	Offset:                  0,
	HasMissingConversations: false,
}

func convs(uid string, offset, limit int64) *Convs {
	conv := db.New("conversation", db.ColumnHook(db.Json("profile")))

	count, err := conv.Count()
	if err != nil || count == 0 {
		return emptyConvs
	}

	rows := conv.Select(db.WhereEq("uid", uid), db.OrderByDesc("updated_at"), db.Limit(int(limit)), db.Offset(int(offset)))
	if rows.Err != nil {
		return emptyConvs
	}

	items := []*Item{}
	for _, row := range rows.List {
		items = append(items, &Item{
			Id:         row.GetString("cid"),
			Title:      row.GetString("title"),
			CreateTime: row.GetString("created_at"),
			UpdateTime: row.GetString("updated_at"),
		})
	}

	return &Convs{
		Items:  items,
		Total:  cast.ToInt(count),
		Limit:  cast.ToInt(limit),
		Offset: cast.ToInt(offset),
	}
}

func getCookieByToken(token string) (string, error) {
	token = strings.ReplaceAll(token, "Bearer ", "")
	user := db.New("user").SelectOne(db.WhereEq("token", token))
	if user.Err != nil {
		Error("getCookieByToken user", user.Err)
		return "", user.Err
	}

	oai := db.New("openai").SelectOne(db.WhereEq("id", user.GetString("oid")))
	if oai.Err != nil {
		Error("getCookieByToken openai", oai.Err)
		return "", oai.Err
	}

	return oai.GetString("session_token"), nil
}

func getAccessToken(token string) string {
	token = strings.ReplaceAll(token, "Bearer ", "")
	_at, ok := accessTokenMap.Load(token)
	Debug("getAccessToken", token, ok)
	if !ok {
		return ""
	}
	return _at.(string)
}
