<script lang="ts">
  import { onMount, tick } from "svelte";
  import { Chart, registerables } from "chart.js";
  import "chartjs-adapter-date-fns";
  Chart.register(...registerables);

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
  let periodBegin = $state<Date>(new Date());
  let periodEnd = $state<Date>(new Date());
  let stats = $state<GetStatsByPeriodResponse[]>([]);
  let chartInstances = $state<Chart[]>([]);

  // Group stats by unit (e.g., "%", "celsius", "bytes")
  let groupedStats = $derived.by(() => {
    const groups: Array<{
      unit: string;
      unitStats: GetStatsByPeriodResponse[];
    }> = [];
    const unitMap: Record<string, GetStatsByPeriodResponse[]> = {};

    for (const stat of stats) {
      const unit = stat.metric_type.unit;
      if (!unitMap[unit]) {
        unitMap[unit] = [];
      }
      unitMap[unit].push(stat);
    }

    for (const [unit, unitStats] of Object.entries(unitMap)) {
      groups.push({ unit, unitStats });
    }

    return groups;
  });

  function getRandomColor() {
    const r = Math.floor(Math.random() * 256);
    const g = Math.floor(Math.random() * 256);
    const b = Math.floor(Math.random() * 256);
    return `rgb(${r},${g},${b})`;
  }

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
        alert("Unable to fetch stats: " + res);
      }
    } catch (e) {
      alert("Network error: " + e);
    }

    await tick();
    drawCharts();
  }

  const BYTES_TO_GB = Math.pow(1024, 3);

  function drawCharts() {
    chartInstances.forEach((chart) => chart.destroy());
    chartInstances = [];

    for (const { unit, unitStats } of groupedStats) {
      // Skip string-type metrics
      if (unit === "string") continue;

      const canvas = document.getElementById(
        `chart-${unit}`,
      ) as HTMLCanvasElement;
      if (!canvas) continue;

      const datasets: any[] = [];

      const componentGroups: Record<string, GetStatsByPeriodResponse[]> = {};
      for (const stat of unitStats) {
        const componentKey = `${stat.component.type}:${stat.component.address || stat.component.type}`;
        if (!componentGroups[componentKey]) {
          componentGroups[componentKey] = [];
        }
        componentGroups[componentKey].push(stat);
      }

      for (const [_, componentStats] of Object.entries(componentGroups)) {
        const component = componentStats[0].component;
        const metricNames = [
          ...new Set(componentStats.map((s) => s.metric_type.name)),
        ];

        for (const metricName of metricNames) {
          const metricData = componentStats
            .filter((s) => s.metric_type.name === metricName)
            .map((m) => ({
              x: new Date(m.timestamp),
              y: unit === "bytes" ? m.value / BYTES_TO_GB : m.value,
            }));

          const label =
            component.type === "cpu" || component.type === "mem"
              ? `${metricName}`
              : `${component.address} - ${metricName}`;

          datasets.push({
            label,
            data: metricData,
            borderColor: getRandomColor(),
            tension: 0.1,
          });
        }
      }

      const chart = new Chart(canvas, {
        type: "line",
        data: { datasets },
        options: {
          responsive: true,
          scales: {
            x: {
              type: "time",
              time: { unit: "minute" },
              title: { display: true, text: "Time" },
            },
            y: {
              title: {
                display: true,
                text: unit === "bytes" ? "GB" : unit,
              },
            },
          },
        },
      });

      chartInstances.push(chart);
    }
  }

  onMount(handleGetServerComponents);
  onMount(handleGetMetricTypes);
</script>

<div class="flex flex-col items-center">
  <h1>Dashboard</h1>
  <!-- Server info -->
  <div class="flex flex-row gap-5 justify-around">
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
      <div class="flex items-center">
        <button class="form-button" onclick={handleGetStatsByPeriod}
          >Get metrics</button
        >
      </div>
    </form>
  </div>
  {#if stats.length == 0}
    <p>No metrics in selected period</p>
  {:else}
    <p>Points: {stats.length / metricTypes.length}</p>
  {/if}
  <div class="flex flex-row flex-wrap justify-around">
    {#each groupedStats as group (group.unit)}
      <div
        class="chart-container m-2 p-2 border border-gray-300 rounded-md min-w-100"
      >
        <h3 class="text-center">{group.unit} metrics</h3>
        <canvas id="chart-{group.unit}" class="min-h-75"></canvas>
      </div>
    {/each}
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
