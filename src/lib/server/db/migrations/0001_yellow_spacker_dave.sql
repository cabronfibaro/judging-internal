PRAGMA foreign_keys=OFF;--> statement-breakpoint
CREATE TABLE `__new_qualifiers` (
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
INSERT INTO `__new_qualifiers`("id", "tournament_id", "participant_id", "status", "score", "created_at", "updated_at") SELECT "id", "tournament_id", "participant_id", "status", "score", "created_at", "updated_at" FROM `qualifiers`;--> statement-breakpoint
DROP TABLE `qualifiers`;--> statement-breakpoint
ALTER TABLE `__new_qualifiers` RENAME TO `qualifiers`;--> statement-breakpoint
PRAGMA foreign_keys=ON;