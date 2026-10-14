<script lang="ts">
  /* eslint-disable svelte/require-each-key */
  import { Badge } from "$src/components/ui/badge";
  import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
  } from "$src/components/ui/card";
  import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
  } from "$src/components/ui/tabs";
  import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
  } from "$src/components/ui/alert-dialog";
  import { capitalize } from "es-toolkit/string";
  import { Stage } from "$src/lib/stage.svelte";
  import {
    ResetBattle,
    ConcludeBattle,
    SyncTournament,
  } from "$wails/go/main/App";
  import { tournament, type repository } from "$wails/go/models";
  import type { Judgement } from "$src/lib/judging.svelte";
  import Button from "../ui/button/button.svelte";
  import { useQueryClient } from "@tanstack/svelte-query";
  import Empty from "../ui/empty/empty.svelte";
  import EmptyHeader from "../ui/empty/empty-header.svelte";
  import EmptyTitle from "../ui/empty/empty-title.svelte";
  import EmptyDescription from "../ui/empty/empty-description.svelte";

  type Battles = repository.ActiveBattle[];

  let {
    tournamentID,
    matches,
    judgement,
    stage,
    activeBattles,
  }: {
    tournamentID: string;
    matches: Tournament.GroupedMatches;
    judgement: Judgement;
    stage: Stage;
    activeBattles: Battles;
  } = $props();

  const queryClient = useQueryClient();

  let selectedMatchId = $state<string | undefined>(undefined);

  function getLead(
    p1Battle: Battles[number],
    p2Battle: Battles[number],
  ): "p1" | "p2" {
    if (p1Battle.leadStatus === -1) return "p2";
    if (p2Battle.leadStatus === -1) return "p1";

    return p1Battle.leadStatus === 1 ? "p1" : "p2";
  }

  function onStartMatch(match: repository.Match[]) {
    stage.clear();
    judgement.clear();

    selectedMatchId = match[0].matchId;
    stage.matchId = selectedMatchId;

    const p1 = match[0];
    const p2 = match[1];

    if (!p1.participantId || !p2.participantId) {
      console.warn("missing participant id");
      return;
    }
    const battles = activeBattles.filter(
      (battle) => battle.matchId === p1.matchId,
    );

    if (battles.length !== 2) {
      console.warn("active battle not found or incomplete");
      return;
    }

    const p1Battle = battles.find(
      (battle) => battle.participantId === p1.participantId,
    );
    const p2Battle = battles.find(
      (battle) => battle.participantId === p2.participantId,
    );

    if (!p1Battle || !p2Battle) {
      console.warn("missing active battle participant");
      return;
    }

    const lead = getLead(p1Battle, p2Battle);

    if (!lead) {
      console.warn("no lead participant");
      return;
    }

    stage.setDrivers(p1.participantId, p2.participantId, lead);
  }

  const sortedMatches = $derived(
    [...matches]
      .sort((a, b) => {
        const roundA = a[0]?.[0]?.round ?? 0;
        const roundB = b[0]?.[0]?.round ?? 0;

        if (roundA === 0) return 1;
        if (roundB === 0) return -1;

        return roundA - roundB;
      })
      .map((round) =>
        round.map((match) => {
          const battles = activeBattles.filter(
            (battle) => battle.matchId === match[0]?.matchId,
          );
          const leadingBattle = battles.find(
            (battle) => battle.leadStatus === 1,
          );
          return {
            match,
            leadingParticipantId: leadingBattle?.participantId ?? null,
          };
        }),
      ),
  );

  async function onHandleOMT(match: repository.Match[]) {
    await ResetBattle(tournamentID, selectedMatchId!);

    await queryClient.invalidateQueries({
      queryKey: ["activeBattles", tournamentID],
    });

    await queryClient.invalidateQueries({
      queryKey: ["matches", tournamentID],
    });

    onStartMatch(match);
  }

  async function onHandleConclude() {
    await ConcludeBattle(tournamentID, selectedMatchId!);
    await SyncTournament(tournamentID);

    await queryClient.invalidateQueries({
      queryKey: ["activeBattles", tournamentID],
    });

    await queryClient.invalidateQueries({
      queryKey: ["matches", tournamentID],
    });

    stage.clear();
    selectedMatchId = undefined;
  }

  let isProcessing = $state(false);
  let omtDialogOpen = $state(false);
  let concludeDialogOpen = $state(false);

  async function handleOMT(match: repository.Match[]) {
    isProcessing = true;

    try {
      await onHandleOMT(match);
      omtDialogOpen = false;
    } finally {
      isProcessing = false;
    }
  }

  async function handleConclude() {
    isProcessing = true;

    try {
      await onHandleConclude();
      concludeDialogOpen = false;
    } finally {
      isProcessing = false;
    }
  }

  function roundKeys() {
    let keys: string[] = [];
    for (let i = 0; i < sortedMatches.length; i++) {
      const matches = sortedMatches[i].length;
      if (matches > 2) {
        keys.push(`Round ${i + 1}`);
      } else {
        if (matches == 2) {
          keys.push("Semifinal");
        } else if (matches == 1) {
          if (sortedMatches[i][0].match[0].round === 0) {
            keys.push("Third Place");
            continue;
          }
          keys.push("Final");
        }
      }
    }
    return keys;
  }
</script>

