<script lang="ts">
  import CriteriaConfigDialog from "$src/components/criteria-config/criteria-dialog.svelte";

  import {
    GetCriteriaConfigs,
    DeleteCriteriaConfig,
    CopyCriteriaConfigById,
  } from "$wails/go/main/App";
  import { repository } from "$wails/go/models";
  import {
    createMutation,
    createQuery,
    useQueryClient,
  } from "@tanstack/svelte-query";

  import * as AlertDialog from "$src/components/ui/alert-dialog";
  import * as Table from "$src/components/ui/table";
  import * as Card from "$src/components/ui/card";
  import { Button } from "$src/components/ui/button";
  import { Badge } from "$src/components/ui/badge";
  import { Separator } from "$src/components/ui/separator";
  import EllipsisVertical from "@lucide/svelte/icons/ellipsis-vertical";
  import * as DropdownMenu from "$src/components/ui/dropdown-menu";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";

  const queryClient = useQueryClient();

  const configsQuery = createQuery(() => ({
    queryKey: ["criteria-configs"],
    queryFn: () => GetCriteriaConfigs(),
  }));

  let configDialogOpen = $state(false);
  let dialogMode = $state<"create" | "update">("create");
  let selectedConfig = $state<repository.CriteriaConfig | undefined>(undefined);

  function openUpdateDialog(config: repository.CriteriaConfig) {
    selectedConfig = config;
    dialogMode = "update";
    configDialogOpen = true;
  }

  const deleteMutation = createMutation(() => ({
    mutationFn: (id: number) => DeleteCriteriaConfig(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["criteria-configs"],
      });
    },
  }));

  const copyMutation = createMutation(() => ({
    mutationFn: (id: number) => CopyCriteriaConfigById(id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: ["criteria-configs"],
      });
    },
  }));

  let configToDelete = $state<repository.CriteriaConfig | undefined>(undefined);

  let deleteDialogOpen = $state(false);

  function deleteConfig(config: repository.CriteriaConfig) {
    configToDelete = config;
    deleteDialogOpen = true;
  }

  function confirmDelete() {
    if (!configToDelete) return;

    deleteMutation.mutate(configToDelete.id, {
      onSuccess: () => {
        deleteDialogOpen = false;
        configToDelete = undefined;
      },
    });
  }

  let mounted = $state(false);
  onMount(() => {
    mounted = true;
  });
</script>

<AlertDialog.Root bind:open={deleteDialogOpen}>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>Delete criteria config?</AlertDialog.Title>
      <AlertDialog.Description>
        Are you sure you want to delete
        <strong>{configToDelete?.name}</strong>? This action cannot be undone.
      </AlertDialog.Description>
    </AlertDialog.Header>

    <AlertDialog.Footer>
      <AlertDialog.Cancel
        disabled={deleteMutation.isPending}
        onclick={() => {
          configToDelete = undefined;
        }}
      >
        Cancel
      </AlertDialog.Cancel>

      <AlertDialog.Action
        onclick={(event) => {
          event.preventDefault();
          confirmDelete();
        }}
        disabled={deleteMutation.isPending}
        class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
      >
        {deleteMutation.isPending ? "Deleting..." : "Delete"}
      </AlertDialog.Action>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>

<CriteriaConfigDialog
  bind:open={configDialogOpen}
  mode={dialogMode}
  config={selectedConfig}
/>

{#if mounted}
  <div in:fade={{ duration: 300 }}>
    <div class="space-y-6">
      <div class="flex justify-between items-end">
        <div>
          <h1 class="text-2xl font-semibold">Criteria Configs</h1>
          <p class="text-muted-foreground mt-1 text-sm">
            Manage the criteria used for judging.
          </p>
        </div>
        <Button
          onclick={() => {
            configDialogOpen = true;
            selectedConfig = undefined;
          }}>+ New Config</Button
        >
      </div>

      <Separator />

      <div>
        <h2 class="text-lg font-semibold">Existing Configs</h2>
        <p class="text-muted-foreground mt-1 text-sm">
          {configsQuery.data?.length ?? 0} configurations
        </p>
      </div>

      {#if deleteMutation.isError}
        <p class="text-destructive text-sm">
          Failed to delete config: {String(deleteMutation.error)}
        </p>
      {/if}

      {#if configsQuery.isPending}
        <p class="text-muted-foreground text-sm">Loading criteria configs...</p>
      {:else if configsQuery.isError}
        <p class="text-destructive text-sm">
          Failed to load configs: {String(configsQuery.error)}
        </p>
      {:else}
        <Card.Root>
          <Card.Content>
            <Table.Root>
              <Table.Header>
                <Table.Row>
                  <Table.Head class="w-50">Name</Table.Head>
                  <Table.Head>Categories</Table.Head>
                  <Table.Head class="w-50 text-right">Actions</Table.Head>
                </Table.Row>
              </Table.Header>

              <Table.Body>
                {#each configsQuery.data ?? [] as config (config.id)}
                  <Table.Row>
                    <Table.Cell class="font-medium">
                      {config.name}
                    </Table.Cell>

                    <Table.Cell>
                      <div class="flex flex-wrap gap-2">
                        {#each [config.categories.common, config.categories.lead, config.categories.chase] as category}
                          <Badge variant="secondary">
                            {category.label} · {Math.round(
                              category.tags.reduce(
                                (total, tag) =>
                                  total +
                                  (tag.type !== "deduction"
                                    ? (tag.max ?? 0)
                                    : 0),
                                0,
                              ),
                            )}
                          </Badge>
                        {/each}
                      </div>
                    </Table.Cell>

                    <Table.Cell class="text-right">
                      <DropdownMenu.Root>
                        <DropdownMenu.Trigger>
                          {#snippet child({ props })}
                            <Button
                              {...props}
                              variant="ghost"
                              size="icon"
                              aria-label="Open actions menu"
                            >
                              <EllipsisVertical class="size-4" />
                            </Button>
                          {/snippet}
                        </DropdownMenu.Trigger>

                        <DropdownMenu.Content align="end">
                          <DropdownMenu.Item
                            onclick={() => openUpdateDialog(config)}
                          >
                            Edit
                          </DropdownMenu.Item>

                          <DropdownMenu.Item
                            onclick={() => copyMutation.mutate(config.id)}
                          >
                            Copy
                          </DropdownMenu.Item>

                          <DropdownMenu.Separator />

                          <DropdownMenu.Item
                            variant="destructive"
                            onclick={() => deleteConfig(config)}
                            disabled={deleteMutation.isPending}
                          >
                            Delete
                          </DropdownMenu.Item>
                        </DropdownMenu.Content>
                      </DropdownMenu.Root>
                    </Table.Cell>
                  </Table.Row>
                {:else}
                  <Table.Row>
                    <Table.Cell
                      colspan={3}
                      class="text-muted-foreground h-24 text-center"
                    >
                      No criteria configs yet.
                    </Table.Cell>
                  </Table.Row>
                {/each}
              </Table.Body>
            </Table.Root>
          </Card.Content>
        </Card.Root>
      {/if}
    </div>
  </div>
{/if}
