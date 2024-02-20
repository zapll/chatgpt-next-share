import Entity from "../utils/entity"
import db from "../utils/db"
import sql from 'sql-bricks'
import { timeDifference, dateFormat } from '../utils/func'

const project = 'openai'
const schema = {
    "tableProps": {
        "page": {
            "ps": 100,
            "sizes": [100, 200]
        }
    },
    "filter": [
        {
            "field": "desc",
            "label": "车号",
            "operator": "like"
        },
    ],
    "afterSelect": [
        function (rows) {
            const count = db.query(`select oid,count(1) as count from user group by oid`).all()
            for (let row of rows) {
                for (let c of count) {
                    if (row.id === c.oid) {
                        row.ref_num = c.count
                    }
                }
            }
            return rows
        },
    ],
    "orderBy": {
        "field": "id",
        "mod": "desc"
    },
    "formItems": [
        {
            "field": "desc",
            "label": "车号",
            "type": "input"
        },
        {
            "field": "session_token",
            "label": "session_token",
        },
    ],
    "headers": [
        {
            "field": "id",
            "label": "ID"
        },
        {
            "field": "desc",
            "label": "车号",
            "type": "input"
        },
        {
            "field": "session_token_exp",
            "label": "Cookie 过期时间",
            "handler": (value, row) => {
                return value ? timeDifference((new Date(value)).toUTCString()) : ''
            }
        },
        {
            "label": "订阅数",
            "fake": true,
            "field": 'ref_num'
        },
        {
            "field": "created_at",
            "label": "创建时间",
            "handler": (value, row) => {
                return value ? dateFormat(new Date(value)) : ''
            }
        },
    ],
    "normalButton": [
        {
            "target": "/openai/form",
            "text": "新建",
            "type": "jump"
        }
    ],
    "rowButton": [
        {
            "props": {
                "type": "primary"
            },
            "target": "/openai/{id}",
            "text": "编辑",
            "type": "jump"
        },
    ],
    "beforeCreate": [
        function (data) {
            return refreshOne(data)
        }
    ],
    "afterUpdate": [
        function (row) {
            refreshOne(row)
        }
    ]
}

const refreshOne = async (data) => {
    const sessionToken = data.session_token

    var myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + sessionToken);

    var requestOptions = {
        method: 'POST',
        headers: myHeaders,
    };

    const response = await fetch(Bun.env.CNS_NINJA + "/auth/refresh_session", requestOptions)
    if (response.status !== 200) {
        return new Error(response.statusText)
    }

    const result = await response.json()

    if (!result || result.code) {
        console.log("refresh error")
        return
    }

    if (data.id) {
        const _sql = sql.update('openai', {
            session_token_exp: result.expires
        }).where({id: data.id}).toString()
    
        db.run(_sql)
    }

    data.session_token_exp = result.expires
    return data 
}

export default function (app) {
    new Entity(project, schema).reg(app)
}