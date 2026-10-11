<script lang="ts">
  import CalendarDays from "@lucide/svelte/icons/calendar-days";
  import ExternalLink from "@lucide/svelte/icons/external-link";
  import Flag from "@lucide/svelte/icons/flag";
  import GitBranch from "@lucide/svelte/icons/git-branch";
  import ArrowUpRight from "@lucide/svelte/icons/arrow-up-right";
  import Trash2 from "@lucide/svelte/icons/trash-2";
  import { Badge } from "$src/components/ui/badge";
  import Button from "../ui/button/button.svelte";
  import Card from "../ui/card/card.svelte";
  import CardContent from "../ui/card/card-content.svelte";
  import Progress from "../ui/progress/progress.svelte";

  let {
    tournament,
    openDeleteDialog,
    position,
  }: {
    tournament: any;
    openDeleteDialog(tournament: any): void;
    position: number;
  } = $props();
</script>

<Card
  class="group w-full overflow-hidden border-border/80 py-0 transition-colors hover:border-foreground/20"
>
  <CardContent class="p-0">
    <div class="grid w-full lg:grid-cols-[minmax(0,1fr)_18rem]">
      <div class="min-w-0 space-y-3 p-3 sm:p-4 lg:pr-6">
        <div class="flex min-w-0 items-center gap-3">
          <div
            class="flex size-9 shrink-0 items-center justify-center rounded-md border bg-muted/40 text-sm font-semibold text-muted-foreground tabular-nums"
            aria-label={`Tournament ${position}`}
          >
            {String(position).padStart(2, "0")}
          </div>

          <div
            class="flex min-w-0 flex-1 flex-wrap items-center gap-x-3 gap-y-1"
          >
            <h2
              class="min-w-0 flex-1 truncate text-base font-semibold tracking-tight sm:text-lg"
            >
              {tournament.name}
            </h2>
            <Badge
              variant={tournament.status === "qualifying"
                ? "secondary"
                : "default"}
            >
              {tournament.status === "qualifying"
                ? "Qualifying"
                : "Elimination"}
            </Badge>
          </div>
        </div>

        <div
          class="flex flex-wrap items-center gap-x-6 gap-y-2 pl-12 text-sm text-muted-foreground"
        >
          <div class="flex min-w-0 items-center gap-2">
            <Flag class="size-4 shrink-0" aria-hidden="true" />
            <span
              class="truncate"
              title={tournament.track ?? "Track unavailable"}
            >
              {tournament.track ?? "Track unavailable"}
            </span>
          </div>
          <div class="flex min-w-0 items-center gap-2">
            <CalendarDays class="size-4 shrink-0" aria-hidden="true" />
            <span class="truncate">{tournament.startingDate}</span>
          </div>
          <a
            href={`https://challonge.com/${tournament.url}`}
            target="_blank"
            rel="noopener noreferrer"
            class="flex min-w-0 items-center gap-2 text-foreground underline-offset-4 transition-colors hover:text-primary hover:underline"
            aria-label={`Open Challonge bracket for ${tournament.name}`}
          >
            <GitBranch class="size-4 shrink-0" aria-hidden="true" />
            <span>Open bracket</span>
            <ExternalLink class="size-3.5 shrink-0" aria-hidden="true" />
          </a>
        </div>
      </div>

      <div
        class="flex min-w-0 flex-col gap-3 border-t border-border bg-muted/20 p-3 sm:p-4 lg:border-t-0 lg:border-l"
      >
        <div class="min-w-0 space-y-2">
          <div class="flex items-center justify-between gap-3 text-xs">
            <span class="font-medium text-muted-foreground">Progress</span>
            <span class="text-foreground tabular-nums"
              >{tournament.progress * 100}% complete</span
            >
          </div>
          <Progress
            value={tournament.progress * 100}
            aria-label={`${tournament.name}: ${tournament.progress * 100}% complete`}
          />
        </div>

        <div class="grid grid-cols-[minmax(0,1fr)_auto] gap-2">
          <Button href={`/tournaments/${tournament.id}`} class="w-full">
            Open tournament
            <ArrowUpRight aria-hidden="true" />
          </Button>
          <Button
            variant="ghost"
            class="text-destructive hover:bg-destructive/10 hover:text-destructive"
            onclick={() => openDeleteDialog(tournament)}
          >
            <Trash2 aria-hidden="true" />
            Delete
          </Button>
        </div>
      </div>
    </div>
  </CardContent>
</Card>
