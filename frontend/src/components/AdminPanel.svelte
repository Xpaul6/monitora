<script lang="ts">
  import { onMount } from "svelte";
  import { getAllServers } from "../lib/api";
  import { type GetAllServersResponse, type Server } from "../lib/models";

  let serverList = $state<Server[]>([]);
  let currentServer = $state<Server>({} as Server);

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

  onMount(handleGetAllServers);
</script>

<main class="w-full flex flex-col items-center">
  <h1>Admin panel</h1>
  <div
    class="m-2 p-2 border border-gray-400 rounded-md flex flex-col items-center w-[75%]"
  >
    <h2>Servers: {serverList.length}</h2>
    {#each serverList as s}
      {@render server(s)}
    {/each}
    <button class="form-button">+</button>
  </div>
</main>

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
    <button class="form-button">Observe</button>
  </div>
{/snippet}
