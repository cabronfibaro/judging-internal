<script lang="ts">
  import type { Judgement } from "$src/lib/judging.svelte";
  import type { Snippet } from "svelte";
  import type { Stage } from "$src/lib/stage.svelte";

  import Empty from "$src/components/ui/empty/empty.svelte";
  import EmptyHeader from "$src/components/ui/empty/empty-header.svelte";
  import EmptyMedia from "$src/components/ui/empty/empty-media.svelte";
  import EmptyTitle from "$src/components/ui/empty/empty-title.svelte";
  import EmptyDescription from "$src/components/ui/empty/empty-description.svelte";
  import Spinner from "$src/components/ui/spinner/spinner.svelte";
  import CommonCriteria from "./common-criteria.svelte";
  import CriteriaWrapper from "./criteria-wrapper.svelte";
  import type { tournament, repository } from "$wails/go/models";

  let {
    tournament,
    judgement,
    isQualifying,
    stage,
    children,
  }: {
    tournament: tournament.Tournament;
    judgement: Judgement;
    isQualifying: boolean;
    stage: Stage;
    children: Snippet<[]>;
  } = $props();

  const criteriaConfig = $derived(
    JSON.parse(tournament.criteriaConfig) as repository.CriteriaCategories,
  );

  const categories = $derived(Object.entries(criteriaConfig));

  const common = $derived(categories.find(([id]) => id === "common"));
  const lead = $derived(categories.find(([id]) => id === "lead"));
  const chase = $derived(categories.find(([id]) => id === "chase"));
</script>

<div class="relative flex flex-col gap-5" aria-busy={judgement.processing}>
  {#if isQualifying || stage.matchId}
    <div class="grid items-start gap-5 xl:grid-cols-[minmax(0,1fr)_18rem]">
      <div class="grid min-w-0 grid-cols-1 gap-3 md:grid-cols-2 md:gap-4">
        {#if common}
          {@const [, category] = common}
          <div class={!isQualifying ? "md:col-span-2" : ""}>
            <CommonCriteria {category} {judgement} {isQualifying} />
          </div>
        {/if}

        {#if lead}
          {@const [, category] = lead}
          <CriteriaWrapper {category} role="lead" {judgement} />
        {/if}

        {#if chase && !isQualifying}
          {@const [, category] = chase}
          <CriteriaWrapper {category} role="chase" {judgement} />
        {/if}
      </div>

      <div class="min-w-0 xl:sticky xl:top-24">{@render children()}</div>
    </div>
  {/if}

  {#if judgement.processing}
    <div
      class="absolute inset-0 z-50 flex items-center justify-center rounded-xl bg-background/95 px-4"
    >
      <Empty class="w-full max-w-md border-0 bg-transparent shadow-none ring-0">
        <EmptyHeader>
          <EmptyMedia variant="icon">
            <Spinner />
          </EmptyMedia>

          <EmptyTitle>Processing the scores</EmptyTitle>
          <EmptyDescription>
            Please wait while we process your request. Do not refresh the page.
          </EmptyDescription>
        </EmptyHeader>
      </Empty>
    </div>
  {/if}
</div>
