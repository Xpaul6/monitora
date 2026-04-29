<script lang="ts">
    import { onMount } from "svelte";

  let status = $state<string>("offline")

  async function ping() {
    try {
      const res = await fetch("/api/ping");
      if (!res.ok) throw new Error("lol no");

      status = "online";
    } catch (e) {
      console.log(e);
      status = "offline";
    }
  }

  onMount(ping)
</script>

<div>
  <h1>monitora</h1>
  <h3>Backend status: { status }</h3>
</div>
