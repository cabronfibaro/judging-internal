<script lang="ts">
  import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
  } from "$src/components/ui/select";
  import Slider from "../ui/slider/slider.svelte";
  import Switch from "../ui/switch/switch.svelte";

  let {
    tag,
    value = 0,
    onValueChange,
  }: {
    tag: Judging.Tag;
    value?: unknown;
    onValueChange: (value: unknown) => void;
  } = $props();
</script>

{#if tag.value_type === "enum"}
  <Select
    type="single"
    onValueChange={(selected) => onValueChange(Number(selected))}
  >
    <SelectTrigger class="w-full sm:w-36" aria-label={tag.label}>
      <SelectValue placeholder="Select" />
    </SelectTrigger>

    <SelectContent>
      {#each tag.options ?? [] as option (option.value)}
        <SelectItem value={String(option.value)}>
          {option.label}
        </SelectItem>
      {/each}
    </SelectContent>
  </Select>
{:else if tag.value_type === "dynamic"}
  <div class="flex items-center justify-end gap-3 sm:justify-center">
    <Slider
      type="single"
      value={Math.abs(Number(value ?? 0)) / (tag.multiplier ?? 1)}
      min={0}
      max={tag.max !== undefined ? tag.max / (tag.multiplier ?? 1) : 100}
      step={1}
      aria-label={tag.label}
      onValueChange={(value) => {
        const multiplier = tag.multiplier ?? 1;
        const amount = value * multiplier;

        onValueChange(tag.type === "deduction" ? -amount : amount);
      }}
      class="w-full sm:w-36"
    />

    <!-- <span class="w-8 shrink-0 text-right text-sm tabular-nums">
			{Math.abs(Number(value ?? 0)) / (tag.multiplier ?? 1)}
		</span> -->
  </div>
{:else if tag.value_type === "switch"}
  <Switch
    aria-label={tag.label}
    checked={Number(value ?? 0) !== 0}
    onCheckedChange={(checked) => {
      const amount = Number(tag.default_value ?? 0);

      if (checked) {
        onValueChange(Number(value ?? 0) + amount);
      } else {
        onValueChange(Number(value ?? 0) - amount);
      }
    }}
  />
{/if}
