<script lang="ts">
  import { connect, gameStatus } from "$lib/socketStore";
  import { goto } from "$app/navigation";
  import { alertStore } from "$lib/alertStore";

  import { onMount } from "svelte";
  import { fade, fly, scale } from "svelte/transition";

  let roomID = "";
  let username = "";
  let lastLeaderboard: any[] = [];
  let showPodium = false;

  onMount(() => {
    const stored = localStorage.getItem("last_leaderboard");
    if (stored) {
      try {
        lastLeaderboard = JSON.parse(stored);
      } catch (e) {
        console.error("Failed to parse last_leaderboard", e);
      }
    }
  });

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
      
      {#if lastLeaderboard && lastLeaderboard.length > 0}
        <button class="btn-secondary" on:click={() => (showPodium = true)}>
          Last Game Podium 🏅
        </button>
      {/if}
    </div>
  </div>

  {#if showPodium}
    <div class="podium-overlay" transition:fade>
      <div class="podium-modal glass-card" in:fly={{ y: 20 }}>
        <header class="podium-header">
          <button class="close-btn" on:click={() => (showPodium = false)}>×</button>
          <h2>Last Game Results</h2>
        </header>

        <div class="podium-list">
          {#each lastLeaderboard as entry, i}
            <div class="podium-item" in:fly={{ x: -20, delay: i * 50 }}>
              <div class="rank-badge {i < 3 ? 'top-' + (i + 1) : ''}">
                {i + 1}
              </div>
              <div class="user-info">
                <span class="user-name">{entry.username}</span>
              </div>
              <div class="user-score">
                {entry.score} <span>pts</span>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}
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

  .btn-secondary {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: white;
    padding: 1rem;
    border-radius: 12px;
    font-weight: 700;
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 1rem;

    &:hover {
      background: rgba(255, 255, 255, 0.2);
      transform: translateY(-2px);
    }
  }

  .podium-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.85);
    backdrop-filter: blur(8px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
    padding: 2rem;
  }

  .podium-modal {
    width: 100%;
    max-width: 800px;
    max-height: 80vh;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    animation: scaleIn 0.3s ease;

    .podium-header {
      display: flex;
      align-items: center;
      margin-bottom: 2rem;
      position: relative;

      h2 {
        flex: 1;
        text-align: center;
        margin: 0;
        font-size: 2.5rem;
        color: #fff;
      }

      .close-btn {
        position: absolute;
        right: -10px;
        top: -10px;
        background: rgba(255, 255, 255, 0.1);
        border: none;
        color: white;
        width: 35px;
        height: 35px;
        border-radius: 50%;
        font-size: 1.5rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s;

        &:hover {
          background: rgba(255, 0, 0, 0.3);
          transform: rotate(90deg);
        }
      }
    }
  }

  .podium-list {
    flex: 1;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 0.8rem;
    padding-right: 0.5rem;

    &::-webkit-scrollbar {
      width: 5px;
    }
    &::-webkit-scrollbar-thumb {
      background: rgba(255, 255, 255, 0.2);
      border-radius: 10px;
    }
  }

  .podium-item {
    display: flex;
    align-items: center;
    gap: 2rem;
    background: rgba(255, 255, 255, 0.08);
    padding: 1.2rem 2rem;
    border-radius: 20px;
    border: 1px solid rgba(255, 255, 255, 0.05);

    .rank-badge {
      width: 50px;
      height: 50px;
      background: rgba(0, 0, 0, 0.3);
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 900;
      font-size: 1.5rem;

      &.top-1 { background: #ffd700; color: #000; box-shadow: 0 0 15px #ffd700; }
      &.top-2 { background: #c0c0c0; color: #000; box-shadow: 0 0 10px #c0c0c0; }
      &.top-3 { background: #cd7f32; color: #351c1c; box-shadow: 0 0 10px #cd7f32; }
    }

    .user-info {
      flex: 1;
      .user-name {
        font-weight: 700;
        font-size: 1.6rem;
      }
    }

    .user-score {
      font-size: 2rem;
      font-weight: 900;
      color: var(--accent);
      span { font-size: 1rem; opacity: 0.6; color: #fff; }
    }
  }

  @keyframes scaleIn {
    from { transform: scale(0.9); opacity: 0; }
    to { transform: scale(1); opacity: 1; }
  }
</style>
