import Database from 'better-sqlite3';
import { drizzle } from 'drizzle-orm/better-sqlite3';
import { DB_URL } from '$env/static/private';

const sqlite = new Database(DB_URL.replace(/^file:/, ''));
export const db = drizzle(sqlite);
