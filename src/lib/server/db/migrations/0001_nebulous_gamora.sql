PRAGMA foreign_keys=OFF;--> statement-breakpoint
CREATE TABLE `__new_tournaments` (
	`id` text PRIMARY KEY NOT NULL,
	`name` text NOT NULL,
	`size` integer NOT NULL,
	`starting_date` integer NOT NULL,
	`track_id` text NOT NULL,
	`qualifying_order` text NOT NULL,
	`randomize_qualifying_order` integer DEFAULT true NOT NULL,
	`created_at` integer NOT NULL,
	`updated_at` integer NOT NULL,
	FOREIGN KEY (`track_id`) REFERENCES `tracks`(`id`) ON UPDATE no action ON DELETE no action,
	CONSTRAINT "size_power_of_2" CHECK("__new_tournaments"."size" > 0 AND ("__new_tournaments"."size" & ("__new_tournaments"."size" - 1)) = 0)
);
--> statement-breakpoint
INSERT INTO `__new_tournaments`("id", "name", "size", "starting_date", "track_id", "qualifying_order", "randomize_qualifying_order", "created_at", "updated_at") SELECT "id", "name", "size", "starting_date", "track_id", "qualifying_order", "randomize_qualifying_order", "created_at", "updated_at" FROM `tournaments`;--> statement-breakpoint
DROP TABLE `tournaments`;--> statement-breakpoint
ALTER TABLE `__new_tournaments` RENAME TO `tournaments`;--> statement-breakpoint
PRAGMA foreign_keys=ON;