<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import {
        playersStore,
        gameStatus,
        currentQuestion,
        answersStore,
        connect,
        sendMessage,
    } from "$lib/socketStore";
    import { fade, fly, scale } from "svelte/transition";

    interface Question {
        text: string;
        options: string[];
        correct_answer: number;
        image_url?: string;
        time_limit?: number;
    }

    interface Quiz {
        id: string;
        title: string;
        questions: Question[];
        room_id: string;
    }

    let quiz: Quiz | null = null;
    let loading = true;
    let error: string | null = null;
    let timeLeft = 0;
    let timerInterval: any;
    let currentQuestionIndex = -1; // -1 means Lobby
    let playerScores: Record<string, number> = {};
    let answersReceived = 0;
    let answeredThisRound = new Set<string>();

    // Initializing scores for new players reactively but safely
    $: {
        const currentPlayers = $playersStore;
        let changed = false;
        currentPlayers.forEach((player) => {
            if (playerScores[player] === undefined) {
                playerScores[player] = 0;
                changed = true;
            }
        });
        if (changed) playerScores = { ...playerScores };
    }

    // Separate logic for handling incoming answers
    $: if ($answersStore && $gameStatus === "QUESTION") {
        const { username, answer_index } = $answersStore;
        if (!answeredThisRound.has(username)) {
            answeredThisRound.add(username);
            answersReceived++;

            if (answer_index === $currentQuestion?.correct_answer) {
                playerScores[username] = (playerScores[username] || 0) + 1000;
                playerScores = { ...playerScores };
            }
        }
    }

    async function fetchQuiz() {
        const id = $page.params.quizId;
        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch(
                `http://127.0.0.1:8081/api/games/${id}?host=true`,
                {
                    headers: { Authorization: `Bearer ${token}` },
                },
            );
            if (res.ok) {
                quiz = await res.json();
                // Connect to socket as HOST
                connect(quiz.room_id, "HOST");
            } else {
                throw new Error("Failed to fetch quiz");
            }
        } catch (err: any) {
            error = err.message;
        } finally {
            loading = false;
        }
    }

    // Subscribe to socket messages for answers
    onMount(() => {
        fetchQuiz();
    });

    onDestroy(() => {
        clearInterval(timerInterval);
    });

    function startRun() {
        if (!quiz) return;
        currentQuestionIndex = 0;
        answersReceived = 0;
        answeredThisRound = new Set();
        const q = quiz.questions[currentQuestionIndex];
        gameStatus.set("QUESTION");
        currentQuestion.set(q);
        timeLeft = q.time_limit || 20;

        sendMessage("start_game", {
            room_id: quiz.room_id,
            question: q,
        });

        startTimer();
    }

    function startTimer() {
        clearInterval(timerInterval);
        timerInterval = setInterval(() => {
            if (timeLeft > 0) {
                timeLeft--;
            } else {
                showResults();
            }
        }, 1000);
    }

    function showResults() {
        clearInterval(timerInterval);
        gameStatus.set("RESULT");
        sendMessage("show_results", { room_id: quiz?.room_id });
    }

    function nextQuestion() {
        if (!quiz) return;
        currentQuestionIndex++;
        if (currentQuestionIndex < quiz.questions.length) {
            const q = quiz.questions[currentQuestionIndex];
            gameStatus.set("QUESTION");
            currentQuestion.set(q);
            timeLeft = q.time_limit || 20;
            answersReceived = 0;
            answeredThisRound = new Set();

            sendMessage("next_question", {
                room_id: quiz.room_id,
                question: q,
            });
            startTimer();
        } else {
            finishGame();
        }
    }

    function finishGame() {
        gameStatus.set("FINISHED");
        sendMessage("finish_game", { room_id: quiz?.room_id });
        saveGameHistory();
    }

    async function saveGameHistory() {
        if (!quiz) return;

        const leaderboard = Object.entries(playerScores)
            .map(([username, score]) => ({ username, score }))
            .sort((a, b) => b.score - a.score);

        const historyData = {
            quiz_id: quiz.id,
            quiz_title: quiz.title,
            room_id: quiz.room_id,
            players_count: $playersStore.length,
            leaderboard: leaderboard,
        };

        const token = localStorage.getItem("admin_token");
        try {
            await fetch("http://localhost:8081/api/history", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(historyData),
            });
            console.log("History saved successfully");
        } catch (e) {
            console.error("Failed to save history:", e);
        }
    }
</script>

