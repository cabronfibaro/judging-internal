CREATE TABLE `active_battles` (
	`id` text PRIMARY KEY NOT NULL,
	`match_id` text NOT NULL,
	`participant_id` text NOT NULL,
	`lead_status` integer DEFAULT false NOT NULL,
	FOREIGN KEY (`participant_id`) REFERENCES `participants`(`id`) ON UPDATE no action ON DELETE cascade
);
