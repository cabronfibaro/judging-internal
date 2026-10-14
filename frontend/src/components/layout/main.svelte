<script lang="ts">
  import type { Snippet } from "svelte";
  import { link } from "svelte-spa-router";
  import { cn } from "$src/utils";
  import { Button } from "$src/components/ui/button";
  import Separator from "$src/components/ui/separator/separator.svelte";
  import ThemeToggle from "./theme-toggle.svelte";
  import Menu from "@lucide/svelte/icons/menu";
  import X from "@lucide/svelte/icons/x";
  import LayoutDashboard from "@lucide/svelte/icons/layout-dashboard";
  import Trophy from "@lucide/svelte/icons/trophy";
  import UsersRound from "@lucide/svelte/icons/users-round";
  import MapPinned from "@lucide/svelte/icons/map-pinned";
  import SliderHorizontal from "@lucide/svelte/icons/sliders-horizontal";
  import ChevronRight from "@lucide/svelte/icons/chevron-right";

  let {
    fixed = false,
    fluid = false,
    class: className = "",
    children,
  }: {
    fixed?: boolean;
    fluid?: boolean;
    class?: string;
    children?: Snippet;
  } = $props();

  let mobileNavOpen = $state(false);

  // const onTournamentRoute = $derived(
  //   page.url.pathname.startsWith("/tournaments"),
  // );
  // const tournamentsActive = $derived(
  //   page.url.pathname === "/" || onTournamentRoute,
  // );

  function closeMobileNav() {
    mobileNavOpen = false;
  }
</script>

