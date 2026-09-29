import { sql } from 'drizzle-orm';
import { integer, real, sqliteTable, text, check, int } from 'drizzle-orm/sqlite-core';
import { nanoid } from 'nanoid';

const timestamps = {
	createdAt: integer('created_at', { mode: 'timestamp' })
		.notNull()
		.$defaultFn(() => new Date()),
	updatedAt: integer('updated_at', { mode: 'timestamp' })
		.notNull()
		.$defaultFn(() => new Date())
};

export const tracks = sqliteTable('tracks', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => nanoid(15)),
	name: text('name').notNull(),
	config: text('config', { mode: 'json' }).$type<Record<string, unknown>>(),
	...timestamps
});

export const tournaments = sqliteTable(
	'tournaments',
	{
		id: text('id')
			.primaryKey()
			.$defaultFn(() => nanoid(15)),
		name: text('name').notNull(),
		size: integer('size').notNull(),
		startingDate: integer('starting_date', { mode: 'timestamp' }).notNull(),
		trackId: text('track_id')
			.notNull()
			.references(() => tracks.id),
		qualifyingOrder: text('qualifying_order').notNull(),
		randomizeQualifyingOrder: integer('randomize_qualifying_order', {
			mode: 'boolean'
		})
			.notNull()
			.default(true),
		status: text('status', {
			enum: ['qualifying', 'elimination']
		})
			.notNull()
			.default('qualifying'),
		...timestamps
	},
	(table) => [
		check('size_power_of_2', sql`${table.size} > 0 AND (${table.size} & (${table.size} - 1)) = 0`)
	]
);

export const participants = sqliteTable('participants', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => nanoid(15)),
	name: text('name').notNull(),
	challongeParticipantId: text('challonge_participant_id'), // NULLABLE
	tournamentId: text('tournament_id')
		.notNull()
		.references(() => tournaments.id, {
			onDelete: 'cascade'
		}),
	...timestamps
});

export const entrants = sqliteTable('entrants', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => nanoid(15)),
	tournamentId: text('tournament_id')
		.notNull()
		.references(() => tournaments.id, {
			onDelete: 'cascade'
		}),
	participantId: text('participant_id')
		.notNull()
		.references(() => participants.id),
	status: text('status', {
		enum: ['qualified', 'not_qualified', 'pending', 'scored']
	})
		.notNull()
		.default('pending'),

	score: real('score').notNull().default(0),
	...timestamps
});

export const brackets = sqliteTable('brackets', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => nanoid(15)),
	tournamentId: text('tournament_id')
		.notNull()
		.references(() => tournaments.id, {
			onDelete: 'cascade'
		}),
	challongeId: text('challonge_id').notNull(),
	url: text('url').notNull(),
	status: text('status', { enum: ['not_started', 'started'] })
		.notNull()
		.default('not_started'),
	...timestamps
});

export const nodes = sqliteTable('nodes', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => nanoid(15)),
	bracketId: text('bracket_id')
		.notNull()
		.references(() => brackets.id, {
			onDelete: 'cascade'
		}),
	participantId: text('participant_id').references(() => participants.id, {
		onDelete: 'cascade'
	}),
	matchId: text('match_id').notNull(),
	status: text('status', { enum: ['pending', 'open'] }).notNull(),
	round: int('round').notNull()
});
