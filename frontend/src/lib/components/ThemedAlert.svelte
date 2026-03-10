<script lang="ts">
    import { alertStore } from "$lib/alertStore";
    import { fade, fly } from "svelte/transition";

    $: ({ show, message, type } = $alertStore);

    function close() {
        alertStore.hideAlert();
    }
</script>

{#if show}
    <div
        class="alert-overlay"
        on:click|self={close}
        on:keydown={(e) => e.key === "Escape" && close()}
        role="button"
        tabindex="0"
        in:fade={{ duration: 200 }}
        out:fade={{ duration: 200 }}
    >
        <div
            class="alert-card glass-card {type}"
            in:fly={{ y: -50, duration: 400 }}
            out:fly={{ y: -50, duration: 200 }}
        >
            <div class="alert-content">
                <div class="icon-wrap">
                    {#if type === "error"}
                        <span class="pixel-icon">!</span>
                    {:else if type === "success"}
                        <span class="pixel-icon">✓</span>
                    {:else}
                        <span class="pixel-icon">?</span>
                    {/if}
                </div>
                <p class="message">{message}</p>
                <button class="btn-primary" on:click={close}>OK</button>
            </div>
        </div>
    </div>
{/if}

<style lang="scss">
    .alert-overlay {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.7);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 9999;
        backdrop-filter: blur(5px);
    }

    .alert-card {
        max-width: 320px;
        width: 90%;
        text-align: center;
        border: 4px solid #fff;
        image-rendering: pixelated;

        &.error {
            border-color: #ff3e3e;
            box-shadow: 0 0 20px rgba(255, 62, 62, 0.4);
        }

        &.success {
            border-color: #39ff14;
            box-shadow: 0 0 20px rgba(57, 255, 20, 0.4);
        }

        &.info {
            border-color: #00d2ff;
            box-shadow: 0 0 20px rgba(0, 210, 255, 0.4);
        }
    }

    .alert-content {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1.5rem;
    }

    .icon-wrap {
        width: 50px;
        height: 50px;
        background: rgba(255, 255, 255, 0.1);
        display: flex;
        align-items: center;
        justify-content: center;
        border: 4px solid currentColor;

        .pixel-icon {
            font-family: var(--font-pixel);
            font-size: 1.5rem;
            font-weight: bold;
        }
    }

    .error .icon-wrap {
        color: #ff3e3e;
    }
    .success .icon-wrap {
        color: #39ff14;
    }
    .info .icon-wrap {
        color: #00d2ff;
    }

    .message {
        font-family: var(--font-pixel);
        font-size: 0.7rem;
        line-height: 1.6;
        color: var(--text-primary);
        text-transform: uppercase;
    }

    .btn-primary {
        width: 100%;
        font-size: 0.7rem;
        padding: 0.8rem;
    }
</style>
