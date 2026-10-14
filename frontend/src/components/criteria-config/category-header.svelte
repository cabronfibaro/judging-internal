<script lang="ts">
  import { Button } from "$src/components/ui/button";

  type CategoryKey = "common" | "lead" | "chase";

  type Props = {
    title: string;
    categoryKey: CategoryKey;
    tagCount: number;
    onAddTag: (category: CategoryKey) => void;
    onCopyCategory: (
      source: "lead" | "chase",
      target: "lead" | "chase",
    ) => void;
  };

  let { title, categoryKey, tagCount, onAddTag, onCopyCategory }: Props =
    $props();
</script>

<div class="flex items-center justify-between gap-4">
  <div class="space-y-1">
    <h3 class="font-semibold">{title} Criteria</h3>
    <p class="text-muted-foreground text-sm">{tagCount} tags</p>
  </div>

  <div class="flex items-center gap-2">
    {#if categoryKey === "lead"}
      <Button
        type="button"
        variant="ghost"
        size="sm"
        onclick={() => onCopyCategory("chase", "lead")}
      >
        Copy from Chase
      </Button>
    {:else if categoryKey === "chase"}
      <Button
        type="button"
        variant="ghost"
        size="sm"
        onclick={() => onCopyCategory("lead", "chase")}
      >
        Copy from Lead
      </Button>
    {/if}

    <Button
      type="button"
      variant="outline"
      size="sm"
      onclick={() => onAddTag(categoryKey)}
    >
      Add Tag
    </Button>
  </div>
</div>
