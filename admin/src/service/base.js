
import { gen, verify } from "../utils/jwt"
import db from "../utils/db"

export default function (app) {
  app.post('/user/login', async (c) => {
    const { username, password } = await c.req.json()
    const result = db.query(`select * from operator where username=$username and password=$password`).get({ $username: username, $password: password })

    if (!result) {
      return c.json({
        "code": 500,
        "data": null
      })
    }

    return c.json({
      "code": 0,
      "data": {
        "name": username,
        "token": await gen({ username })
      }
    })
  })

  app.get('/user/info', async (c) => {
    const token = c.req.header('x-token')

    if (!token) {
      return c.json({
        code: 401
      })
    }

    const paylod = await verify(token)
    return c.json({
      "code": 0,
      "data": {
        "id": 1,
        "name": paylod.username,
        "resource": null,
        "env": "prod",
        "website": {
          "title": "GptAdmin"
        }
      }
    })
  })

  app.get('/user/routes', async (c) => {
    return c.json({
      "code": 0,
      "data": [
        {
          "id": 0,
          "routes": [
            {
              "module_id": 0,
              "name": "账号管理",
              "type": 2,
              "path": "/openai",
              "icon": "ra-office-supplies",
              "page_type": 7,
            },
            {
              "module_id": 0,
              "name": "用户管理",
              "type": 2,
              "path": "/user",
              "icon": "ra-office-supplies",
              "page_type": 7,
            },
            {
              "module_id": 0,
              "name": "会话管理",
              "type": 2,
              "path": "/conversation",
              "icon": "ra-office-supplies",
              "page_type": 7,
            },
          ]
        }
      ]
    })
  })


  app.get('/user/logout', async (c) => {
    return c.json({
      "code": 0,
      "data": null
    })
  })

  app.get("/form_mutex", async (c) => {
    return c.json({
      code: 0
    })
  })
}