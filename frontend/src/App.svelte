<script lang="ts">
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { ModeWatcher } from "mode-watcher";
  import Router from "svelte-spa-router";
  import { wrap } from "svelte-spa-router/wrap";

  import Home from "./components/pages/home.svelte";

  const routes = {
    "/": Home,
    "/tournaments/:id": wrap({
      asyncComponent: () => import("./components/pages/tournament.svelte"),
    }),
  };

  const queryClient = new QueryClient();
</script>

<ModeWatcher defaultMode="dark" />
<QueryClientProvider client={queryClient}>
  <SvelteQueryDevtools />
  <Router {routes} />
</QueryClientProvider>
