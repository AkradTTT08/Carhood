<script lang="ts">
    import {
        playersStore,
        currentUser,
        gameStatus,
        currentQuestion,
        sendMessage,
        socketStore,
        disconnect,
    } from "$lib/socketStore";
    import { fade, fly, scale } from "svelte/transition";

    let selectedOption: number | null = null;
    let submitted = false;

    $: if ($gameStatus === "QUESTION") {
        submitted = false;
        selectedOption = null;
    }

    function submitAnswer(index: number) {
        if (submitted) return;
        selectedOption = index;
        submitted = true;

        sendMessage("submit_answer", {
            answer_index: index,
        });
    }

    function leaveGame() {
        disconnect();
        window.location.href = "/";
    }
</script>

<div class="player-container" in:fade>
    {#if $gameStatus === "LOBBY"}
        <div class="lobby-view" in:scale>
            <div class="status-card glass-card">
                <h1 class="logo">TTT Space Quiz</h1>
                <div class="loader"></div>
                <h2>Preparing to Start...</h2>
                <div class="joined-players">
                    <p class="join-count">
                        {$playersStore.length} player{$playersStore.length !== 1
                            ? "s"
                            : ""} joined
                    </p>
                    <div class="players-list-scroll">
                        {#each $playersStore as player}
                            <div
                                class="player-tag {player === $currentUser
                                    ? 'is-self'
                                    : ''}"
                                in:fly={{ y: 10 }}
                            >
                                <span class="dot"></span>
                                {player}
                                {player === $currentUser ? "(You)" : ""}
                            </div>
                        {/each}
                    </div>
                </div>
                <div class="lobby-actions">
                    <button class="btn-3d red small" on:click={leaveGame}
                        >LEAVE ROOM</button
                    >
                </div>
                <p class="instruction">See your name on the host screen!</p>
            </div>
        </div>
    {:else if $gameStatus === "QUESTION"}
        <div class="question-view">
            <header class="q-header">
                <h2>Choose wisely...</h2>
            </header>

            <div class="options-grid">
                {#each $currentQuestion?.options || [] as opt, i}
                    <button
                        class="option-btn opt-{i} {selectedOption === i
                            ? 'selected'
                            : ''} {submitted && selectedOption !== i
                            ? 'dimmed'
                            : ''}"
                        on:click={() => submitAnswer(i)}
                        disabled={submitted}
                    >
                        <span class="shape"></span>
                        <span class="text">{opt}</span>
                    </button>
                {/each}
            </div>

            {#if submitted}
                <div class="submitted-overlay" in:fade>
                    <div class="check">✔</div>
                    <p>Answer Submitted!</p>
                </div>
            {/if}
        </div>
    {:else if $gameStatus === "RESULT"}
        <div class="result-view" in:fade>
            <div class="status-card glass-card">
                <h1>Check the screen!</h1>
                <p>Did you get it right?</p>
            </div>
        </div>
    {:else if $gameStatus === "FINISHED"}
        <div class="finished-view" in:scale>
            <div class="status-card glass-card">
                <h1>Game Over! 🏁</h1>
                <p>Check the leaderboard on the host screen.</p>
                <button
                    class="btn-3d purple"
                    on:click={() => (window.location.href = "/")}
                    >Play Again</button
                >
            </div>
        </div>
    {/if}
</div>

<style lang="scss">
    .player-container {
        height: 100vh;
        width: 100%;
        display: flex;
        flex-direction: column;
        background: radial-gradient(circle at top left, #46178f, #250818);
        color: white;
    }

    .lobby-view,
    .result-view,
    .finished-view {
        margin: auto;
        padding: 2rem;
        width: 100%;
        max-width: 400px;
    }

    .status-card {
        text-align: center;
        padding: 3rem 2rem;
        border-radius: 30px;

        .logo {
            font-size: 2.5rem;
            margin-bottom: 2rem;
        }
        h1 {
            margin: 0;
            font-size: 2rem;
        }
        h2 {
            margin: 1rem 0;
            opacity: 0.8;
        }
        p.instruction {
            opacity: 0.6;
            margin-top: 1.5rem;
            font-size: 0.9rem;
        }

        .joined-players {
            margin: 1.5rem 0;
            background: rgba(255, 255, 255, 0.05);
            padding: 1.5rem;
            border-radius: 20px;
            border: 1px solid rgba(255, 255, 255, 0.1);

            .join-count {
                font-weight: 800;
                color: var(--accent);
                margin-bottom: 1rem;
                font-size: 1.1rem;
                opacity: 1;
            }

            .players-list-scroll {
                display: flex;
                flex-wrap: wrap;
                justify-content: center;
                gap: 0.8rem;
                max-height: 200px;
                overflow-y: auto;
                padding: 0.5rem;

                &::-webkit-scrollbar {
                    width: 4px;
                }
                &::-webkit-scrollbar-thumb {
                    background: rgba(255, 255, 255, 0.2);
                    border-radius: 10px;
                }
            }

            .player-tag {
                background: rgba(255, 255, 255, 0.1);
                padding: 0.6rem 1.2rem;
                border-radius: 50px;
                font-size: 0.95rem;
                font-weight: 700;
                display: flex;
                align-items: center;
                gap: 0.6rem;
                border: 1px solid rgba(255, 255, 255, 0.1);
                transition: all 0.3s ease;

                &.is-self {
                    background: rgba(19, 104, 206, 0.2);
                    border-color: rgba(19, 104, 206, 0.5);
                    box-shadow: 0 0 15px rgba(19, 104, 206, 0.3);
                }

                .dot {
                    width: 8px;
                    height: 8px;
                    background: #4caf50;
                    border-radius: 50%;
                    box-shadow: 0 0 10px #4caf50;
                }
            }
        }

        .lobby-actions {
            display: flex;
            justify-content: center;
            margin-top: 2rem;
            margin-bottom: 1rem;
        }

        .loader {
            width: 60px;
            height: 60px;
            border: 5px solid rgba(255, 255, 255, 0.1);
            border-top: 5px solid var(--accent);
            border-radius: 50%;
            margin: 0 auto 2rem;
            animation: spin 1s linear infinite;
        }
    }

    .question-view {
        flex: 1;
        display: flex;
        flex-direction: column;
        padding: 1rem;
        position: relative;
    }

    .q-header {
        text-align: center;
        padding: 2rem 0;
    }

    .options-grid {
        flex: 1;
        display: grid;
        grid-template-columns: 1fr 1fr;
        grid-template-rows: 1fr 1fr;
        gap: 1rem;
        padding-bottom: 2rem;
    }

    .option-btn {
        border: none;
        border-radius: 15px;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        color: white;
        font-size: 1.5rem;
        font-weight: 800;
        cursor: pointer;
        transition: all 0.2s;

        &:disabled {
            cursor: default;
        }

        &.selected {
            transform: scale(0.95);
            opacity: 1 !important;
            border: 5px solid white;
        }
        &.dimmed {
            opacity: 0.3;
        }

        .shape {
            width: 40px;
            height: 40px;
            border: 4px solid white;
            margin-bottom: 1rem;
        }
    }

    .opt-0 {
        background: #e21b3c;
        .shape {
            border-radius: 5px;
        }
    }
    .opt-1 {
        background: #1368ce;
        .shape {
            border-radius: 50%;
        }
    }
    .opt-2 {
        background: #26890c;
        .shape {
            border-radius: 50%;
            transform: rotate(45deg);
        }
    }
    .opt-3 {
        background: #ffa602;
        .shape {
            border-radius: 0;
            transform: rotate(45deg);
        }
    }

    .submitted-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.8);
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        z-index: 10;
        .check {
            font-size: 5rem;
            color: #26890c;
            margin-bottom: 1rem;
        }
        p {
            font-size: 1.5rem;
            font-weight: 700;
        }
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
