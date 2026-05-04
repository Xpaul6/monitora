<script lang="ts">
  import { onMount } from "svelte";
  import {
    getMetricTypes,
    getServerComponents,
    getStatsByPeriod,
  } from "../lib/api";
  import type {
    MetricType,
    Component,
    GetServerComponentsRequest,
    Server,
    GetStatsByPeriodRequest,
    GetStatsByPeriodResponse,
  } from "../lib/models";

  let { server = $bindable<Server>() } = $props();
  let serverComponents = $state<Component[]>([]);
  let metricTypes = $state<MetricType[]>([]);
  let periodBegin = $state<Date>({} as Date);
  let periodEnd = $state<Date>({} as Date);
  let stats = $state<GetStatsByPeriodResponse[]>([]);

  let serverComponentMap = $derived.by<Map<number, Component>>(() => {
    let ret = new Map<number, Component>();
    for (let c of serverComponents) {
      ret.set(c.ID, c);
    }
    return ret;
  });

  let metricTypeMap = $derived.by<Map<number, MetricType>>(() => {
    let ret = new Map<number, MetricType>();
    for (let m of metricTypes) {
      ret.set(m.ID, m);
    }
    return ret;
  });

  async function handleGetServerComponents() {
    const body: GetServerComponentsRequest = { id: server.ID };
    try {
      const res = await getServerComponents(body);
      if (typeof res != "string") {
        serverComponents = res;
      } else {
        alert("Unable to fetch components: " + res);
      }
    } catch (e) {
      alert("Network error: " + e);
    }
  }

  async function handleGetMetricTypes() {
    try {
      const res = await getMetricTypes();
      if (typeof res != "string") {
        metricTypes = res;
      } else {
        alert("Unable to fetch metric types: " + res);
      }
    } catch (e) {
      alert("Network error: " + e);
    }
  }

  async function handleGetStatsByPeriod() {
    const body: GetStatsByPeriodRequest = {
      server_id: server.ID,
      period_begin: new Date(periodBegin).toISOString(),
      period_end: new Date(periodEnd).toISOString(),
    };
    try {
      const res = await getStatsByPeriod(body);
      if (typeof res != "string") {
        stats = res;
      } else {
        alert("Unable to fetch metric types: " + res);
      }
    } catch (e) {
      alert("Network error: " + e);
    }
  }

  onMount(handleGetServerComponents);
  onMount(handleGetMetricTypes);
</script>

<div>
  <h1>Dashboard</h1>
  <!-- Server info -->
  <div class="flex flex-row gap-5">
    <div>
      <h2>Server:</h2>
      <p>ID: {server.ID}</p>
      <p>Name: {server.name}</p>
      <p>IP: {server.ip}</p>
      <p>Status: {server.status}</p>
    </div>
    <div>
      <h2>Components:</h2>
      {#each serverComponents as c (c.ID)}
        {@render component(c)}
      {/each}
    </div>
  </div>
  <!-- Stats -->
  <div>
    <form
      onsubmit={(e) => e.preventDefault()}
      class="flex flex-row border border-s rounded-md m-5 p-5 gap-1"
    >
      <label for="begin" class="flex items-center">Period start</label>
      <input
        type="date"
        id="begin"
        class="form-input"
        bind:value={periodBegin}
      />
      <label for="end" class="flex items-center">Period end</label>
      <input type="date" id="end" class="form-input" bind:value={periodEnd} />
      <div class="flex flex-row justify-around mt-3">
        <button class="form-button" onclick={handleGetStatsByPeriod}
          >Get metrics</button
        >
      </div>
    </form>
  </div>
</div>

{#snippet component(c: Component)}
  <div>
    {#if c.type == "cpu" || c.type == "mem"}
      <p>{c.address}</p>
    {:else}
      <p>{c.type}: {c.address}</p>
    {/if}
  </div>
{/snippet}
