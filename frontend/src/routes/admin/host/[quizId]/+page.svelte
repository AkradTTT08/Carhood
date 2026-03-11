<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import {
        playersStore,
        gameStatus,
        currentQuestion,
        answersStore,
        allAnswersStore,
        connect,
        sendMessage,
    } from "$lib/socketStore";
    import { fade, fly, scale } from "svelte/transition";
    import Modal from "$lib/components/Modal.svelte";
    import { API_URL } from "$lib/api";

    interface Question {
        text: string;
        options: string[];
        correct_answer: number;
        image_url?: string;
        time_limit?: number;
        points?: number;
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
    let pendingScores: Record<string, number> = {};
    let answersReceived = 0;
    let answeredThisRound = new Set<string>();
    let showCancelModal = false;

    // Reactively save state whenever crucial variables change
    $: {
        if (quiz && $gameStatus !== "LOBBY" && typeof window !== "undefined") {
            const state = {
                status: $gameStatus,
                qIndex: currentQuestionIndex,
                scores: playerScores,
                ansReceived: answersReceived,
                answered: Array.from(answeredThisRound),
                timeLeft: timeLeft,
                allAnswers: $allAnswersStore,
            };
            localStorage.setItem(
                `host_state_${quiz.id}_${quiz.room_id}`,
                JSON.stringify(state),
            );
        }
    }

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

    // Subscription for handling incoming answers reliably
    let answersUnsubscribe: () => void;

    function subscribeToAnswers() {
        if (answersUnsubscribe) answersUnsubscribe();

        answersUnsubscribe = answersStore.subscribe(($answers) => {
            if (!$answers || $gameStatus !== "QUESTION") return;

            const { username, answer_index } = $answers;

            // Allow 1s grace period after timer hits 0
            if (
                !answeredThisRound.has(username) &&
                (timeLeft > 0 || (timeLeft === 0 && !timerInterval))
            ) {
                answeredThisRound.add(username);
                answersReceived++;

                if (answer_index === $currentQuestion?.correct_answer) {
                    const basePoints = $currentQuestion?.points || 1000;
                    const totalTime = $currentQuestion?.time_limit || 20;

                    // High precision: Base * 1000 + Speed Bonus (0-500)
                    // This ensures even with 1 base point, we can rank by speed.
                    const speedBonus = Math.round(500 * (timeLeft / totalTime));
                    const pointsToAdd = basePoints * 1000 + speedBonus;

                    pendingScores[username] =
                        (pendingScores[username] || 0) + pointsToAdd;
                }
            }
        });
    }

    async function fetchQuiz() {
        const urlParams = new URLSearchParams(window.location.search);
        const quizId = urlParams.get("quiz");
        const pin = urlParams.get("pin");

        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch(`${API_URL}/api/games/${quizId}`, {
                headers: { Authorization: `Bearer ${token}` },
            });
            if (res.ok) {
                const data = await res.json();
                // Override the template's room_id with our dynamic session PIN
                if (pin) data.room_id = pin;
                quiz = data;

                // Connect to socket as HOST
                connect(data.room_id, "HOST");

                // Restore state if exists
                if (typeof window !== "undefined") {
                    const saved = localStorage.getItem(
                        `host_state_${quiz.id}_${quiz.room_id}`,
                    );
                    if (saved && quiz) {
                        try {
                            const state = JSON.parse(saved);
                            if (state.status && state.status !== "LOBBY") {
                                currentQuestionIndex = state.qIndex;
                                playerScores = state.scores || {};
                                answersReceived = state.ansReceived || 0;
                                answeredThisRound = new Set(
                                    state.answered || [],
                                );
                                if (state.allAnswers) {
                                    allAnswersStore.set(state.allAnswers);
                                }
                                gameStatus.set(state.status);

                                if (
                                    state.status === "QUESTION" ||
                                    state.status === "RESULT"
                                ) {
                                    currentQuestion.set(
                                        quiz.questions[currentQuestionIndex],
                                    );
                                    if (state.status === "QUESTION") {
                                        timeLeft =
                                            state.timeLeft ||
                                            quiz.questions[currentQuestionIndex]
                                                .time_limit ||
                                            20;
                                        startTimer();
                                    }
                                }
                            }
                        } catch (e) {
                            console.error("Failed to restore state", e);
                        }
                    }
                }
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
        subscribeToAnswers();
    });

    onDestroy(() => {
        clearInterval(timerInterval);
        if (answersUnsubscribe) answersUnsubscribe();
    });

    function startRun() {
        if (!quiz) return;
        currentQuestionIndex = 0;
        answersReceived = 0;
        answeredThisRound = new Set();
        pendingScores = {};
        const q = quiz.questions[currentQuestionIndex];
        gameStatus.set("QUESTION");
        currentQuestion.set(q);
        timeLeft = q.time_limit || 20;

        sendMessage("start_game", {
            room_id: quiz.room_id,
            question: q,
            progress: {
                current: currentQuestionIndex + 1,
                total: quiz.questions.length,
            },
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
        timerInterval = null; // Important for grace period check
        gameStatus.set("RESULT");

        // Commit pending scores to the main scores board
        Object.entries(pendingScores).forEach(([username, points]) => {
            playerScores[username] = (playerScores[username] || 0) + points;
        });

        // Trigger Svelte refresh and clear buffer
        playerScores = { ...playerScores };
        pendingScores = {};

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
            pendingScores = {};

            sendMessage("next_question", {
                room_id: quiz.room_id,
                question: q,
                progress: {
                    current: currentQuestionIndex + 1,
                    total: quiz.questions.length,
                },
            });
            startTimer();
        } else {
            finishGame();
        }
    }

    function finishGame() {
        const leaderboard = Object.entries(playerScores)
            .filter(([username]) => $playersStore.includes(username))
            .map(([username, score]) => ({ username, score }))
            .sort((a, b) => b.score - a.score);

        gameStatus.set("FINISHED");
        sendMessage("finish_game", {
            room_id: quiz?.room_id,
            leaderboard: leaderboard,
        });
        saveGameHistory(leaderboard);
    }

    function confirmCancel() {
        if (!quiz) return;
        sendMessage("cancel_game", {
            room_id: quiz.room_id,
        });
        localStorage.removeItem(`host_state_${quiz.id}_${quiz.room_id}`);
        goto("/admin/start");
    }

    async function saveGameHistory(
        leaderboard: { username: string; score: number }[],
    ) {
        if (!quiz) return;

        const historyData = {
            quiz_id: quiz.id,
            quiz_title: quiz.title,
            room_id: quiz.room_id,
            players_count: $playersStore.length,
            leaderboard: leaderboard,
        };

        const token = localStorage.getItem("admin_token");
        try {
            await fetch(`${API_URL}/api/history`, {
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
                        on:click={() => (showCancelModal = true)}
                    >
                        CANCEL ✖
                    </button>
                    <button
                        class="btn-3d blue run-btn"
                        on:click={startRun}
                        disabled={!quiz || $playersStore.length === 0}
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

                <div class="game-content-wrapper">
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
                                        : ''} {$gameStatus === 'RESULT' &&
                                    i !== $currentQuestion?.correct_answer
                                        ? 'dimmed'
                                        : ''}"
                                >
                                    <span class="opt-shape"></span>
                                    <span class="opt-text">{opt}</span>
                                    {#if $gameStatus === "RESULT"}
                                        <span class="opt-count" in:scale>
                                            👤 {Object.values(
                                                $allAnswersStore,
                                            ).filter((v) => v === i).length}
                                        </span>
                                    {/if}
                                </div>
                            {/each}
                        </div>
                    </main>

                    <div class="right-column">
                        <aside class="live-ranking-sidebar glass-card">
                            <h3>Top 10 Live Ranking 🏆</h3>
                            <div class="ranking-list">
                                {#each Object.entries(playerScores)
                                    .filter( ([username]) => $playersStore.includes(username), )
                                    .sort((a, b) => b[1] - a[1])
                                    .slice(0, 10) as [username, score], i}
                                    <div class="rank-item">
                                        <span class="rank-num">#{i + 1}</span>
                                        <span class="rank-name">{username}</span
                                        >
                                        <span class="rank-score">
                                            {Math.floor(score / 1000)}
                                        </span>
                                    </div>
                                {/each}
                            </div>
                        </aside>

                        {#if $gameStatus === "RESULT"}
                            <div class="result-actions" in:fly={{ y: 20 }}>
                                <button
                                    class="btn-3d purple next-btn"
                                    on:click={nextQuestion}
                                >
                                    {currentQuestionIndex + 1 <
                                    quiz.questions.length
                                        ? "Next Question ➡️"
                                        : "Show Podium 🏆"}
                                </button>
                            </div>
                        {/if}
                    </div>
                </div>
            </div>
        {:else if $gameStatus === "FINISHED"}
            <!-- FINISHED / RANKING VIEW -->
            <div class="finished-view" in:scale>
                <h1>Game Finished! 🏁</h1>

                <div class="leaderboard-section">
                    <h2>Final Ranking 🏆</h2>
                    <div class="leaderboard-list">
                        {#each Object.entries(playerScores)
                            .filter( ([username]) => $playersStore.includes(username), )
                            .sort((a, b) => b[1] - a[1]) as [username, score], i}
                            <div
                                class="player-rank glass-card"
                                in:fly={{ x: -20, delay: i * 100 }}
                            >
                                <div class="rank-info">
                                    <span class="rank-num">#{i + 1}</span>
                                    <span class="player-name">{username}</span>
                                </div>
                                <div class="score">
                                    {Math.floor(score / 1000)}
                                    <span class="pts">pts</span>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>

                <div class="actions">
                    <button
                        class="btn-3d purple"
                        on:click={() => {
                            if (quiz)
                                localStorage.removeItem(
                                    `host_state_${quiz.id}_${quiz.room_id}`,
                                );
                            goto("/admin");
                        }}>Back to Dashboard</button
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

    <Modal
        bind:show={showCancelModal}
        title="Confirm Cancellation"
        message="คุณต้องการยกเลิกห้องเล่นเกมนี้หรือไม่?"
        confirmText="ตกลง"
        cancelText="ไม่ใช่"
        type="danger"
        on:confirm={confirmCancel}
    />
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
        padding-right: 420px;
        position: relative;
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

    .game-content-wrapper {
        display: flex;
        gap: 1rem;
        width: 100%;
        padding: 0;
        height: calc(100vh - 180px);
    }

    .right-column {
        width: 380px;
        min-width: 380px;
        display: flex;
        flex-direction: column;
        gap: 0;
        height: calc(100vh - 3rem);
        position: absolute;
        top: 1.5rem;
        right: 1.5rem;
        bottom: 1.5rem;
    }

    .live-ranking-sidebar {
        flex: 1;
        background: rgba(0, 0, 0, 0.4);
        border-radius: 20px;
        padding: 2rem;
        display: flex;
        flex-direction: column;
        overflow: hidden;
        border: 1px solid rgba(255, 255, 255, 0.1);
        h3 {
            margin-top: 0;
            text-align: center;
            font-size: 1.5rem;
            margin-bottom: 1.5rem;
            color: #fff;
        }
        .ranking-list {
            overflow-y: auto;
            flex: 1;
            display: flex;
            flex-direction: column;
            gap: 0.8rem;
            padding-right: 0.5rem;
        }
        .rank-item {
            display: flex;
            align-items: center;
            justify-content: space-between;
            background: rgba(255, 255, 255, 0.1);
            padding: 1rem;
            border-radius: 10px;
            font-weight: bold;
            .rank-num {
                opacity: 0.5;
                width: 30px;
            }
            .rank-name {
                flex: 1;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                margin: 0 1rem;
            }
            .rank-score {
                color: var(--accent);
            }
        }
    }
    .question-main {
        flex: 1;
        min-height: 0;
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
        width: 100%;
        overflow: hidden;

        .question-card {
            background: white;
            color: black;
            padding: 1.5rem 2rem;
            text-align: center;
            border-radius: 20px;
            flex-shrink: 0;
            h1 {
                margin: 0;
                font-size: 1.8rem;
                line-height: 1.2;
            }
        }

        .media-area {
            flex: 1;
            min-height: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            overflow: hidden;
            margin-bottom: 1.5rem;
            img {
                max-width: 100%;
                max-height: 100%;
                height: 100%;
                width: auto;
                object-fit: contain;
                border-radius: 15px;
            }
        }
        .options-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 0.8rem;
            margin-top: auto;
            margin-bottom: 1.5rem;
        }

        .option-item {
            padding: 0.8rem 1.5rem;
            border-radius: 12px;
            display: flex;
            align-items: center;
            gap: 1.2rem;
            font-size: 1.2rem;
            font-weight: 700;
            opacity: 0.8;
            position: relative;
            &.correct {
                opacity: 1;
                border: 4px solid white;
                box-shadow: 0 0 30px rgba(255, 255, 255, 0.5);
            }
            &.dimmed {
                filter: grayscale(100%);
                opacity: 0.4;
            }
            .opt-shape {
                width: 25px;
                height: 25px;
                border: 2px solid white;
                flex-shrink: 0;
            }
            .opt-count {
                margin-left: auto;
                background: rgba(0, 0, 0, 0.4);
                padding: 0.5rem 1rem;
                border-radius: 20px;
                font-size: 1.2rem;
                display: flex;
                align-items: center;
                gap: 0.5rem;
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
        width: 100%;
        display: flex;
        justify-content: center;
        padding: 0;
        margin-top: 1.5rem;

        .next-btn {
            width: 100%;
            padding: 1.5rem;
            font-size: 1.2rem;
            border-radius: 10px;
        }
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
