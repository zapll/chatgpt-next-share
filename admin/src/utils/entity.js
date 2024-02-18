import db from "./db"
import sql from 'sql-bricks'
import dayjs from "dayjs";

var select = sql.select, insert = sql.insert, update = sql.update;
var _in = sql.in, or = sql.or, like = sql.like, lt = sql.lt, eq = sql.eq, gt = sql.gt, gte = sql.gte, lte = sql.lte;

var handleValue = sql._handleValue;

// Insert & Update OR clauses (SQLite dialect)
update.defineClause('or', function (opts) { return this._or ? `OR ${this._or}` : '' }, { after: 'update' });
insert.defineClause('or', function (opts) { return this._or ? `OR ${this._or}` : '' }, { after: 'insert' });

var or_methods = {
    'orReplace': 'REPLACE', 'orRollback': 'ROLLBACK',
    'orAbort': 'ABORT', 'orFail': 'FAIL'
};
Object.keys(or_methods).forEach(function (method) {
    insert.prototype[method] = update.prototype[method] = function () {
        this._or = or_methods[method]; return this;
    };
});

select.prototype.limit = function (val) {
    this._limit = val;
    return this;
};
select.prototype.offset = function (val) {
    this._offset = val;
    return this;
};

select.defineClause(
    'limit',
    function (opts) { return this._limit != null ? `LIMIT ${handleValue(this._limit, opts)}` : '' },
    { after: 'orderBy' }
);

select.defineClause(
    'offset',
    function (opts) { return this._offset != null ? `OFFSET ${handleValue(this._offset, opts)}` : '' },
    { after: 'limit' }
);

class Entity {
    project;
    schema;
    constructor(project, schema) {
        this.project = project
        this.schema = schema
    }

    makeFilter = (c, query) => {
        if (this.schema.filter) {
            for (let item of this.schema.filter) {
                const value = c.req.query(item.field)
                if (value === undefined) {
                    continue;
                }
                item.operator = item.operator || '='
                switch (item.operator) {
                    case '=':
                        query = query.where(item.field, value)
                        break;
                    case '>':
                        query = query.where(gt(item.field, value))
                        break;
                    case '<':
                        query = query.where(lt(item.field, value))
                        su.lt(item.field, value)
                        break;
                    case '>=':
                        query = query.where(gte(item.field, value))
                        break;
                    case '<=':
                        query = query.where(lte(item.field, value))
                        break;
                    case 'like':
                        query = query.where(like(item.field, `%${value}%`))
                }
            }
        }
    }

    list = async (c) => {
        const table = select().from(this.project)

        let query = table.select('count(1) as count')
        this.makeFilter(c, query)

        console.log(query.toString())

        const result = db.query(query.toString()).get()

        if (!result) {
            return c.json({
                code: 500,
                message: 'select error'
            })
        }

        const count = result.count
        const page = (c.req.query('_pn') || 1) * 1
        const size = (c.req.query('_ps') || 20) * 1
        if (count == 0) {
            return c.json({
                "code": 0,
                "data": {
                    "page": {
                        "pn": page,
                        "ps": size,
                        "total": 0
                    },
                    "list": []
                }
            })
        }

        let field = []
        for (let item of (this.schema.headers || [])) {
            if (item.fake || !item.field) {
                continue
            }
            field.push(item.field)
        }
        query =  select().from(this.project).select(field)
        this.makeFilter(c, query)
        const from = (page - 1) * size
        query = query.limit(size).offset(from)
        if (this.schema.orderBy) {
            const { field, mod } = this.schema.orderBy
            query.orderBy(field, mod)
        }

        console.log(query.toString())
        let data = db.query(query.toString()).all()

        data = data || []

        for (let h of (this.schema.join || [])) {
            data = ((rows) => {
                let keys = rows.map(e => e[h.foreign_key]);
                let _sql = select().from(h.table).select(h.select.concat([h.local_key])).where(_in(h.local_key, keys)).toString()
                const record = db.query(_sql).all()
                rows = rows.map(row => {
                    for (let item of record) {
                        if (row[h.foreign_key] === item[h.local_key]) {
                            row[h.table] = item
                            return row
                        }
                    }
                    return row
                })

                return rows
            })(data)
        }

        for (let h of (this.schema.afterSelect || [])) {
            data = h(data)
        }

        data.map((row) => {
            for (let field of (this.schema.headers || [])) {
                if (field.field === undefined) {
                    continue
                }
                const value = row[field.field]
                row[field.field] = field.handler ? field.handler(value, row) : value
            }
        })

        return c.json({
            "code": 0,
            "data": {
                "page": {
                    "pn": page,
                    "ps": size,
                    "total": count
                },
                "list": data
            }
        })
    }


    get = async (c) => {
        const id = await c.req.param("id")
        const data = db.query(select().from(this.project).where(eq('id', id)).toString()).get()

        if (!data) {
            return c.json({
                code: 500,
                message: `Not Found`
            })
        }

        return c.json({
            "code": 0,
            "data": data
        })
    }

    body = async (c) => {
        const body = await c.req.json()
        const _body = {}
        for (let item of (this.schema.formItems || [])) {
            if (item.field === undefined) {
                continue
            }
            _body[item.field] = body[item.field]
        }
        return _body
    }

    update = async (c) => {
        const id = await c.req.param("id")
        const body = await this.body(c)

        db.run(update(this.project, body).where({ 'id': id }).toString())

        for (let h of (this.schema.afterUpdate || [])) {
            const res = h({ ...body, id })
            if (res instanceof Promise) {
                await res
            }
        }

        return c.json({
            "code": 0,
        })
    }

    delete = async (c) => {
        const id = await c.req.param("id")
        db.run(delete (this.project).where({ id: id }))

        for (let h of (this.schema.afterDelete || [])) {
            const res = h(body)
            if (res instanceof Promise) {
                await res
            }
        }

        return c.json({
            "code": 0,
        })
    }

    create = async (c) => {
        let body = await c.req.json()
        for (let h of (this.schema.beforeCreate || [])) {
            let res = h(body)
            if (res instanceof Promise) {
                res = await res
            }
            if (res) {
                body = res
            }
        }

        db.run(insert(this.project, body).toString())

        for (let h of (this.schema.afterCreate || [])) {
            const res = h(body)
            if (res instanceof Promise) {
                await res
            }
        }

        return c.json({
            "code": 0,
        })
    }


    selectOptions = async (c) => {
        const kw = await c.req.query('kw')
        const table = select().from(this.project)
        // if kw is number
        if (/^\d+$/.test(kw)) {
            table.where(eq('id', kw))
        } else {
            const field = await c.req.query('field')
            const operatro = (await c.req.query('operator')) || 'like'

            if (field) {
                switch (operatro) {
                    case 'like':
                        table.where(like(field, `%${kw}%`))
                        break
                    case 'eq':
                        table.where(eq(field, kw))
                }
            }
        }

        table.limit(100)

        const data = db.query(table.toString()).all()
        return c.json({
            "code": 0,
            "data": data
        })
    }

    reg(app, fn) {
        app.get(`/schema/${this.project}`, async (c) => {
            return c.json({
                "code": 0,
                "data": this.schema
            })
        })
        app.get(`/${this.project}/list`, this.list)
        app.get(`/${this.project}/get/:id`, this.get)
        app.post(`/${this.project}/update/:id`, this.update)
        app.post(`/${this.project}/del/:id`, this.delete)
        app.post(`/${this.project}/create`, this.create)
        app.get(`/${this.project}/options`, this.selectOptions)
        fn && fn(app)
    }
}

export default Entity
