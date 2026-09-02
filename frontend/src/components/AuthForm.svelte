<script lang="ts">
  import { type AuthRequest, type LoginResponse } from "../lib/models";
  import { login, register } from "../lib/api";

  let email = $state<string>("");
  let password = $state<string>("");

  async function handleRegister() {
    let body: AuthRequest = { email: email, password: password };
    try {
      const errorMessage: string | null = await register(body);
      if (!errorMessage) {
        alert("Registered, sign in");
        password = "";
        email = "";
      } else {
        alert("Registration failed: " + errorMessage);
      }
    } catch (e) {
      alert("Network error: " + e);
    }
  }

  async function handleLogin() {
    let body: AuthRequest = { email: email, password: password };
    try {
      const res: LoginResponse | string = await login(body);
      if (typeof res != "string") {
        localStorage.setItem("token", res.token);
        window.location.href = "/";
      } else {
        alert("Login failed: " + res);
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
  <label for="email">Login</label>
  <input type="text" id="email" class="form-input" bind:value={email} />
  <label for="password">Password</label>
  <input
    type="password"
    id="password"
    class="form-input"
    bind:value={password}
  />
  <div class="flex flex-row justify-around mt-3">
    <button class="form-button" onclick={handleRegister}>Sign up</button>
    <button class="form-button" onclick={handleLogin}>Sign in</button>
  </div>
</form>
