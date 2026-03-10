<script>
    import { createEventDispatcher } from "svelte";
    const dispatch = createEventDispatcher();

    export let show = false;
    export let title = "Confirm";
    export let message = "";
    export let confirmText = "OK";
    export let cancelText = "Cancel";
    export let type = "info"; // info, warning, danger
    export let classic = false; // New prop for professional style

    function handleConfirm() {
        dispatch("confirm");
    }

    function handleCancel() {
        dispatch("cancel");
        show = false;
    }
</script>

{#if show}
    <div class="modal-backdrop" on:click={handleCancel}>
        <div
            class="modal-container glass-card {classic ? 'classic' : ''}"
            on:click|stopPropagation
        >
            <div class="modal-header">
                <h2 class={type}>{title}</h2>
                <button class="close-btn" on:click={handleCancel}
                    >&times;</button
                >
            </div>
            <div class="modal-body">
                <p>{message}</p>
            </div>
            <div class="modal-footer">
                <button
                    class={classic ? "btn-classic gray" : "btn-3d"}
                    on:click={handleCancel}>{cancelText}</button
                >
                <button
                    class="{classic ? 'btn-classic' : 'btn-3d'} {type ===
                    'danger'
                        ? 'red'
                        : type === 'warning'
                          ? 'yellow'
                          : 'blue'}"
                    on:click={handleConfirm}
                >
                    {confirmText}
                </button>
            </div>
        </div>
    </div>
{/if}

<style lang="scss">
    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100vw;
        height: 100vh;
        background: rgba(0, 0, 0, 0.6);
        backdrop-filter: blur(4px);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 2000;
        animation: fadeIn 0.3s ease;
    }

    .modal-container {
        width: 90%;
        max-width: 450px;
        background: rgba(255, 255, 255, 0.95);
        border: 1px solid rgba(255, 255, 255, 0.3);
        padding: 2rem;
        color: #333;
        transform: scale(1);
        animation: popIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);

        &.classic {
            border-radius: 12px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
            font-family: var(--font-main);
            background: #fff;
            border: 1px solid #ddd;

            &::after {
                display: none;
            }

            h2 {
                font-family: var(--font-main);
                font-weight: 700;
                letter-spacing: 0;
                font-size: 1.4rem;
            }

            p {
                font-family: var(--font-main);
                color: #444;
            }

            .btn-classic {
                font-family: var(--font-main);
                font-weight: 600;
                padding: 0.7rem 1.4rem;
                border-radius: 8px;
                border: none;
                cursor: pointer;
                transition: all 0.2s;
                font-size: 0.95rem;

                &.gray {
                    background: #f2f2f2;
                    color: #555;
                    &:hover {
                        background: #e8e8e8;
                    }
                }

                &.blue {
                    background: #1368ce;
                    color: white;
                    &:hover {
                        background: #0f56aa;
                    }
                }

                &.red {
                    background: #e21b3c;
                    color: white;
                    &:hover {
                        background: #d01937;
                    }
                }

                &.yellow {
                    background: #d89e00;
                    color: white;
                    &:hover {
                        background: #c08d00;
                    }
                }
            }
        }

        .modal-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1.5rem;

            h2 {
                margin: 0;
                font-size: 1.5rem;
                &.danger {
                    color: #e21b3c;
                }
                &.warning {
                    color: #d89e00;
                }
                &.info {
                    color: #1368ce;
                }
            }

            .close-btn {
                background: none;
                border: none;
                font-size: 2rem;
                cursor: pointer;
                color: #666;
                padding: 0;
                line-height: 1;
                &:hover {
                    color: #000;
                }
            }
        }

        .modal-body {
            margin-bottom: 2rem;
            p {
                font-size: 1.1rem;
                line-height: 1.6;
                color: #555;
            }
        }

        .modal-footer {
            display: flex;
            justify-content: flex-end;
            gap: 1rem;

            button {
                padding: 0.8rem 1.5rem;
                font-size: 0.9rem;
                min-width: 100px;
                border: none;
                cursor: pointer;
            }
        }
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }

    @keyframes popIn {
        from {
            opacity: 0;
            transform: scale(0.9);
        }
        to {
            opacity: 1;
            transform: scale(1);
        }
    }
</style>
