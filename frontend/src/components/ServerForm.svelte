<script lang="ts">
  import { addServer } from "../lib/api";
  import type { AddServerRequest } from "../lib/models";

  let { panelState = $bindable() } = $props();

  let name = $state<string>("");
  let ip = $state<string>("");

  async function handleAddServer() {
    const body: AddServerRequest = { name: name, ip: ip };
    try {
      const res = await addServer(body);
      if (typeof res != "string") {
        panelState = 'main';
      } else {
        alert("Failed to add server: " + res);
        panelState = 'main';
      }
    } catch (e) {
      alert("Network error: " + e);
    }
  }
</script>

<form
  onsubmit={(e) => e.preventDefault()}
  class="flex flex-col border border-s rounded-md m-5 p-5 gap-1"
>
  <label for="name">Name</label>
  <input type="text" id="name" class="form-input" bind:value={name} />
  <label for="ip">IP:port</label>
  <input type="text" id="ip" class="form-input" bind:value={ip} />
  <div class="flex flex-row justify-around mt-3">
    <button class="form-button" onclick={handleAddServer}>Add server</button>
  </div>
</form>
