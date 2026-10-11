<script lang="ts">
  import CardContent from "../ui/card/card-content.svelte";
  import CardHeader from "../ui/card/card-header.svelte";
  import Card from "../ui/card/card.svelte";
  import Skeleton from "../ui/skeleton/skeleton.svelte";
</script>

<Card class="py-0">
  <CardHeader class="sr-only">Loading tournament</CardHeader>

  <CardContent class="p-0">
    <div class="grid lg:grid-cols-[minmax(0,1fr)_18rem]">
      <div class="space-y-3 p-3 sm:p-4">
        <div class="flex items-center gap-3">
          <Skeleton class="size-9 shrink-0 rounded-md" />
          <div class="flex flex-1 items-center justify-between gap-3">
            <Skeleton class="h-5 w-48" />
            <Skeleton class="h-5 w-24" />
          </div>
        </div>
        <div class="flex flex-wrap gap-6 pl-12">
          <Skeleton class="h-4 w-28" />
          <Skeleton class="h-4 w-36" />
          <Skeleton class="h-4 w-28" />
        </div>
      </div>
      <div
        class="flex flex-col gap-3 border-t border-border bg-muted/20 p-3 sm:p-4 lg:border-t-0 lg:border-l"
      >
        <div class="space-y-2">
          <div class="flex justify-between">
            <Skeleton class="h-3 w-14" />
            <Skeleton class="h-3 w-16" />
          </div>
          <Skeleton class="h-1.5 w-full" />
        </div>
        <div class="grid-cols[minmax(0,1fr)_auto] grid gap-2">
          <Skeleton class="h-9 w-full" />
          <Skeleton class="h-9 w-20" />
        </div>
      </div>
    </div></CardContent
  >
</Card>
