<script lang="ts">
  import EllipsisVertical from "@lucide/svelte/icons/ellipsis-vertical";
  import { Badge } from "$src/components/ui/badge";
  import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
  } from "$src/components/ui/dropdown-menu";
  import { Button } from "$src/components/ui/button";
  import { createMutation, useQueryClient } from "@tanstack/svelte-query";
  import {
    UpdateTournament,
    UpdateEntrant,
    UpdateTournamentState,
    SyncTournament,
    ConcludeQualifers,
  } from "$wails/go/main/App";
  import type { entrant, tournament } from "$wails/go/models";
  import AlertDialog from "../ui/alert-dialog/alert-dialog.svelte";
  import AlertDialogTrigger from "../ui/alert-dialog/alert-dialog-trigger.svelte";
  import AlertDialogContent from "../ui/alert-dialog/alert-dialog-content.svelte";
  import AlertDialogHeader from "../ui/alert-dialog/alert-dialog-header.svelte";
  import AlertDialogTitle from "../ui/alert-dialog/alert-dialog-title.svelte";
  import AlertDialogDescription from "../ui/alert-dialog/alert-dialog-description.svelte";
  import AlertDialogFooter from "../ui/alert-dialog/alert-dialog-footer.svelte";
  import AlertDialogCancel from "../ui/alert-dialog/alert-dialog-cancel.svelte";
  import AlertDialogAction from "../ui/alert-dialog/alert-dialog-action.svelte";
  import Empty from "../ui/empty/empty.svelte";
  import EmptyHeader from "../ui/empty/empty-header.svelte";
  import EmptyTitle from "../ui/empty/empty-title.svelte";
  import EmptyDescription from "../ui/empty/empty-description.svelte";

  let {
    tournament,
    entrants,
  }: {
    tournament: tournament.Tournament;
    entrants: entrant.Entrant[];
  } = $props();

  const queryClient = useQueryClient();

  const updateTournamentMutation = createMutation(() => ({
    mutationFn: async ({ qualifyingOrder }: { qualifyingOrder: string[] }) =>
      UpdateTournament(tournament.id, { qualifyingOrder }),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["tournament", tournament.id],
      });
    },
  }));

  function moveToFirst(idx: number) {
    const [entrant] = orderedEntrants.splice(idx, 1);

    if (!entrant) return;

    const reordered = [entrant, ...orderedEntrants];
    updateTournamentMutation.mutate({
      qualifyingOrder: reordered.map((entrant) => entrant.participantId),
    });
  }

  const updateQualifierMutation = createMutation(() => ({
    mutationFn: async ({
      qualifierId,
      status,
    }: {
      qualifierId: string;
      status: entrant.Entrant["status"];
    }) =>
      UpdateEntrant({
        TournamentID: tournament.id,
        EntrantID: qualifierId,
        Status: status,
      }),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["entrants", tournament.id],
      });
    },
  }));

  function markAs(id: string, status: entrant.Entrant["status"]) {
    updateQualifierMutation.mutate({
      qualifierId: id,
      status,
    });
  }

  let concludeDialogOpen = $state(false);
  const concludeQualifiersMutation = createMutation(() => ({
    mutationFn: async () => ConcludeQualifers(tournament.id),
    onSuccess: async () => {
      await UpdateTournamentState(tournament.id, "start");
      await SyncTournament(tournament.id);

      await queryClient.invalidateQueries({
        queryKey: ["matches", tournament.id],
      });
      await queryClient.invalidateQueries({
        queryKey: ["activeBattles", tournament.id],
      });

      await queryClient.invalidateQueries({
        queryKey: ["tournament", tournament.id],
      });

      concludeDialogOpen = false;
    },
    onError: (error) => {
      console.error("Failed to conclude qualifiers:", error);
    },
  }));

  const orderedEntrants = $derived.by(() => {
    const order = tournament.qualifying_order.split(",").map((id) => id.trim());
    const byParticipantId = new Map(
      entrants.map((entrant) => [entrant.participantId, entrant]),
    );
    return order
      .map((id) => byParticipantId.get(id))
      .filter((entrant): entrant is entrant.Entrant => entrant !== undefined);
  });

  const front = $derived(orderedEntrants[0]?.participantId);

  const concluding = $derived.by(() => {
    return !entrants.some((entrant) => entrant.status === "pending");
  });
</script>

