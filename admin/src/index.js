import { Hono } from 'hono'
import base from './service/base'
import openai from './service/openai'
import user from './service/user'
import conversation from './service/conversation'
import { verify } from './utils/jwt'
import { serveStatic } from 'hono/bun'

const app = new Hono()

const openApi = [
    '/api/user/login'
]

app.use('*', async (c, next) => {
    if (openApi.includes(c.req.path)) {
        return await next()
    }
    const token = c.req.header('x-token')

    if (!token) {
        return c.json({
            code: 401
        })
    }

    const paylod = await verify(token)

    if (!paylod) {
        return c.json({
            code: 401,
            token
        })
    }

    await next()
})

base(app)
openai(app)
user(app)
conversation(app)

const root = new Hono()
root.route('/api', app)
root.use('*', serveStatic({ root: './ui' }))

export default {
    port: 3001,
    fetch: root.fetch,
}  
