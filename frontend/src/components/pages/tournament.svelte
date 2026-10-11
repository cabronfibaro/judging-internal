<script lang="ts">
  import { useRoute } from "@dvcol/svelte-simple-router/router";
  import { createQuery } from "@tanstack/svelte-query";

  const { location } = $derived(useRoute());
  const id = $derived(location?.params.id);

  const tournamentQuery = createQuery(() => ({
    queryKey: ["tournament", id],
    enabled: !!id,
    queryFn: async () => {
      await new Promise((resolve) => setTimeout(resolve, 1000));

      return {
        id,
        name: `Tournament ${id}`,
        status: "ongoing",
        participants: 16,
      };
    },
  }));
</script>

{#if !id}
  <p>No tournament ID found.</p>
{:else if tournamentQuery.isPending}
  <p>Loading tournament...</p>
{:else if tournamentQuery.isError}
  <p>Failed to load tournament.</p>
{:else}
  <h1>{tournamentQuery.data.name}</h1>
  <p>ID: {tournamentQuery.data.id}</p>
  <p>Status: {tournamentQuery.data.status}</p>
  <p>Participants: {tournamentQuery.data.participants}</p>
{/if}
