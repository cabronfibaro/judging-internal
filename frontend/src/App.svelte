<script lang="ts">
  import {
    RouterView,
    type Route,
    type RouterOptions,
  } from "@dvcol/svelte-simple-router";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { ModeWatcher } from "mode-watcher";

  import Home from "./components/pages/home.svelte";
  import Tournament from "./components/pages/tournament.svelte";

  const RouteName = {
    Home: "home",
    Tournaments: "Tournaments",
  } as const;

  type RouteNames = (typeof RouteName)[keyof typeof RouteName];
  const routes: Readonly<Route<RouteNames>[]> = [
    {
      name: RouteName.Home,
      path: "/",
      component: Home,
    },
    {
      name: RouteName.Tournaments,
      path: "/tournaments/:id",
      component: Tournament,
    },
  ] as const;

  const options: RouterOptions<RouteNames> = {
    routes,
  } as const;

  const queryClient = new QueryClient();
</script>

<ModeWatcher defaultMode="dark" />
<QueryClientProvider client={queryClient}>
  <RouterView {options} />
</QueryClientProvider>
