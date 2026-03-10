<script>
  import "../app.scss";
  import ThemedAlert from "$lib/components/ThemedAlert.svelte";
</script>

<div class="main-bg">
  <img src="/image/BGADMIN.svg" alt="Background" class="bg-image" />
</div>

<main class="app-content">
  <slot />
</main>

<ThemedAlert />

<style lang="scss">
  :global(body) {
    margin: 0;
    padding: 0;
    background: #000;
  }

  .app-content {
    position: relative;
    z-index: 10;
  }
</style>
