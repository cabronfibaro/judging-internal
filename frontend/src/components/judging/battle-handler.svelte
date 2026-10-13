<script lang="ts">
  import {
    UpdateMatchSingleParticipant,
    ChangeLeader,
  } from "$wails/go/main/App";
  import type { Judgement } from "$src/lib/judging.svelte";
  import Button from "../ui/button/button.svelte";
  import CardContent from "../ui/card/card-content.svelte";
  import CardHeader from "../ui/card/card-header.svelte";
  import CardTitle from "../ui/card/card-title.svelte";
  import Card from "../ui/card/card.svelte";
  import type { Stage } from "$src/lib/stage.svelte";
  import { useQueryClient } from "@tanstack/svelte-query";

  let {
    judgement,
    stage,
  }: {
    judgement: Judgement;
    stage: Stage;
  } = $props();

  const queryClient = useQueryClient();
  let submitting = $state(false);
  const currentScore = $derived(judgement.finalizedScore());

  async function onSubmitJudgement() {
    if (submitting) return;

    judgement.processing = true;
    submitting = true;

    try {
      const scores = judgement.finalizedScore();
      stage.setScore(scores);
      const [p1, p2] = stage.getScores();

      await Promise.all([
        UpdateMatchSingleParticipant({
          TournamentID: stage.tournamentId,
          MatchID: stage.matchId!,
          ParticipantID: p1.participantId,
          AppendScore: true,
          Score: p1.score,
        }),
        UpdateMatchSingleParticipant({
          TournamentID: stage.tournamentId,
          MatchID: stage.matchId!,
          ParticipantID: p2.participantId,
          AppendScore: true,
          Score: p2.score,
        }),
      ]);

      await ChangeLeader(stage.tournamentId, stage.matchId!);

      await queryClient.invalidateQueries({
        queryKey: ["activeBattles", stage.tournamentId],
      });

      await queryClient.invalidateQueries({
        queryKey: ["matches", stage.tournamentId],
      });

      stage.switchLeader();
      judgement.clear();
    } finally {
      submitting = false;
      judgement.processing = false;
    }
  }
</script>

<Card class="gap-0 border-border/70 shadow-none">
  <CardHeader class="pb-4">
    <CardTitle>Judgement</CardTitle>
  </CardHeader>

  <CardContent class="space-y-3 pt-0">
    <div
      class="flex items-center justify-between gap-4 rounded-md border bg-muted/30 px-3 py-2.5"
    >
      <span class="text-sm font-medium text-muted-foreground"
        >Current score</span
      >
      <div class="flex items-center gap-3 text-sm tabular-nums">
        <span>Lead <strong>{currentScore.lead}</strong></span>
        <span class="text-muted-foreground">/</span>
        <span>Chase <strong>{currentScore.chase}</strong></span>
      </div>
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
