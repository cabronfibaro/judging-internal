<script lang="ts">
  import { AlertCircle, RefreshCw, Trophy, Plus } from "@lucide/svelte";
  import { createMutation, useQueryClient } from "@tanstack/svelte-query";

  import { Button } from "$src/components/ui/button";
  import {
    Empty,
    EmptyDescription,
    EmptyHeader,
    EmptyMedia,
    EmptyTitle,
  } from "$src/components/ui/empty";
  import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
  } from "$src/components/ui/dialog";
  import { Badge } from "$src/components/ui/badge";

  import CardSkel from "./card-skel.svelte";
  import Card from "./card.svelte";
  import type { tournament } from "$wails/go/models";
  import { DeleteTournament } from "$wails/go/main/App";

  let {
    tournaments = [],
    isLoading = false,
    isError = false,
    error = null,
    onRetry,
    onCreateTournament,
  }: {
    tournaments: tournament.Tournament[];
    isLoading?: boolean;
    isError?: boolean;
    error?: Error | null;
    onRetry?: () => void;
    onCreateTournament?: () => void;
  } = $props();

  const queryClient = useQueryClient();

  const deleteTournamentMutation = createMutation(() => ({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      DeleteTournament(id, status),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["tournaments"],
      });
      deleteDialogOpen = false;
      selectedTournament = null;
    },
    onError: (err) => {
      console.error(`Error while deleting tournament: ${err}`);
    },
  }));

  let deleteDialogOpen = $state(false);
  let selectedTournament = $state<tournament.Tournament | null>(null);

  function openDeleteDialog(tournament: tournament.Tournament) {
    selectedTournament = tournament;
    deleteDialogOpen = true;
  }

  function onDeleteTournament() {
    if (!selectedTournament) return;
    const { id, status } = selectedTournament;
    deleteTournamentMutation.mutate({ id, status });
  }

  const tournamentCountLabel = $derived(
    tournaments.length === 1
      ? "1 tournament"
      : `${tournaments.length} tournaments`,
  );
  const activeCount = $derived(
    tournaments.filter((t) => t.status === "elimination").length,
  );
  const qualifyingCount = $derived(
    tournaments.filter((t) => t.status === "qualifying").length,
  );
</script>

<!-- Delete Dialog -->
<Dialog bind:open={deleteDialogOpen}>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Delete tournament?</DialogTitle>
      <DialogDescription>
        This will permanently delete
        <strong>{selectedTournament?.name}</strong>
        and its qualifier records. This action cannot be undone.
      </DialogDescription>
    </DialogHeader>

    <DialogFooter>
      <Button variant="outline">Cancel</Button>
      <Button
        disabled={deleteTournamentMutation.isPending}
        onclick={onDeleteTournament}
        variant="destructive"
      >
        {deleteTournamentMutation.isPending ? "Deleting..." : "Delete"}
      </Button>
    </DialogFooter>
  </DialogContent>
</Dialog>

<!-- Content Views -->
{#if isLoading}
  <div class="w-full space-y-4" aria-label="Loading tournaments">
    {#each Array.from({ length: 4 }, (_, i) => i) as idx (idx)}
      <CardSkel />
    {/each}
  </div>
{:else if isError}
  <Empty class="min-h-64 border border-dashed border-destructive/40">
    <EmptyHeader>
      <EmptyMedia variant="icon"><AlertCircle aria-hidden="true" /></EmptyMedia>
      <EmptyTitle>We couldn't load tournaments</EmptyTitle>
      <EmptyDescription>
        {error?.message || "The tournament list could not be loaded."}
      </EmptyDescription>
    </EmptyHeader>
    {#if onRetry}
      <Button variant="outline" onclick={onRetry}>
        <RefreshCw aria-hidden="true" />
        Try again
      </Button>
    {/if}
  </Empty>
{:else if tournaments.length === 0}
  <Empty class="min-h-80 border border-dashed">
    <EmptyHeader>
      <EmptyMedia variant="icon"><Trophy aria-hidden="true" /></EmptyMedia>
      <EmptyTitle>Your competition board is empty</EmptyTitle>
      <EmptyDescription>
        Create your first tournament to start qualifying drivers and running
        head-to-head battles.
      </EmptyDescription>
    </EmptyHeader>
    {#if onCreateTournament}
      <Button onclick={onCreateTournament}>
        <Plus aria-hidden="true" />
        Create tournament
      </Button>
    {/if}
  </Empty>
{:else}
  <div class="w-full space-y-6">
    <div
      class="flex flex-col gap-4 border-b border-border pb-5 sm:flex-row sm:items-end sm:justify-between"
    >
      <div class="min-w-0">
        <div class="flex flex-wrap items-center gap-3">
          <h2 class="text-lg font-semibold tracking-tight">
            Competition board
          </h2>
          <Badge variant="outline">{tournamentCountLabel}</Badge>
        </div>
        <p class="mt-1 max-w-2xl text-sm text-muted-foreground">
          A live view of every tournament and its bracket progress.
        </p>
      </div>
      <div
        class="flex shrink-0 flex-wrap items-center gap-2 text-xs text-muted-foreground"
      >
        <span class="rounded-full bg-primary/10 px-2.5 py-1 text-primary">
          {activeCount} active
        </span>
        <span class="rounded-full bg-muted px-2.5 py-1">
          {qualifyingCount} qualifying
        </span>
      </div>
    </div>

    <div class="w-full space-y-3">
      {#each tournaments as tournament, index (tournament.id)}
        <Card {tournament} {openDeleteDialog} position={index + 1} />
      {/each}
    </div>
  </div>
{/if}
