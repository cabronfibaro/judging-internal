import { defineEnvVars } from '@sveltejs/kit/env';

export const variables = defineEnvVars({
	CHALLONGE_API_KEY: {
		static: true
	}
});