<div class="space-y-4 pt-4">
  <div>
    <h3 class="font-semibold">Bracket</h3>
    <p class="text-sm text-muted-foreground">Drift battle matches by round</p>
  </div>

  <Tabs value="round-1">
    {#if sortedMatches.length === 0}
      <Empty class="min-h-48 border-dashed py-10">
        <EmptyHeader>
          <EmptyTitle>No bracket matches yet</EmptyTitle>
          <EmptyDescription>
            Matches will appear here after qualifying is concluded and the
            bracket is synced.
          </EmptyDescription>
        </EmptyHeader>
      </Empty>
    {:else}
      <TabsList
        class="overflow-x w-full justify-start gap-1 overscroll-x-contain"
      >
        {#each sortedMatches as round, i}
          {@const roundNumber = round[0]?.match[0]?.round ?? 0}

          {#if roundNumber === 0}
            <TabsTrigger class="min-w-fit flex-none" value="third-place"
              >3rd Place</TabsTrigger
            >
          {:else}
            <TabsTrigger
              class="min-w-fit flex-none"
              value={`round-${roundNumber}`}
            >
              {roundKeys()[i]}
            </TabsTrigger>
          {/if}
        {/each}
      </TabsList>
    {/if}

    {#each sortedMatches as round, i}
      {@const roundNumber = round[0]?.match[0]?.round ?? 0}
      {@const tabValue =
        roundNumber === 0 ? "third-place" : `round-${roundNumber}`}

      <TabsContent value={tabValue} class="mt-4 space-y-3">
        <div>
          <h4 class="font-medium">
            {roundKeys()[i]}
          </h4>

          <p class="text-sm text-muted-foreground">
            {round.length}
            {round.length === 1 ? "battle" : "battles"}
          </p>
        </div>

        <div class="grid gap-3">
          {#each round as matchData}
            {@const match = matchData.match}
            {@const currentMatch = match[0]}
            {@const hasBothParticipants =
              match[0].participantId && match[1].participantId}

            <Card class="border-border/80">
              <CardHeader class="border-b pb-3">
                <div class="flex items-center justify-between">
                  <CardTitle class="text-sm"
                    >Match {round.indexOf(matchData) + 1}</CardTitle
                  >

                  <div class="flex items-center gap-2">
                    {#if currentMatch.status !== "open"}
                      <Badge variant="secondary">
                        {capitalize(currentMatch.status)}
                      </Badge>
                    {/if}
                  </div>
                </div>
              </CardHeader>

              <CardContent class="space-y-3 pt-4">
                {#each match as entrant}
                  <div
                    class={[
                      "flex items-center justify-between gap-3 rounded-md border px-3 py-3 sm:px-4",
                      entrant.winnerId &&
                      entrant.winnerId === entrant.challongeParticipantId
                        ? "border-green-500/30 bg-green-500/10"
                        : "",
                    ]}
                  >
                    <div class="flex min-w-0 items-center gap-3">
                      <span class="truncate font-medium">
                        {entrant.participant_name ?? "TBD"}
                      </span>

                      {#if hasBothParticipants && matchData.leadingParticipantId === entrant.participantId && currentMatch.status !== "complete"}
                        <Badge variant="outline">Lead</Badge>
                      {/if}
                    </div>

                    <span class="shrink-0 text-lg font-semibold tabular-nums"
                      >{entrant.score}</span
                    >

                    {#if entrant.winnerId !== null && entrant.winnerId === entrant.participantId}
                      <Badge>Winner</Badge>
                    {/if}
                  </div>
                {/each}

                <div class="flex justify-end gap-4">
                  {#if hasBothParticipants && selectedMatchId !== currentMatch.matchId && currentMatch.status !== "complete"}
                    <Button onclick={() => onStartMatch(match)}>Load</Button>
                  {/if}

                  {#if currentMatch.bothRunFinished}
                    {#if selectedMatchId === currentMatch.matchId}
                      <AlertDialog bind:open={omtDialogOpen}>
                        <AlertDialogTrigger>
                          <Button variant="outline" disabled={isProcessing}
                            >OMT</Button
                          >
                        </AlertDialogTrigger>

                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Start OMT?</AlertDialogTitle>
                            <AlertDialogDescription>
                              This will reset the battle and start an OMT. Are
                              you sure you want to continue?
                            </AlertDialogDescription>
                          </AlertDialogHeader>

                          <AlertDialogFooter>
                            <AlertDialogCancel disabled={isProcessing}
                              >Cancel</AlertDialogCancel
                            >

                            <AlertDialogAction
                              disabled={isProcessing}
                              onclick={() => handleOMT(match)}
                            >
                              {#if isProcessing}
                                Processing...
                              {:else}
                                Continue
                              {/if}
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>

                      <AlertDialog bind:open={concludeDialogOpen}>
                        <AlertDialogTrigger>
                          <Button variant="outline" disabled={isProcessing}
                            >Conclude</Button
                          >
                        </AlertDialogTrigger>

                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Conclude match?</AlertDialogTitle>
                            <AlertDialogDescription>
                              This will finalize the result of this match. Are
                              you sure?
                            </AlertDialogDescription>
                          </AlertDialogHeader>

                          <AlertDialogFooter>
                            <AlertDialogCancel disabled={isProcessing}
                              >Cancel</AlertDialogCancel
                            >

                            <AlertDialogAction
                              disabled={isProcessing}
                              onclick={() => handleConclude()}
                            >
                              {#if isProcessing}
                                Processing...
                              {:else}
                                Conclude
                              {/if}
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
                    {/if}
                  {/if}
                </div>
              </CardContent>
            </Card>
          {/each}
        </div>
      </TabsContent>
    {/each}
  </Tabs>
</div>
