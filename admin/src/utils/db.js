import { Database } from "bun:sqlite";

const db = new Database(Bun.env.CNS_DATA + "/db.sqlite", { create: true });

export default db