<script lang="ts">
  import { SetChallongeAPIKey, HasChallongeAPIKey } from "$wails/go/main/App";
  import { Button } from "$src/components/ui/button";
  import { Input } from "$src/components/ui/input";
  import { Label } from "$src/components/ui/label";
  import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
  } from "$src/components/ui/dialog";

  let { open = $bindable() }: { open: boolean } = $props();

  let apiKey = $state("");
  let submitting = $state(false);
  let error = $state("");

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    const key = apiKey.trim();
    if (!key || submitting) return;

    submitting = true;
    error = "";

    try {
      await SetChallongeAPIKey(key);
      const hasAPIKey = await HasChallongeAPIKey();

      if (!hasAPIKey) {
        error = "The API key could not be verified. Please try again.";
        return;
      }

      open = false;
      apiKey = "";
    } catch (err) {
      console.error("Failed to configure Challonge API key:", err);
      error = "Failed to save or verify the API key. Please try again.";
    } finally {
      submitting = false;
    }
  }
</script>

<Dialog bind:open>
  <DialogContent
    class="sm:max-w-md"
    onInteractOutside={(event) => event.preventDefault()}
    showCloseButton={false}
  >
    <DialogHeader>
      <DialogTitle>Configure Challonge API Key</DialogTitle>
      <DialogDescription>
        A Challonge API key is required to manage tournament brackets and
        synchronize participants. Enter your API key to continue.
      </DialogDescription>
    </DialogHeader>

    <form onsubmit={handleSubmit} class="space-y-4">
      <div class="space-y-2">
        <Label for="challonge-api-key">API Key</Label>
        <Input
          id="challonge-api-key"
          type="password"
          bind:value={apiKey}
          placeholder="Enter your Challonge API key"
          autocomplete="off"
          required
          disabled={submitting}
        />
      </div>

      {#if error}
        <p class="text-sm text-destructive" role="alert">
          {error}
        </p>
      {/if}

      <Button
        type="submit"
        class="w-full"
        disabled={submitting || !apiKey.trim()}
      >
        {submitting ? "Verifying..." : "Save API Key"}
      </Button>
    </form>
  </DialogContent>
</Dialog>
