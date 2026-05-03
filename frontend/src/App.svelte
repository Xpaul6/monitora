<script lang="ts">
  import { onMount } from "svelte";

  import AuthForm from "./components/AuthForm.svelte";
  import { checkAuth } from "./lib/api";

  let status = $state<string>("offline");
  let isLoggedIn = $state<boolean>(false);

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

  async function getToken() {
    let token = localStorage.getItem("token");
    if (token == "" || token == null) {
      isLoggedIn = false;
      return;
    }
    try {
      isLoggedIn = await checkAuth();
    } catch (e) {
      alert(e);
      isLoggedIn = false;
    }
  }

  onMount(ping);
  onMount(getToken);
</script>

<div class="flex flex-col items-center w-full">
  <h1>monitora</h1>
  <h2>Backend status: {status}</h2>
  <h2>Auth status: {isLoggedIn}</h2>
  {#if !isLoggedIn}
    <hr />
    <AuthForm />
  {/if}
</div>
