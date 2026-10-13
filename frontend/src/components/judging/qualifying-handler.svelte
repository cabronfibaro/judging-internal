<script lang="ts">
  import { createMutation, useQueryClient } from "@tanstack/svelte-query";
  import { UpdateTournament, UpdateEntrant } from "$wails/go/main/App";
  import type { Judgement } from "$src/lib/judging.svelte";
  import Button from "../ui/button/button.svelte";
  import CardContent from "../ui/card/card-content.svelte";
  import CardHeader from "../ui/card/card-header.svelte";
  import CardTitle from "../ui/card/card-title.svelte";
  import Card from "../ui/card/card.svelte";
  import type { entrant, tournament } from "$wails/go/models";

  let {
    tournament = $bindable(),
    judgement,
    entrants,
  }: {
    tournament: tournament.Tournament;
    judgement: Judgement;
    entrants: entrant.Entrant[];
  } = $props();

  const queryClient = useQueryClient();

  let submitting = $state(false);
  const currentScore = $derived(judgement.finalizedScore());

  const updateTournamentMutation = createMutation(() => ({
    mutationFn: async ({ qualifyingOrder }: { qualifyingOrder: string[] }) =>
      UpdateTournament(tournament.id, { qualifyingOrder }),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["tournament", tournament.id],
      });
    },
  }));

  const updateEntrantMutation = createMutation(() => ({
    mutationFn: async ({
      entrantId,
      score,
      status,
    }: {
      entrantId: string;
      score: number;
      status: entrant.Entrant["status"];
    }) =>
      UpdateEntrant({
        TournamentID: tournament.id,
        EntrantID: entrantId,
        Score: score,
        Status: status,
      }),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["entrants", tournament.id],
      });
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

  async function onSubmitJudgement() {
    if (submitting) return;

    judgement.processing = true;
    submitting = true;

    try {
      const { lead } = judgement.finalizedScore();

      await handleQualifying(lead);

      judgement.clear();
    } finally {
      submitting = false;
      judgement.processing = false;
    }
  }

  async function handleQualifying(score: number) {
    const front = orderedEntrants[0];

    await updateEntrantMutation.mutateAsync({
      entrantId: front.id,
      score,
      status: "scored",
    });

    const order = tournament.qualifying_order.split(",").map((id) => id.trim());

    const first = order.shift();
    if (first) {
      order.push(first);
    }

    await updateTournamentMutation.mutateAsync({
      qualifyingOrder: order,
    });
  }
</script>

<Card class="gap-0 border-border/70 shadow-none">
  <CardHeader class="p-4">
    <CardTitle>Judgement</CardTitle>
  </CardHeader>

  <CardContent class="space-y-3 pt-0">
    <div
      class="flex items-center justify-between gap-4 rounded-md border bg-muted/30 px-3 py-2.5"
    >
      <span class="text-sm font-medium text-muted-foreground"
        >Current score</span
      >
      <strong class="text-sm tabular-nums">{currentScore.lead}</strong>
    </div>

    <Button
      class="h-11 w-full"
      onclick={onSubmitJudgement}
      disabled={submitting}
    >
      {#if submitting}
        Submitting...
      {:else}
        Submit Judgment
      {/if}
    </Button>

    <Button
      variant="outline"
      type="reset"
      class="w-full"
      onclick={() => judgement.clear()}
      disabled={submitting}
    >
      Reset
    </Button>
  </CardContent>
</Card>
