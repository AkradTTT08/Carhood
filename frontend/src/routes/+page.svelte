<script lang="ts">
  import { connect } from "$lib/socketStore";
  import { goto } from "$app/navigation";

  let roomID = "";
  let username = "";

  function joinGame() {
    if (roomID && username) {
      connect(roomID, username);
      goto("/game");
    }
  }
</script>

<main class="container">
  <div class="glass-card main-card">
    <h1 class="logo">CARHOOD</h1>

    <div class="form-group">
      <input
        type="text"
        placeholder="GAME PIN"
        bind:value={roomID}
        class="input-field pin-input"
      />
      <input
        type="text"
        placeholder="USERNAME"
        bind:value={username}
        class="input-field"
      />
      <button class="btn-primary" on:click={joinGame}>Enter</button>
    </div>
  </div>
</main>

<style lang="scss">
  .container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    padding: 1rem;
  }

  .main-card {
    max-width: 400px;
    width: 100%;
    text-align: center;
    animation: fadeIn 0.8s ease-out;
  }

  .logo {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    background: linear-gradient(135deg, #fff, #aaa);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: 2px;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .pin-input {
    text-align: center;
    font-weight: 900;
    font-size: 1.5rem;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
