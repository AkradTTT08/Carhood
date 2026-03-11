<script lang="ts">
    import {
        playersStore,
        currentUser,
        gameStatus,
        currentQuestion,
        sendMessage,
        socketStore,
        disconnect,
        allAnswersStore,
        questionProgressStore,
        leaderboardStore,
    } from "$lib/socketStore";
    import { fade, fly, scale } from "svelte/transition";
    import { onDestroy } from "svelte";

    let selectedOption: number | null = null;
    let submitted = false;
    let timeLeft: number = 0;
    let timerInterval: any;
    let showPodium = false;

    $: if ($leaderboardStore && $leaderboardStore.length > 0) {
        localStorage.setItem("last_leaderboard", JSON.stringify($leaderboardStore));
    }

    $: if ($gameStatus === "QUESTION" && $currentQuestion) {
        submitted = false;
        selectedOption = null;
        startTimer($currentQuestion.time_limit || 20);
    } else {
        stopTimer();
    }

    function startTimer(duration: number) {
        stopTimer();
        timeLeft = duration;
        timerInterval = setInterval(() => {
            if (timeLeft > 0) {
                timeLeft--;
            } else {
                stopTimer();
            }
        }, 1000);
    }

    function stopTimer() {
        if (timerInterval) {
            clearInterval(timerInterval);
            timerInterval = null;
        }
    }

    onDestroy(() => {
        stopTimer();
    });

    function submitAnswer(index: number) {
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
        <div class="game-layout">
            <div class="question-view">
                <header class="q-header">
                    <div class="progress-indicator">
                        {$questionProgressStore.current} / {$questionProgressStore.total}
                    </div>
                    <div
                        class="timer {timeLeft <= 5 && timeLeft > 0
                            ? 'danger-pulse'
                            : ''}"
                    >
                        {timeLeft}
                    </div>
                    <h2 class="question-text">
                        {$currentQuestion?.text || "Choose wisely..."}
                    </h2>
                    {#if $currentQuestion?.image_url}
                        <div class="question-image-container">
                            <img
                                src={$currentQuestion.image_url}
                                alt="Question Media"
                                class="question-image"
                            />
                        </div>
                    {/if}
                </header>

                <div class="options-grid">
                    {#each $currentQuestion?.options || [] as opt, i}
                        <button
                            class="option-btn opt-{i} {selectedOption === i
                                ? 'selected'
                                : ''}"
                            on:click={() => submitAnswer(i)}
                        >
                            <span class="shape"></span>
                            <span class="text">{opt}</span>
                        </button>
                    {/each}
                </div>
            </div>

            <div class="stats-panel glass-card">
                <h3>Live Answers</h3>
                <div class="stats-list">
                    {#each $currentQuestion?.options || [] as _, i}
                        <div class="stat-item opt-{i}">
                            <div class="stat-shape"></div>
                            <span class="stat-count">
                                {Object.values($allAnswersStore).filter(
                                    (v) => v === i,
                                ).length}
                            </span>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    {:else if $gameStatus === "RESULT"}
        <div class="result-view" in:fade>
            <div class="status-card glass-card">
                <h1>
                    {selectedOption === $currentQuestion?.correct_answer
                        ? "Correct! 🎉"
                        : "Incorrect! ❌"}
                </h1>
                <p>The correct answer was:</p>
                <div
                    class="correct-answer-display opt-{$currentQuestion?.correct_answer}"
                >
                    <span class="correct-text"
                        >{$currentQuestion?.options[
                            $currentQuestion?.correct_answer
                        ]}</span
                    >
                </div>
            </div>
        </div>
    {:else if $gameStatus === "FINISHED"}
        <div class="finished-view" in:scale>
            {#if !showPodium}
                <div class="status-card glass-card" in:fade>
                    <div class="finish-icon">🏆</div>
                    <h1>Game Over! 🏁</h1>
                    <p>Great job participating in the quiz!</p>
                    
                    <div class="finish-actions">
                        <button
                            class="btn-3d purple"
                            on:click={leaveGame}
                        >
                            Play Again 🎮
                        </button>
                        <button
                            class="btn-3d blue"
                            on:click={() => (showPodium = true)}
                        >
                            View Podium 🏅
                        </button>
                    </div>
                </div>
            {:else}
                <div class="podium-view" in:fly={{ y: 20 }}>
                    <header class="podium-header-scaled">
                        <button class="back-link-scaled" on:click={() => showPodium = false}>
                            ← Back
                        </button>
                        <h1>Game Finished! 🏁</h1>
                        <h2>Final Ranking 🏆</h2>
                    </header>

                    <div class="podium-list">
                        {#each $leaderboardStore as entry, i}
                            <div class="podium-item {entry.username === $currentUser ? 'is-self' : ''}" 
                                 in:fly={{ x: -20, delay: i * 50 }}>
                                <div class="rank-badge {i < 3 ? 'top-' + (i + 1) : ''}">
                                    {i + 1}
                                </div>
                                <div class="user-info">
                                    <span class="user-name">{entry.username}</span>
                                    {#if entry.username === $currentUser}
                                        <span class="self-label">(You)</span>
                                    {/if}
                                </div>
                                <div class="user-score">
                                    {entry.score} <span>pts</span>
                                </div>
                            </div>
                        {/each}
                    </div>

                    <div class="podium-footer-scaled">
                        <button
                            class="btn-3d purple"
                            on:click={() => (window.location.href = "/")}
                        >
                            Play Again 🎮
                        </button>
                    </div>
                </div>
            {/if}
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
        max-width: 800px;
        transition: max-width 0.5s ease;
    }

    .finished-view {
        .finish-icon {
            font-size: 5rem;
            margin-bottom: 1rem;
            animation: bounce 2s infinite;
        }

        .finish-actions {
            display: flex;
            flex-direction: column;
            gap: 1rem;
            margin-top: 2rem;
            width: 100%;

            button {
                width: 100%;
                padding: 1.2rem;
                font-size: 1.2rem;
            }
        }
    }

    .podium-view {
        width: 100%;
        max-width: 1000px;
        margin: 0 auto;
        display: flex;
        flex-direction: column;
        align-items: center;

        .podium-header-scaled {
            text-align: center;
            margin-bottom: 3rem;
            position: relative;
            width: 100%;

            .back-link-scaled {
                position: absolute;
                left: 0;
                top: 50%;
                transform: translateY(-50%);
                background: rgba(255, 255, 255, 0.1);
                border: none;
                color: white;
                padding: 0.8rem 1.5rem;
                border-radius: 12px;
                font-weight: 700;
                cursor: pointer;
                transition: all 0.2s;

                &:hover {
                    background: rgba(255, 255, 255, 0.2);
                    transform: translateY(-50%) translateX(-5px);
                }
            }

            h1 {
                font-size: 4rem;
                margin: 0 0 1rem 0;
                text-shadow: 0 0 20px rgba(255, 255, 255, 0.3);
            }

            h2 {
                font-size: 2rem;
                margin: 0;
                opacity: 0.8;
                color: #fff;
            }
        }

        .podium-list {
            width: 100%;
            display: flex;
            flex-direction: column;
            gap: 1.2rem;
            margin-bottom: 3rem;
        }

        .podium-item {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 2rem;
            background: rgba(255, 255, 255, 0.1);
            padding: 1.5rem 3rem;
            border-radius: 25px;
            border: 1px solid rgba(255, 255, 255, 0.1);
            transition: all 0.3s ease;

            &:hover {
                background: rgba(255, 255, 255, 0.15);
                transform: scale(1.01);
            }

            &.is-self {
                background: rgba(19, 104, 206, 0.2);
                border-color: rgba(19, 104, 206, 0.4);
                box-shadow: 0 0 20px rgba(19, 104, 206, 0.2);
            }

            .rank-badge {
                width: 60px;
                height: 60px;
                background: rgba(255, 255, 255, 0.05);
                border-radius: 50%;
                display: flex;
                align-items: center;
                justify-content: center;
                font-weight: 900;
                font-size: 2rem;
                color: rgba(255, 255, 255, 0.3);

                &.top-1 { background: #ffd700; color: #000; box-shadow: 0 0 20px #ffd700; opacity: 1; }
                &.top-2 { background: #c0c0c0; color: #000; box-shadow: 0 0 15px #c0c0c0; opacity: 1; }
                &.top-3 { background: #cd7f32; color: #351c1c; box-shadow: 0 0 15px #cd7f32; opacity: 1; }
            }

            .user-info {
                flex: 1;
                display: flex;
                align-items: center;
                gap: 1rem;
                
                .user-name {
                    font-weight: 700;
                    font-size: 2rem;
                }
                .self-label {
                    background: var(--accent);
                    color: white;
                    padding: 0.2rem 0.8rem;
                    border-radius: 20px;
                    font-size: 0.9rem;
                    font-weight: 800;
                }
            }

            .user-score {
                font-size: 2.5rem;
                font-weight: 900;
                color: var(--accent);
                span { font-size: 1.2rem; opacity: 0.5; font-weight: normal; color: #fff; }
            }
        }

        .podium-footer-scaled {
            width: 100%;
            display: flex;
            justify-content: center;
            button {
                padding: 1.5rem 4rem;
                font-size: 1.5rem;
            }
        }
    }

    @keyframes bounce {
        0%, 20%, 50%, 80%, 100% {transform: translateY(0);}
        40% {transform: translateY(-20px);}
        60% {transform: translateY(-10px);}
    }

    .game-layout {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: 100%;
        gap: 2rem;
        padding: 2rem;
        box-sizing: border-box;
    }

    .question-view {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: center;
        max-width: 800px;
        margin: 0 auto;
        width: 100%;
    }

    .stats-panel {
        width: 250px;
        min-width: 250px;
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 1rem;

        h3 {
            margin: 0;
            text-align: center;
            font-size: 1.2rem;
            opacity: 0.8;
            border-bottom: 2px solid rgba(255, 255, 255, 0.1);
            padding-bottom: 0.5rem;
        }

        .stats-list {
            display: flex;
            flex-direction: column;
            gap: 1rem;
            flex: 1;
            justify-content: center;
        }

        .stat-item {
            display: flex;
            align-items: center;
            gap: 1rem;
            padding: 1rem;
            border-radius: 10px;
            background: rgba(255, 255, 255, 0.1);

            .stat-shape {
                width: 20px;
                height: 20px;
                border: 2px solid white;
            }

            .stat-count {
                font-size: 1.5rem;
                font-weight: 800;
                margin-left: auto;
            }

            &.opt-0 {
                border-left: 5px solid #e21b3c;
                .stat-shape { border-radius: 4px; }
            }
            &.opt-1 {
                border-left: 5px solid #1368ce;
                .stat-shape { border-radius: 50%; }
            }
            &.opt-2 {
                border-left: 5px solid #26890c;
                .stat-shape { border-radius: 50%; transform: rotate(45deg); }
            }
            &.opt-3 {
                border-left: 5px solid #ffa602;
                .stat-shape { transform: rotate(45deg); }
            }
        }
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
            }

            .players-list-scroll {
                display: flex;
                flex-wrap: wrap;
                justify-content: center;
                gap: 1rem;
                max-height: 300px;
                overflow-y: auto;
                padding: 0.5rem;

                &::-webkit-scrollbar { width: 4px; }
                &::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.2); border-radius: 10px; }
            }

            .player-tag {
                background: rgba(255, 255, 255, 0.1);
                padding: 0.8rem 1.5rem;
                border-radius: 50px;
                font-size: 1rem;
                font-weight: 700;
                display: flex;
                align-items: center;
                gap: 0.8rem;
                border: 1px solid rgba(255, 255, 255, 0.1);

                &.is-self {
                    background: rgba(19, 104, 206, 0.2);
                    border-color: rgba(19, 104, 206, 0.5);
                }

                .dot {
                    width: 10px;
                    height: 10px;
                    background: #4caf50;
                    border-radius: 50%;
                    box-shadow: 0 0 10px #4caf50;
                }
            }
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

        .lobby-actions {
            margin-top: 2rem;
            display: flex;
            justify-content: center;
        }
    }

    .q-header {
        text-align: center;
        padding-bottom: 2rem;

        .progress-indicator {
            background: rgba(0, 0, 0, 0.5);
            padding: 0.8rem 1.5rem;
            border-radius: 15px;
            font-size: 1.2rem;
            font-weight: 800;
            margin-bottom: 1rem;
            display: inline-block;
        }
    }

    .timer {
        font-size: 3.5rem;
        font-weight: 900;
        background: rgba(0, 0, 0, 0.4);
        width: 110px;
        height: 110px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 4px solid var(--accent);
        margin: 0 auto 1.5rem;
    }

    .danger-pulse {
        border-color: #e21b3c;
        color: #e21b3c;
        animation: shakeDanger 0.2s infinite alternate;
    }

    @keyframes shakeDanger {
        from { transform: translateX(-2px); }
        to { transform: translateX(2px); }
    }

    .question-text {
        font-size: 1.8rem;
        margin-bottom: 2rem;
    }

    .question-image-container {
        max-width: 500px;
        margin: 0 auto 2rem;
        border-radius: 15px;
        overflow: hidden;
        border: 4px solid rgba(255, 255, 255, 0.2);
    }

    .question-image {
        width: 100%;
        max-height: 250px;
        object-fit: contain;
        background: #fff;
    }

    .options-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }

    .option-btn {
        padding: 2rem;
        border: none;
        border-radius: 15px;
        color: white;
        font-size: 1.5rem;
        font-weight: 700;
        cursor: pointer;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1rem;

        &.selected {
            border: 4px solid white;
            transform: scale(0.98);
        }

        .shape {
            width: 30px;
            height: 30px;
            border: 3px solid white;
        }
    }

    .opt-0 { background: #e21b3c; .shape { border-radius: 4px; } }
    .opt-1 { background: #1368ce; .shape { border-radius: 50%; } }
    .opt-2 { background: #26890c; .shape { border-radius: 50%; transform: rotate(45deg); } }
    .opt-3 { background: #ffa602; .shape { border-radius: 0; transform: rotate(45deg); } }

    .correct-answer-display {
        padding: 2rem;
        border-radius: 15px;
        font-size: 2rem;
        font-weight: 900;
        margin-top: 2rem;
        border: 4px solid white;
    }

    @keyframes spin {
        to { transform: rotate(360deg); }
    }
</style>
