PRAGMA foreign_keys=OFF;--> statement-breakpoint
CREATE TABLE `__new_active_battles` (
	`id` text PRIMARY KEY NOT NULL,
	`match_id` text NOT NULL,
	`participant_id` text NOT NULL,
	`lead_status` integer DEFAULT 0 NOT NULL,
	FOREIGN KEY (`participant_id`) REFERENCES `participants`(`id`) ON UPDATE no action ON DELETE cascade
);
--> statement-breakpoint
INSERT INTO `__new_active_battles`("id", "match_id", "participant_id", "lead_status") SELECT "id", "match_id", "participant_id", "lead_status" FROM `active_battles`;--> statement-breakpoint
DROP TABLE `active_battles`;--> statement-breakpoint
ALTER TABLE `__new_active_battles` RENAME TO `active_battles`;--> statement-breakpoint
PRAGMA foreign_keys=ON;