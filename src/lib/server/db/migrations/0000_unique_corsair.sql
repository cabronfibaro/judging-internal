CREATE TABLE `brackets` (
	`id` text PRIMARY KEY NOT NULL,
	`tournament_id` text NOT NULL,
	`challonge_id` text NOT NULL,
	`url` text NOT NULL,
	`status` text DEFAULT 'not_started' NOT NULL,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL,
	FOREIGN KEY (`tournament_id`) REFERENCES `tournaments`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `entrants` (
	`id` text PRIMARY KEY NOT NULL,
	`tournament_id` text NOT NULL,
	`participant_id` text NOT NULL,
	`status` text DEFAULT 'pending' NOT NULL,
	`score` real DEFAULT 0 NOT NULL,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL,
	FOREIGN KEY (`tournament_id`) REFERENCES `tournaments`(`id`) ON UPDATE no action ON DELETE cascade,
	FOREIGN KEY (`participant_id`) REFERENCES `participants`(`id`) ON UPDATE no action ON DELETE no action
);
--> statement-breakpoint
CREATE TABLE `nodes` (
	`id` text PRIMARY KEY NOT NULL,
	`bracket_id` text NOT NULL,
	`participant_id` text,
	`match_id` text NOT NULL,
	`status` text NOT NULL,
	`round` integer NOT NULL,
	FOREIGN KEY (`bracket_id`) REFERENCES `brackets`(`id`) ON UPDATE no action ON DELETE cascade,
	FOREIGN KEY (`participant_id`) REFERENCES `participants`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `participants` (
	`id` text PRIMARY KEY NOT NULL,
	`name` text NOT NULL,
	`challonge_participant_id` text,
	`tournament_id` text NOT NULL,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL,
	FOREIGN KEY (`tournament_id`) REFERENCES `tournaments`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
CREATE TABLE `tournaments` (
	`id` text PRIMARY KEY NOT NULL,
	`name` text NOT NULL,
	`size` integer NOT NULL,
	`starting_date` integer NOT NULL,
	`track_id` text NOT NULL,
	`qualifying_order` text NOT NULL,
	`randomize_qualifying_order` integer DEFAULT true NOT NULL,
	`status` text DEFAULT 'qualifying' NOT NULL,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL,
	FOREIGN KEY (`track_id`) REFERENCES `tracks`(`id`) ON UPDATE no action ON DELETE no action,
	CONSTRAINT "size_power_of_2" CHECK("tournaments"."size" > 0 AND ("tournaments"."size" & ("tournaments"."size" - 1)) = 0)
);
--> statement-breakpoint
CREATE TABLE `tracks` (
	`id` text PRIMARY KEY NOT NULL,
	`name` text NOT NULL,
	`config` text,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL
);
