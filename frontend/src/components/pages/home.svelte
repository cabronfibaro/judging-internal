<script lang="ts">
  import { createQuery } from "@tanstack/svelte-query";
  import Main from "../layout/main.svelte";
  import CreateDialog from "../tournaments/create-dialog.svelte";
  import Grid from "../tournaments/grid.svelte";

  import { GetTournaments } from "$wails/go/main/App";

  let createTournamentOpen = $state(false);

  const tournamentsQuery = createQuery(() => ({
    queryKey: ["tournaments"],
    queryFn: GetTournaments,
  }));
</script>

<Main>
  <div class="w-full space-y-8">
    <header
      class="flex flex-col gap-5 border-b border-border pb-6 sm:flex-row sm:items-end sm:justify-between"
    >
      <div class="space-y-2">
        <p class="text-sm font-medium text-muted-foreground">
          Tournament management
        </p>
        <h1 class="text-3xl font-semibold tracking-tight">Tournaments</h1>
        <p class="max-w-xl text-sm text-muted-foreground">
          Create a competition, manage qualifying participants, and open its
          bracket when you are ready.
        </p>
      </div>
      <CreateDialog bind:open={createTournamentOpen} />
    </header>

    <Grid
      tournaments={tournamentsQuery.data ?? []}
      isLoading={tournamentsQuery.isPending}
      isError={tournamentsQuery.isError}
      error={tournamentsQuery.error}
      onRetry={() => tournamentsQuery.refetch()}
    />
  </div>
</Main>