<div class="min-h-svh bg-muted/20 text-foreground">
  <aside
    class="fixed inset-y-0 left-0 z-40 hidden w-64 flex-col border-r border-sidebar-border bg-sidebar lg:flex"
    aria-label="Primary navigation"
  >
    <div
      class="flex h-16 items-center gap-3 border-b border-sidebar-border px-5"
    >
      <div
        class="flex size-9 items-center justify-center rounded-md bg-sidebar-primary text-sidebar-primary-foreground"
      >
        <LayoutDashboard class="size-4" aria-hidden="true" />
      </div>
      <div class="min-w-0">
        <p class="truncate text-sm font-semibold text-sidebar-foreground">
          Judging Console
        </p>
        <p class="truncate text-xs text-sidebar-foreground/60">
          Tournament operations
        </p>
      </div>
    </div>

    <nav class="flex-1 space-y-1 p-3" aria-label="Main">
      <p
        class="px-3 pt-2 pb-2 text-[11px] font-semibold tracking-[0.12em] text-sidebar-foreground/50 uppercase"
      >
        Workspace
      </p>

      <a
        href="/"
        // aria-current={tournamentsActive ? "page" : undefined}
        class={cn(
          "group flex h-10 items-center gap-3 rounded-md px-3 text-sm font-medium transition-colors",
          // tournamentsActive
          //   ? "bg-sidebar-accent text-sidebar-accent-foreground"
          //   : "text-sidebar-foreground/70 hover:bg-sidebar-accent/70 hover:text-sidebar-foreground",
        )}
      >
        <Trophy class="size-4" aria-hidden="true" />
        <span class="flex-1">Tournaments</span>
        <!-- {#if tournamentsActive}
          <ChevronRight
            class="size-3.5 text-sidebar-foreground/50"
            aria-hidden="true"
          />
        {/if} -->
      </a>

      <button
        type="button"
        disabled
        class="flex h-10 w-full cursor-not-allowed items-center gap-3 rounded-md px-3 text-left text-sm font-medium text-sidebar-foreground/45"
        aria-label="Drivers, coming soon"
      >
        <UsersRound class="size-4" aria-hidden="true" />
        <span class="flex-1">Drivers</span>
        <span
          class="text-[10px] font-medium tracking-wide text-sidebar-foreground/45 uppercase"
          >Soon</span
        >
      </button>

      <button
        type="button"
        disabled
        class="flex h-10 w-full cursor-not-allowed items-center gap-3 rounded-md px-3 text-left text-sm font-medium text-sidebar-foreground/45"
        aria-label="Tracks, coming soon"
      >
        <MapPinned class="size-4" aria-hidden="true" />
        <span class="flex-1">Tracks</span>
        <span
          class="text-[10px] font-medium tracking-wide text-sidebar-foreground/45 uppercase"
          >Soon</span
        >
      </button>
      <a
        href="/criteria-configs"
        use:link
        class={cn(
          "group flex h-10 items-center gap-3 rounded-md px-3 text-sm font-medium transition-colors",
          // tournamentsActive
          //   ? "bg-sidebar-accent text-sidebar-accent-foreground"
          //   : "text-sidebar-foreground/70 hover:bg-sidebar-accent/70 hover:text-sidebar-foreground",
        )}
      >
        <SliderHorizontal class="size-4" aria-hidden="true" />
        <span class="flex-1">Criteria Configs</span>
      </a>
    </nav>

    <Separator class="bg-sidebar-border" />
    <div class="flex items-center justify-between gap-3 p-3">
      <div class="min-w-0">
        <p class="text-xs font-medium text-sidebar-foreground">Appearance</p>
        <p class="truncate text-xs text-sidebar-foreground/60">
          Switch between light and dark
        </p>
      </div>
      <ThemeToggle />
    </div>
  </aside>

  {#if mobileNavOpen}
    <button
      type="button"
      class="fixed inset-0 z-40 bg-background/80 lg:hidden"
      aria-label="Close navigation"
      onclick={closeMobileNav}
    ></button>
  {/if}

  <div class="flex min-h-svh flex-col lg:pl-64">
    <header
      class="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-border bg-background/95 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/80 lg:px-8"
    >
      <div class="flex items-center gap-3">
        <Button
          variant="ghost"
          size="icon"
          class="lg:hidden"
          onclick={() => (mobileNavOpen = !mobileNavOpen)}
          aria-label={mobileNavOpen ? "Close navigation" : "Open navigation"}
        >
          {#if mobileNavOpen}
            <X aria-hidden="true" />
          {:else}
            <Menu aria-hidden="true" />
          {/if}
        </Button>
        <div class="lg:hidden">
          <p class="text-sm font-semibold">Judging Console</p>
          <p class="text-xs text-muted-foreground">Tournament operations</p>
        </div>
        <div
          class="hidden items-center gap-2 text-sm text-muted-foreground lg:flex"
        >
          <LayoutDashboard class="size-4" aria-hidden="true" />
          <!-- <span>{onTournamentRoute ? "Tournament workspace" : "Overview"}</span> -->
        </div>
      </div>

      <div class="flex items-center gap-1 lg:hidden">
        <ThemeToggle />
      </div>
    </header>

    {#if mobileNavOpen}
      <aside
        class="fixed inset-y-0 left-0 z-50 flex w-72 flex-col border-r border-sidebar-border bg-sidebar shadow-xl lg:hidden"
        aria-label="Primary navigation"
      >
        <div
          class="flex h-16 items-center justify-between border-b border-sidebar-border px-5"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex size-9 items-center justify-center rounded-md bg-sidebar-primary text-sidebar-primary-foreground"
            >
              <LayoutDashboard class="size-4" aria-hidden="true" />
            </div>
            <p class="text-sm font-semibold text-sidebar-foreground">
              Judging Console
            </p>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onclick={closeMobileNav}
            aria-label="Close navigation"
          >
            <X aria-hidden="true" />
          </Button>
        </div>
        <nav class="space-y-1 p-3" aria-label="Main">
          <a
            href="/"
            onclick={closeMobileNav}
            class="flex h-10 items-center gap-3 rounded-md bg-sidebar-accent px-3 text-sm font-medium text-sidebar-accent-foreground"
          >
            <Trophy class="size-4" aria-hidden="true" />
            Tournaments
          </a>
          <div
            class="flex h-10 items-center gap-3 rounded-md px-3 text-sm font-medium text-sidebar-foreground/45"
          >
            <UsersRound class="size-4" aria-hidden="true" />
            <span class="flex-1">Drivers</span><span
              class="text-[10px] uppercase">Soon</span
            >
          </div>
          <div
            class="flex h-10 items-center gap-3 rounded-md px-3 text-sm font-medium text-sidebar-foreground/45"
          >
            <MapPinned class="size-4" aria-hidden="true" />
            <span class="flex-1">Tracks</span><span
              class="text-[10px] uppercase">Soon</span
            >
          </div>
        </nav>
        <div
          class="mt flex items-center justify-between gap-3 border-t border-sidebar-border p-3"
        >
          <div class="min-w-0">
            <p class="text-xs font-medium text-sidebar-foreground">
              Appearance
            </p>
            <p class="truncate text-xs text-sidebar-foreground/60">
              Switch between light and dark
            </p>
          </div>
          <ThemeToggle />
        </div>
      </aside>
    {/if}

    <main
      data-layout={fixed ? "fixed" : ""}
      class={cn(
        "flex-1 px-4 py-6 sm:px-6 lg:px-8 lg:py-8",
        fixed && "flex grow flex-col overflow-hidden",
        !fluid && "w-full max-w-none",
        className,
      )}
    >
      {@render children?.()}
    </main>
  </div>
</div>
