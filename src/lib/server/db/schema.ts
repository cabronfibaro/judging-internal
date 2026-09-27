import { relations } from 'drizzle-orm';
import { integer, real, sqliteTable, text } from 'drizzle-orm/sqlite-core';
import { nanoid } from 'nanoid';

const timestamps = {
	createdAt: integer('created_at', { mode: 'timestamp' })
		.notNull()
		.$defaultFn(() => new Date()),

	updatedAt: integer('updated_at', { mode: 'timestamp' })
		.notNull()
		.$defaultFn(() => new Date())
};

export const participants = sqliteTable('participants', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => nanoid(15)),

	name: text('name').notNull(),

	...timestamps
});

export const tracks = sqliteTable('tracks', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => nanoid(15)),

	name: text('name').notNull(),

	// JSON configuration for the track
	config: text('config', { mode: 'json' }).$type<Record<string, unknown>>(),

	...timestamps
});

export const tournaments = sqliteTable('tournaments', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => nanoid(15)),

	name: text('name').notNull(),

	startingDate: integer('starting_date', { mode: 'timestamp' }).notNull(),

	trackId: text('track_id')
		.notNull()
		.references(() => tracks.id),

	...timestamps
});

export const qualifiers = sqliteTable('qualifiers', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => nanoid(15)),

	tournamentId: text('tournament_id')
		.notNull()
		.references(() => tournaments.id),

	participantId: text('participant_id')
		.notNull()
		.references(() => participants.id),

	status: text('status', {
		enum: ['qualified', 'not_qualified', 'pending']
	})
		.notNull()
		.default('not_qualified'),

	score: real('score').notNull().default(0),

	...timestamps
});

export const participantsRelations = relations(participants, ({ many }) => ({
	qualifiers: many(qualifiers)
}));

export const tracksRelations = relations(tracks, ({ many }) => ({
	tournaments: many(tournaments)
}));

export const tournamentsRelations = relations(tournaments, ({ one, many }) => ({
	track: one(tracks, {
		fields: [tournaments.trackId],
		references: [tracks.id]
	}),

	qualifiers: many(qualifiers)
}));

export const qualifiersRelations = relations(qualifiers, ({ one }) => ({
	tournament: one(tournaments, {
		fields: [qualifiers.tournamentId],
		references: [tournaments.id]
	}),

	participant: one(participants, {
		fields: [qualifiers.participantId],
		references: [participants.id]
	})
}));
