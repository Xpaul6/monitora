<script lang="ts">
    import { derived } from "svelte/store";
  import type { Component, MetricType } from "../lib/models";

  interface Props {
    newLimitComponent: Component,
    newLimitMetric: MetricType,
    newLimitValue: number,
    components: Component[],
    metrics: MetricType[]
  };

  let {
    newLimitComponent = $bindable(),
    newLimitMetric = $bindable(),
    newLimitValue = $bindable(),
    components,
    metrics
  }: Props = $props();
  let filteredMetrics = $derived(metrics.filter(m => m.name.startsWith(newLimitComponent.type)));
</script>

<form
  onsubmit={(e) => e.preventDefault()}
  class="flex flex-col border border-s rounded-md m-5 p-5 gap-1"
>
  <label for="component">Component</label>
  <select name="component" id="component" class="form-input" bind:value={newLimitComponent}>
    {#each components as c (c.ID)}
      <option value={c}>{c.address}</option>
    {/each}
  </select>
  <label for="metric">Metric type</label>
  <select name="metric" id="metric" class="form-input" bind:value={newLimitMetric}>
    {#each filteredMetrics as m (m.ID)}
      <option value={m}>{m.name} ({m.unit})</option>
    {/each}
  </select>
  <label for="value">Value</label>
  <input
    type="number"
    id="value"
    class="form-input"
    bind:value={newLimitValue}
  />
</form>