<div class="space-4">
  <div class="flex items-end justify-between gap-4 pb-4">
    <div>
      <h3 class="font-semibold tracking-tight">Entrants</h3>
      <p class="text-sm text-muted-foreground">Current qualification scores</p>
    </div>
    {#if orderedEntrants.length > 0}
      <span class="text-xs font-medium text-muted-foreground">
        {orderedEntrants.length}
        {orderedEntrants.length === 1 ? "entrant" : "entrants"}
      </span>
    {/if}
  </div>

  {#if orderedEntrants.length === 0}
    <Empty class="min-h-40 border-dashed py-10">
      <EmptyHeader>
        <EmptyTitle>No entrants yet</EmptyTitle>
        <EmptyDescription>
          Entrants will appear here once they are added to this tournament.
        </EmptyDescription>
      </EmptyHeader>
    </Empty>
  {:else}
    <div class="divide-y overflow-hidden rounded-lg border bg-card">
      {#each orderedEntrants as entrant, idx (entrant.id)}
        <button
          type="button"
          class={[
            "flex w-full items-center justify-between gap-4 px-3 py-3.5 text-left transition-colors sm:px-4",
            front === entrant.participantId
              ? "bg-accent text-accent-foreground"
              : "hover:bg-muted/50",
          ]}
        >
          <div class="flex min-w-0 items-center gap-3">
            <span
              class="flex size-7 shrink-0 items-center justify-center rounded-md border bg-muted/50 text-xs font-semibold tabular-nums text-muted-foreground"
              aria-label={`Rank ${idx + 1}`}
            >
              {idx + 1}
            </span>

            <div class="min-w-0 space-y-1">
              <p class="truncate font-medium">{entrant.name}</p>

              {#if entrant.status === "pending"}
                {#if idx === 0}
                  <Badge variant="default">On Going</Badge>
                {:else}
                  <Badge variant="outline">Pending</Badge>
                {/if}
              {:else if entrant.status === "qualified"}
                <Badge>Qualified</Badge>
              {:else if entrant.status === "not_qualified"}
                <Badge variant="destructive">Not Qualified</Badge>
              {:else if entrant.status === "scored"}
                <Badge variant="secondary">Scored</Badge>
              {/if}
            </div>
          </div>

          <div class="flex shrink-0 items-center gap-2 sm:gap-3">
            <span
              class="min-w-10 text-right text-lg font-semibold tabular-nums"
            >
              {entrant.status === "pending" ? "TBA" : entrant.score}
            </span>

            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div onclick={(event) => event.stopPropagation()}>
              <DropdownMenu>
                <DropdownMenuTrigger>
                  <Button
                    variant="ghost"
                    size="icon"
                    aria-label="Qualifier actions"
                  >
                    <EllipsisVertical />
                  </Button>
                </DropdownMenuTrigger>

                <DropdownMenuContent align="end" class="w-48">
                  <DropdownMenuItem
                    onclick={() => moveToFirst(idx)}
                    disabled={entrant.status === "not_qualified" ||
                      updateTournamentMutation.isPending ||
                      idx === 0}
                  >
                    Move to front
                  </DropdownMenuItem>

                  <DropdownMenuItem
                    onclick={() => markAs(entrant.id, "not_qualified")}
                    disabled={updateQualifierMutation.isPending}
                  >
                    Set as not qualified
                  </DropdownMenuItem>

                  <DropdownMenuItem
                    onclick={() => markAs(entrant.id, "pending")}
                    disabled={updateQualifierMutation.isPending}
                  >
                    Set as pending
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </button>
      {/each}
    </div>
  {/if}

  {#if concluding}
    <AlertDialog bind:open={concludeDialogOpen}>
      <AlertDialogTrigger class="w-full pt-4">
        <Button
          class="w-full"
          variant="default"
          disabled={concludeQualifiersMutation.isPending}
        >
          Conclude
        </Button>
      </AlertDialogTrigger>

      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Conclude qualifying?</AlertDialogTitle>

          <AlertDialogDescription>
            After concluding qualifying, the qualified entrants will be
            finalized and will be seeded to the tournament bracket. This action
            cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>

          <AlertDialogAction
            onclick={() => {
              concludeQualifiersMutation.mutate();
            }}
            disabled={concludeQualifiersMutation.isPending}
          >
            {concludeQualifiersMutation.isPending
              ? "Concluding..."
              : "Conclude"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  {/if}
</div>
