<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
	let status: string = $state("offline");

	async function ping() {
    try {
      const response = await fetch("/api/ping");
      if (!response.ok) throw new Error("server is offline");

      const data = await response.json();
      status = data.status;
    } catch (e) {
      console.log(e)
      status = "offline"
    }
	}

	ping()
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div>
  <h1>Backend check</h1>
  <span>{status}</span>
  {@render children()}
</div>
