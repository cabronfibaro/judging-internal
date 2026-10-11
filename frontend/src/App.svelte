<script lang="ts">
  import "./app.css";
  import {
    RouterView,
    type Route,
    type RouterOptions,
  } from "@dvcol/svelte-simple-router";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import Home from "./components/layouts/home.svelte";
  import Tournament from "./components/layouts/tournament.svelte";

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

<QueryClientProvider client={queryClient}>
  <RouterView {options} />
</QueryClientProvider>
