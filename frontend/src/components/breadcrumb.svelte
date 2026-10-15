<script lang="ts">
  import * as Breadcrumb from "$src/components/ui/breadcrumb";
  import { link } from "svelte-spa-router";

  type BreadcrumbItem = {
    label: string;
    href: string;
  };

  let { items }: { items: BreadcrumbItem[] } = $props();
</script>

<Breadcrumb.Root>
  <Breadcrumb.List>
    {#each items as item, index (item.href)}
      {#if index > 0}
        <Breadcrumb.Separator />
      {/if}

      <Breadcrumb.Item>
        {#if index === items.length - 1}
          <Breadcrumb.Page>{item.label}</Breadcrumb.Page>
        {:else}
          <a
            href={item.href}
            use:link
            class="text-muted-foreground transition-colors hover:text-foreground"
          >
            {item.label}
          </a>
        {/if}
      </Breadcrumb.Item>
    {/each}
  </Breadcrumb.List>
</Breadcrumb.Root>
