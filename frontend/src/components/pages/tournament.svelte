<script lang="ts">
  // TOOD (important): LOADING STATE FOR tournament, matches, activeBattles, entrants
  import { useRoute } from "@dvcol/svelte-simple-router/router";
  import { createQuery } from "@tanstack/svelte-query";
  import DefaultCriteria from "$src/data/criteria/default-criteria.json";

  import Main from "../layout/main.svelte";
  import {
    GetTournament,
    GetEntrants,
    GetMatches,
    GetActiveMatches,
  } from "$wails/go/main/App";
  import EntrantList from "../judging/entrant-list.svelte";
  import StageCard from "../judging/stage-card.svelte";
  import Panel from "../judging/panel.svelte";
  import { Stage } from "$src/lib/stage.svelte";
  import { Judgement } from "$src/lib/judging.svelte";
  import QualifyingHandler from "../judging/qualifying-handler.svelte";
  import Matches from "../tournaments/matches.svelte";
  import BattleHandler from "../judging/battle-handler.svelte";

  const route = $derived(useRoute());
  const location = $derived(route.location);
  const tournamentID = $derived(location?.params?.id as string);

  let stage = new Stage();
  let judgement = new Judgement();

  const criteriaConfig = $derived(DefaultCriteria as Judging.CriteriaConfig);
  const categories = $derived(Object.entries(criteriaConfig.categories));

  const tournamentQuery = createQuery(() => ({
    queryKey: ["tournament", tournamentID],
    enabled: !!tournamentID,
    queryFn: async () => GetTournament(tournamentID),
  }));

  const matchesQuery = createQuery(() => ({
    queryKey: ["matches", tournamentID],
    enabled: !!tournamentID,
    queryFn: async (): Promise<Tournament.GroupedMatches> =>
      GetMatches(tournamentID),
  }));

  const activeBattles = createQuery(() => ({
    queryKey: ["activeBattles", tournamentID],
    enabled: !!tournamentID,
    queryFn: async () => GetActiveMatches(tournamentID),
  }));

  const entrantsQuery = createQuery(() => ({
    queryKey: ["entrants", tournamentID],
    enabled: !!tournamentID,
    queryFn: () => GetEntrants(tournamentID),
  }));
</script>

<Main>
  {#if tournamentQuery.data}
    {@const tournament = tournamentQuery.data}
    {@const isQualifying = tournament.status === "qualifying"}
    <div class="mx-auto max-w-6xl space-y-4 p-6">
      <StageCard>
        {#if isQualifying}
          <EntrantList
            tournament={tournamentQuery.data ?? {}}
            entrants={entrantsQuery.data ?? []}
          />
        {:else}
          <Matches
            {tournamentID}
            {stage}
            matches={matchesQuery.data ?? []}
            activeBattles={activeBattles.data ?? []}
            {judgement}
          />
        {/if}
      </StageCard>

      <Panel {stage} {judgement} {categories} {isQualifying}>
        {#if isQualifying}
          <QualifyingHandler
            {tournament}
            {judgement}
            entrants={entrantsQuery.data ?? []}
          />
        {:else}
          <BattleHandler {tournamentID} {judgement} {stage} />
        {/if}
      </Panel>
    </div>
  {/if}
</Main>
