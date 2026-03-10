<script lang="ts">
  import { connect, gameStatus } from "$lib/socketStore";
  import { goto } from "$app/navigation";
  import { alertStore } from "$lib/alertStore";

  let roomID = "";
  let username = "";

  $: if ($gameStatus === "LOBBY") {
    goto("/game");
  }

  function joinGame() {
    if (!roomID) {
      alertStore.showAlert("กรุณากรอก Game PIN");
      return;
    }
    if (!username) {
      alertStore.showAlert("กรุณากรอก NickName");
      return;
    }
    connect(roomID, username);
  }
</script>

<main class="container">
  <div class="glass-card main-card">
    <div class="logo-wrapper">
      <img
        src="/image/logotttspacequiz.png"
        alt="TTT SPACE QUIZ"
        class="system-logo"
      />
    </div>

    <div class="form-group">
      <input
        type="text"
        placeholder="GAME PIN"
        bind:value={roomID}
        class="input-field pin-input"
      />
      <input
        type="text"
        placeholder="NICKNAME"
        bind:value={username}
        class="input-field nickname-input"
      />
      <div class="remark">
        <p>REMARK: MAX 12 CHARS, NO SPECIALS</p>
      </div>
      <button class="btn-primary" on:click={joinGame}>Enter Mission</button>
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

  .logo-wrapper {
    margin-bottom: 2rem;
    display: flex;
    justify-content: center;
    animation: float 4s ease-in-out infinite;
  }

  .system-logo {
    max-width: 180px;
    height: auto;
    filter: drop-shadow(0 0 15px rgba(0, 210, 255, 0.4));
  }

  @keyframes float {
    0%,
    100% {
      transform: translateY(0) scale(1.02);
    }
    50% {
      transform: translateY(-12px) scale(0.98);
    }
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .pin-input,
  .nickname-input {
    text-align: center;
    font-weight: 900;
  }

  .remark {
    margin-top: -0.5rem;
    padding: 0 0.5rem;

    p {
      color: var(--text-secondary);
      font-size: 0.55rem;
      font-family: var(--font-pixel);
      line-height: 1.4;
      opacity: 0.8;
      text-transform: uppercase;
    }
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