<div class="host-container" in:fade>
    {#if loading}
        <div class="center">
            <div class="spinner"></div>
            <p>Setting up room...</p>
        </div>
    {:else if error}
        <div class="center">
            <p class="error">{error}</p>
            <button class="btn-3d purple" on:click={() => goto("/admin/start")}
                >Go Back</button
            >
        </div>
    {:else if quiz}
        {#if currentQuestionIndex === -1}
            <!-- LOBBY VIEW -->
            <div class="lobby-view" in:scale>
                <header class="lobby-header">
                    <div class="join-info">
                        <h2>Join with Game PIN:</h2>
                        <h1 class="room-pin">{quiz.room_id}</h1>
                    </div>
                </header>

                <div class="players-section">
                    <div class="players-count">
                        <span class="count">{$playersStore.length}</span>
                        <p>Players joined</p>
                    </div>

                    <div class="players-grid">
                        {#each $playersStore as player, i}
                            <div
                                class="player-card"
                                in:fly={{ y: 20, delay: i * 50 }}
                            >
                                <span class="avatar">👤</span>
                                <h3>{player}</h3>
                            </div>
                        {:else}
                            <div class="waiting-text">
                                <p>Waiting for players to join...</p>
                            </div>
                        {/each}
                    </div>
                </div>

                <footer class="lobby-footer">
                    <button
                        class="btn-3d purple cancel-btn"
                        on:click={() => {
                            sendMessage("cancel_game", {
                                room_id: quiz.room_id,
                            });
                            goto("/admin/start");
                        }}
                    >
                        CANCEL ✖
                    </button>
                    <button
                        class="btn-3d blue run-btn"
                        on:click={startRun}
                        disabled={$playersStore.length === 0}
                    >
                        START RUN 🚀
                    </button>
                </footer>
            </div>
        {:else if $gameStatus === "QUESTION" || $gameStatus === "RESULT"}
            <!-- GAME VIEW -->
            <div class="game-view" in:fade>
                <header class="game-header">
                    <div class="q-progress">
                        Question {currentQuestionIndex + 1} of {quiz.questions
                            .length}
                    </div>
                    <div class="timer {timeLeft < 5 ? 'critical' : ''}">
                        {timeLeft}
                    </div>
                    <div class="ans-count">
                        <span class="num">{answersReceived}</span>
                        <span>Answers</span>
                    </div>
                </header>

                <main class="question-main">
                    <div class="question-card glass-card">
                        <h1>{$currentQuestion?.text}</h1>
                    </div>

                    {#if $currentQuestion?.image_url}
                        <div class="media-area">
                            <img
                                src={$currentQuestion.image_url}
                                alt="Question"
                            />
                        </div>
                    {/if}

                    <div class="options-grid">
                        {#each $currentQuestion?.options || [] as opt, i}
                            <div
                                class="option-item opt-{i} {$gameStatus ===
                                    'RESULT' &&
                                i === $currentQuestion?.correct_answer
                                    ? 'correct'
                                    : ''}"
                            >
                                <span class="opt-shape"></span>
                                <span class="opt-text">{opt}</span>
                            </div>
                        {/each}
                    </div>
                </main>

                {#if $gameStatus === "RESULT"}
                    <div class="result-actions" in:fly={{ y: 50 }}>
                        <button
                            class="btn-3d purple next-btn"
                            on:click={nextQuestion}
                        >
                            {currentQuestionIndex + 1 < quiz.questions.length
                                ? "Next Question ➡️"
                                : "Show Podium 🏆"}
                        </button>
                    </div>
                {/if}
            </div>
        {:else if $gameStatus === "FINISHED"}
            <!-- FINISHED / RANKING VIEW -->
            <div class="finished-view" in:scale>
                <h1>Game Finished! 🏁</h1>

                <div class="leaderboard-section">
                    <h2>Final Ranking 🏆</h2>
                    <div class="leaderboard-list">
                        {#each Object.entries(playerScores).sort((a, b) => b[1] - a[1]) as [username, score], i}
                            <div
                                class="player-rank glass-card"
                                in:fly={{ x: -20, delay: i * 100 }}
                            >
                                <div class="rank-info">
                                    <span class="rank-num">#{i + 1}</span>
                                    <span class="player-name">{username}</span>
                                </div>
                                <div class="score">
                                    {score} <span class="pts">pts</span>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>

                <div class="actions">
                    <button
                        class="btn-3d purple"
                        on:click={() => goto("/admin")}
                        >Back to Dashboard</button
                    >
                    <button
                        class="btn-3d blue"
                        on:click={() => goto("/admin/history")}
                        >View history</button
                    >
                </div>
            </div>
        {/if}
    {/if}
</div>

<style lang="scss">
    .host-container {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: radial-gradient(circle at top left, #46178f, #250818);
        z-index: 1000;
        display: flex;
        flex-direction: column;
        color: white;
    }

    .center {
        margin: auto;
        text-align: center;
        .error {
            color: #ff4d4d;
            font-size: 1.2rem;
            margin-bottom: 2rem;
        }
    }

    // Lobby styles
    .lobby-view {
        height: 100%;
        display: flex;
        flex-direction: column;
    }

    .lobby-header {
        background: rgba(0, 0, 0, 0.3);
        padding: 3rem;
        text-align: center;
        .join-info {
            h2 {
                margin: 0;
                opacity: 0.7;
                font-size: 1.5rem;
            }
            .room-pin {
                font-size: 5rem;
                margin: 0.5rem 0;
                font-weight: 900;
                letter-spacing: 5px;
            }
        }
    }

    .players-section {
        flex: 1;
        padding: 3rem;
        display: flex;
        flex-direction: column;
        align-items: center;

        .players-count {
            text-align: center;
            margin-bottom: 3rem;
            .count {
                font-size: 4rem;
                font-weight: 800;
                background: var(--primary);
                padding: 0.5rem 2rem;
                border-radius: 20px;
                box-shadow: 0 10px 0 rgba(0, 0, 0, 0.2);
            }
            p {
                margin: 1rem 0 0;
                font-weight: 700;
                font-size: 1.2rem;
            }
        }

        .players-grid {
            width: 100%;
            display: flex;
            flex-wrap: wrap;
            gap: 1.5rem;
            justify-content: center;
        }

        .player-card {
            background: rgba(255, 255, 255, 0.1);
            padding: 1rem 2rem;
            border-radius: 12px;
            display: flex;
            align-items: center;
            gap: 1rem;
            font-size: 1.2rem;
            font-weight: 700;
            .avatar {
                font-size: 1.5rem;
            }
        }
    }

    .lobby-footer {
        padding: 3rem;
        display: flex;
        justify-content: center;
        gap: 2rem;
        .run-btn,
        .cancel-btn {
            padding: 1.5rem 4rem;
            font-size: 1.5rem;
            min-width: 250px;
        }
    }

    // Game view styles
    .game-view {
        flex: 1;
        display: flex;
        flex-direction: column;
        padding: 2rem;
    }

    .game-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 3rem;

        .timer {
            width: 100px;
            height: 100px;
            background: #46178f;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 3rem;
            font-weight: 900;
            border: 5px solid rgba(255, 255, 255, 0.2);
            &.critical {
                background: #e21b3c;
                animation: pulse 0.5s infinite alternate;
            }
        }

        .q-progress,
        .ans-count {
            font-size: 1.5rem;
            font-weight: 800;
        }
        .ans-count {
            text-align: right;
            .num {
                display: block;
                font-size: 2.5rem;
                color: var(--accent);
            }
        }
    }

    .question-main {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 2rem;
        max-width: 1100px;
        margin: 0 auto;
        width: 100%;

        .question-card {
            background: white;
            color: black;
            padding: 2.5rem;
            text-align: center;
            border-radius: 20px;
            h1 {
                margin: 0;
                font-size: 2.2rem;
                line-height: 1.2;
            }
        }

        .options-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 1rem;
            margin-top: auto;
            margin-bottom: 2rem;
        }

        .option-item {
            padding: 1.5rem 2rem;
            border-radius: 12px;
            display: flex;
            align-items: center;
            gap: 1.5rem;
            font-size: 1.4rem;
            font-weight: 700;
            opacity: 0.8;
            &.correct {
                opacity: 1;
                border: 4px solid white;
                box-shadow: 0 0 30px rgba(255, 255, 255, 0.5);
            }
        }

        .opt-0 {
            background: #e21b3c;
        }
        .opt-1 {
            background: #1368ce;
        }
        .opt-2 {
            background: #26890c;
        }
        .opt-3 {
            background: #ffa602;
        }
    }

    .result-actions {
        position: absolute;
        bottom: 3rem;
        right: 3rem;
    }

    .finished-view {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 4rem 2rem;
        overflow-y: auto;
        h1 {
            font-size: 4rem;
            margin-bottom: 3rem;
            text-shadow: 0 0 20px rgba(255, 255, 255, 0.3);
        }
    }

    .leaderboard-section {
        width: 100%;
        max-width: 800px;
        margin-bottom: 4rem;
        h2 {
            text-align: center;
            margin-bottom: 2rem;
            font-size: 2rem;
            opacity: 0.8;
        }
    }

    .leaderboard-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .player-rank {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.5rem 2.5rem;
        border-radius: 20px;
        background: rgba(255, 255, 255, 0.1);
        border: 1px solid rgba(255, 255, 255, 0.1);

        .rank-info {
            display: flex;
            align-items: center;
            gap: 2rem;
            .rank-num {
                font-size: 2rem;
                font-weight: 900;
                opacity: 0.3;
                width: 50px;
            }
            .player-name {
                font-size: 1.8rem;
                font-weight: 700;
            }
        }

        .score {
            font-size: 2.2rem;
            font-weight: 900;
            color: var(--accent);
            .pts {
                font-size: 1.2rem;
                opacity: 0.5;
            }
        }

        &:first-child {
            background: rgba(255, 215, 0, 0.1);
            border-color: rgba(255, 215, 0, 0.3);
            transform: scale(1.05);
            margin-bottom: 1rem;
            .rank-num {
                opacity: 1;
                color: #ffd700;
            }
            .score {
                color: #ffd700;
            }
        }
    }

    .actions {
        display: flex;
        gap: 2rem;
        margin-top: auto;
        button {
            padding: 1.5rem 3rem;
            font-size: 1.2rem;
        }
    }

    @keyframes pulse {
        to {
            transform: scale(1.1);
        }
    }
    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
