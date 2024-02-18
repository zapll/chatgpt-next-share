import Entity from "../utils/entity"

const project = 'user'
const schema = {
    "tableProps": {
        "page": {
            "ps": 100,
            "sizes": [100, 200]
        }
    },
    "join": [
        {
            "table": "openai",
            "foreign_key": "oid",
            "local_key": "id",
            "select": ['desc']
        },
    ],
    "orderBy": {
        "field": "id",
        "mod": "desc"
    },
    "formItems": [
        {
            "field": "name",
            "label": "名称",
        },
        {
            "field": "token",
            "label": "登录秘钥",
            "type": "user_key"
        },
        {
            "field": "expire_at",
            "label": "过期时间",
            "type": "datetime"
        },
        {
            "field": "oid",
            "label": "主账号",
            "type": "select",
            "props": {
                "valueKey": "id",
                "labelKey": "desc",
                "selectApi": "/openai/options?field=desc",
            }
        },
        {
            "type": "VShow",
            "props": {
                "tpl": "登录地址 your.chatgpt.com <br/>登录秘钥: {token}"
            }
        }
    ],
    "filter": [
        {
            "field": "token",
            "label": "登录秘钥",
        },
        {
            "field": "oid",
            "label": "主账号",
            "type": "select",
            "props": {
                "valueKey": "id",
                "labelKey": "desc",
                "selectApi": "/openai/options?field=desc",
            }
        },
    ],
    "headers": [
        {
            "field": "id",
            "label": "ID"
        },
        {
            "field": "name",
            "label": "名称"
        },
        {
            "field": "token",
            "label": "登录秘钥",
        },
        {
            "field": "oid",
            "hidden": true,
        },
        {
            "field": "openai.desc",
            "label": "订阅",
            "fake": true,
        },
        {
            "field": "expire_at",
            "label": "过期时间",
        },
    ],
    "normalButton": [
        {
            "target": "/user/form",
            "text": "新建",
            "type": "jump"
        }
    ],
    "rowButton": [
        {
            "props": {
                "type": "primary"
            },
            "target": "/user/{id}",
            "text": "编辑",
            "type": "jump"
        },
    ],
}

export default function (app) {
   new Entity(project, schema).reg(app)
}