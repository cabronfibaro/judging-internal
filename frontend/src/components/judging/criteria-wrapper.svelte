<script lang="ts">
  import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
  } from "$src/components/ui/card";
  import type { Judgement } from "$src/lib/judging.svelte";
  import Separator from "../ui/separator/separator.svelte";
  import CriterionRow from "./criterion-row.svelte";

  let {
    category,
    role,
    judgement,
  }: {
    category: Judging.CriteriaCategory;
    role: "lead" | "chase";
    judgement: Judgement;
  } = $props();
</script>

<Card class="gap-0 border-border/70 shadow-none">
  <CardHeader class="flex-row items-baseline justify-between gap-3 pb-4">
    <CardTitle class="min-w-0 truncate">{category.label}</CardTitle>
    <span class="shrink-0 text-xs font-medium text-muted-foreground capitalize"
      >{role}</span
    >
  </CardHeader>
  <Separator />
  <CardContent class="space-y-1 pt-2">
    {#each category.tags as tag, index (tag.id)}
      <CriterionRow
        {tag}
        value={judgement.get(tag.id)?.[role]}
        showSeparator={index > 0}
        onValueChange={(value) => judgement.set(tag.id, role, value)}
      />
    {/each}
  </CardContent>
</Card>
