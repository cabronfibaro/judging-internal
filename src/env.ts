import { defineEnvVars } from '@sveltejs/kit/env';

export const variables = defineEnvVars({
	DB_URL: {
		static: true
	},
	CHALLONGE_API_KEY: {
		static: true
	}
});
