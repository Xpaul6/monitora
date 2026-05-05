<script lang="ts">
  import { getAllServers } from "../lib/api";
  import type { GetAllServersResponse, Server } from "../lib/models";
  import ServerDashboard from "./ServerDashboard.svelte";
  import ServerForm from "./ServerForm.svelte";

  let serverList = $state<Server[]>([]);
  let panelState = $state<"main" | "form" | "dashboard">("main");
  let currentServer = $state<Server | null>(null);

  async function handleGetAllServers() {
    try {
      const res: GetAllServersResponse | string = await getAllServers();
      if (typeof res != "string") {
        serverList = res.servers;
      } else {
        alert("Unable to get server list: " + res);
      }
    } catch (e) {
      alert("Network error: " + e);
    }
  }

  function setCurrentServer(s: Server) {
    currentServer = s;
    panelState = "dashboard";
  }

  function handleLogout() {
    localStorage.setItem("token", "");
    window.location.href = "/";
  }

  $effect(() => {
    if (panelState == "main") {
      handleGetAllServers().then();
    }
  });
</script>

<!-- Main component section -->
<main class="w-full flex flex-col items-center">
  {#if panelState == "main"}
    <button class="form-button" onclick={handleLogout}>Logout</button>
    <h1>Admin panel</h1>
    <div
      class="m-2 p-2 border border-gray-400 rounded-md flex flex-col items-center w-[75%]"
    >
      <h2>Servers: {serverList.length}</h2>
      {#each serverList as s (s.ID)}
        {@render server(s)}
      {/each}
      <button
        class="form-button"
        onclick={() => {
          panelState = "form";
        }}>+</button
      >
    </div>
  {:else if panelState == "form"}
    <button class="form-button" onclick={() => (panelState = "main")}
      >Go back</button
    >
    <ServerForm bind:panelState />
  {:else if panelState == "dashboard"}
    <button class="form-button" onclick={() => (panelState = "main")}
      >Go back</button
    >
    <ServerDashboard server={currentServer} bind:panelState />
  {/if}
</main>

<!-- Snippets -->
{#snippet server(s: Server)}
  <div
    class="flex flex-row justify-between w-full p-1 m-1 border border-transparent
      hover:border hover:border-gray-500 rounded-md transition-all duration-150"
  >
    <div class="flex items-center">{s.name}: {s.ip}</div>
    <div
      class={`flex items-center ${s.status == "Online" ? "text-green-500" : "text-red-500"}`}
    >
      {s.status || "Waiting"}
    </div>
    <button class="form-button" onclick={() => setCurrentServer(s)}
      >Observe</button
    >
  </div>
{/snippet}
