import Entity from "../utils/entity"
const project = 'conversation'
const schema = {
    "orderBy": {
        "field": "updated_at",
        "mod": "desc"
    },
    "tableProps": {
        "page": {
            "ps": 100,
            "sizes": [100, 200]
        }
    },
    "headers": [
        {
            "field": "id",
            "label": "ID",
        },
        {
            "field": "cid",
            "label": "会话ID",
        },
        {
            "field": "uid",
            "label": "用户ID",
        },
        {
            "field": "title",
            "label": "标题",
        },
        {
            "field": "msg_num",
            "label": "消息数",
        },
        {
            "field": "created_at",
            "label": "创建时间",
        },
        {
            "field": "updated_at",
            "label": "更新时间",
        },
    ],
    "rowButton": [
        {
            "props": {
                "type": "danger",

            },
            "extra": {
                "method": "POST"
            },
            "target": "/conversation/del/{id}",
            "text": "删除",
            "type": "api"
        },
    ]
}

export default function (app) {
    new Entity(project, schema).reg(app)
}