<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { fade, fly, scale } from "svelte/transition";
    import { API_URL, resolveImageUrl } from "$lib/api";

    interface Question {
        text: string;
        image_url?: string;
        options: string[];
        correct_answer: number;
        time_limit?: number;
    }

    interface Quiz {
        title: string;
        questions: Question[];
    }

    let quiz: Quiz | null = null;
    let loading = true;
    let error: string | null = null;
    let currentQuestionIndex = 0;
    let timeLeft = 0;
    let timerInterval: any;
    let showAnswer = false;
    let selectedAnswer: number | null = null;
    let answersCount = 0;

    $: currentQuestion = quiz?.questions
        ? quiz.questions[currentQuestionIndex]
        : null;

    onMount(() => {
        const fetchQuiz = async () => {
            const id = $page.params.id;
            const token = localStorage.getItem("admin_token");
            try {
                const res = await fetch(`${API_URL}/api/games/${id}`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });
                if (!res.ok) throw new Error("Failed to fetch quiz");
                const data = (await res.json()) as Quiz;
                quiz = data;
                resetTimer();
            } catch (err: any) {
                error = err.message;
            } finally {
                loading = false;
            }
        };

        fetchQuiz();

        return () => clearInterval(timerInterval);
    });

    function resetTimer() {
        clearInterval(timerInterval);
        const q = quiz?.questions?.[currentQuestionIndex];
        if (q) {
            timeLeft = q.time_limit || 20;
            showAnswer = false;
            selectedAnswer = null;
            answersCount = Math.floor(Math.random() * 10) + 5;
            timerInterval = setInterval(() => {
                if (timeLeft > 0) {
                    timeLeft--;
                } else {
                    clearInterval(timerInterval);
                    showAnswer = true;
                }
            }, 1000);
        }
    }

    function nextQuestion() {
        if (quiz && currentQuestionIndex < quiz.questions.length - 1) {
            currentQuestionIndex++;
            resetTimer();
        }
    }

    function prevQuestion() {
        if (currentQuestionIndex > 0) {
            currentQuestionIndex--;
            resetTimer();
        }
    }

    function selectAnswer(i: number) {
        if (!showAnswer) {
            selectedAnswer = i;
            answersCount++;
        }
    }

    const shapes = ["triangle", "diamond", "circle", "square"];
    const colors = ["red", "blue", "yellow", "green"];
</script>

<div class="v1-preview-root">
    {#if loading}
        <div class="loading-state">
            <div class="spinner"></div>
            <p>Loading Preview...</p>
        </div>
    {:else if error}
        <div class="error-state">
            <p>{error}</p>
            <button class="btn-3d purple" on:click={() => goto("/admin")}
                >Back to Dashboard</button
            >
        </div>
    {:else if quiz && currentQuestion}
        <header class="v1-header">
            <div class="left">
                <span class="badge">Preview Mode</span>
                <h1>{quiz.title}</h1>
            </div>
            <div class="right">
                <button
                    class="btn-3d red"
                    on:click={() => (window.location.href = "/admin")}
                >
                    EXIT PREVIEW
                </button>
            </div>
        </header>

        <main class="v1-main">
            <div class="canvas-panel">
                <div class="q-card" in:fly={{ y: -30 }}>
                    <h2>{currentQuestion.text}</h2>
                </div>

                <div
                    class="visual-area"
                    class:no-media={!currentQuestion.image_url}
                >
                    <div class="stat-circle timer" class:danger={timeLeft <= 5}>
                        <span class="num">{timeLeft}</span>
                    </div>

                    {#if currentQuestion.image_url}
                        <div class="media-center">
                            <img
                                src={resolveImageUrl(currentQuestion.image_url)}
                                alt="Question media"
                                in:scale
                            />
                        </div>
                    {/if}

                    <div class="stat-circle answers">
                        <span class="num">{answersCount}</span>
                        <small>ANSWERS</small>
                    </div>
                </div>

                <div class="v1-answer-grid">
                    {#each currentQuestion.options as option, i}
                        <button
                            class="ans-pill {colors[i]}"
                            class:selected={selectedAnswer === i}
                            class:correct={showAnswer &&
                                i === currentQuestion.correct_answer}
                            class:dimmed={showAnswer &&
                                i !== currentQuestion.correct_answer &&
                                selectedAnswer !== i}
                            on:click={() => selectAnswer(i)}
                            disabled={showAnswer}
                            in:fly={{ y: 20, delay: i * 100 }}
                        >
                            <div class="shape {shapes[i]}"></div>
                            <span class="text">{option}</span>
                            {#if showAnswer && i === currentQuestion.correct_answer}
                                <div class="check">✔</div>
                            {/if}
                        </button>
                    {/each}
                </div>
            </div>
        </main>

        <footer class="v1-footer">
            <div class="nav-btns">
                <button
                    class="btn-3d purple small"
                    on:click={prevQuestion}
                    disabled={currentQuestionIndex === 0}>PREV</button
                >
                <div class="progress">
                    <strong>{currentQuestionIndex + 1}</strong> / {quiz
                        .questions.length}
                </div>
                <button
                    class="btn-3d purple small"
                    on:click={nextQuestion}
                    disabled={currentQuestionIndex ===
                        quiz.questions.length - 1}>NEXT</button
                >
            </div>
        </footer>
    {/if}
</div>

<style lang="scss">
    .v1-preview-root {
        position: fixed;
        inset: 0;
        background: radial-gradient(circle at top left, #29084a, #050510);
        color: white;
        z-index: 10000;
        display: flex;
        flex-direction: column;
        font-family: var(--font-main);
    }

    .v1-header {
        background: rgba(0, 0, 0, 0.5);
        backdrop-filter: blur(10px);
        padding: 1rem 2rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        border-bottom: 4px solid rgba(255, 255, 255, 0.1);

        .left {
            display: flex;
            align-items: center;
            gap: 2rem;

            .badge {
                font-family: var(--font-pixel);
                font-size: 0.6rem;
                background: var(--accent);
                color: #000;
                padding: 4px 8px;
                border-radius: 4px;
                animation: pulse 2s infinite;
            }

            h1 {
                font-family: var(--font-pixel);
                font-size: 1rem;
                margin: 0;
                color: white;
                text-shadow: 0 0 10px rgba(0, 210, 255, 0.5);
            }
        }
    }

    .v1-main {
        flex: 1;
        display: flex;
        justify-content: center;
        align-items: flex-start;
        padding: 2rem;
        overflow-y: auto;

        &::-webkit-scrollbar {
            width: 8px;
        }
        &::-webkit-scrollbar-track {
            background: rgba(0, 0, 0, 0.2);
        }
        &::-webkit-scrollbar-thumb {
            background: var(--primary);
            border-radius: 4px;
        }
    }

    .canvas-panel {
        width: 100%;
        max-width: 1000px;
        display: flex;
        flex-direction: column;
        gap: 2rem;
        padding-bottom: 4rem;
    }

    .q-card {
        background: white;
        padding: 2.5rem;
        border-radius: 8px;
        box-shadow: 0 10px 0 rgba(0, 0, 0, 0.3);
        border: 4px solid #fff;

        h2 {
            font-family: var(--font-pixel);
            color: #333;
            text-align: center;
            margin: 0;
            font-size: 1.5rem;
            line-height: 1.6;
        }
    }

    .visual-area {
        height: 300px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        position: relative;
        transition: all 0.3s ease;

        &.no-media {
            justify-content: center;
            gap: 15rem;
        }
    }

    .stat-circle {
        width: 120px;
        height: 120px;
        border-radius: 50%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        box-shadow: 0 0 30px rgba(0, 0, 0, 0.5);
        border: 4px solid #fff;
        font-family: var(--font-pixel);

        &.timer {
            background: var(--secondary);
            .num {
                font-size: 2.2rem;
                font-weight: 900;
            }
            &.danger {
                animation: emergencyPulse 0.5s infinite;
                background: #ff0000;
            }
        }

        &.answers {
            background: #1368ce;
            .num {
                font-size: 1.8rem;
                font-weight: 900;
            }
            small {
                font-size: 0.5rem;
                margin-top: 4px;
                opacity: 0.8;
            }
        }
    }

    @keyframes emergencyPulse {
        0%,
        100% {
            transform: scale(1);
            box-shadow: 0 0 30px rgba(255, 0, 0, 0.5);
        }
        50% {
            transform: scale(1.15);
            box-shadow: 0 0 50px rgba(255, 0, 0, 0.8);
        }
    }

    .media-center {
        flex: 1;
        max-width: 500px;
        height: 100%;
        margin: 0 2rem;
        background: rgba(255, 255, 255, 0.05);
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        overflow: hidden;
        border: 4px solid rgba(255, 255, 255, 0.1);
        backdrop-filter: blur(5px);

        img {
            width: 100%;
            height: 100%;
            object-fit: contain;
            image-rendering: auto;
        }
        .placeholder {
            font-size: 6rem;
            filter: drop-shadow(0 0 20px rgba(255, 255, 255, 0.3));
        }
    }

    .v1-answer-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.5rem;
    }

    .ans-pill {
        position: relative;
        display: flex;
        align-items: center;
        padding: 1.8rem 2.5rem;
        border-radius: 8px;
        border: 4px solid #fff;
        cursor: pointer;
        transition: all 0.1s;
        top: 0;

        &.red {
            background: #e21b3c;
            box-shadow: 0 8px 0 #8b0000;
        }
        &.blue {
            background: #1368ce;
            box-shadow: 0 8px 0 #0e4e9a;
        }
        &.yellow {
            background: #d89e00;
            box-shadow: 0 8px 0 #8b8b00;
        }
        &.green {
            background: #26890c;
            box-shadow: 0 8px 0 #1b6208;
        }

        &:hover:not(:disabled) {
            top: -4px;
            filter: brightness(1.1);
        }

        &:active:not(:disabled) {
            top: 4px;
            box-shadow: 0 2px 0 rgba(0, 0, 0, 0.3);
        }

        &.selected {
            outline: 6px solid var(--accent);
            outline-offset: 4px;
        }

        &.dimmed {
            opacity: 0.2;
            filter: grayscale(0.8);
        }

        &.correct {
            transform: scale(1.05);
            z-index: 10;
            box-shadow: 0 0 30px var(--success);
        }

        .shape {
            width: 40px;
            height: 40px;
            background: white;
            flex-shrink: 0;
            filter: drop-shadow(0 4px 0 rgba(0, 0, 0, 0.2));

            &.triangle {
                clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
            }
            &.diamond {
                transform: rotate(45deg);
                width: 32px;
                height: 32px;
            }
            &.circle {
                border-radius: 50%;
            }
        }

        .text {
            color: white;
            font-family: var(--font-pixel);
            font-size: 0.8rem;
            margin-left: 2rem;
            text-align: left;
            flex: 1;
            line-height: 1.4;
        }

        .check {
            position: absolute;
            right: 1.5rem;
            background: white;
            color: #26890c;
            width: 40px;
            height: 40px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 1.5rem;
            font-weight: 900;
            box-shadow: 0 4px 0 rgba(0, 0, 0, 0.2);
            animation: popIn 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
        }
    }

    @keyframes popIn {
        from {
            transform: scale(0);
        }
        to {
            transform: scale(1);
        }
    }

    .v1-footer {
        height: 100px;
        background: rgba(0, 0, 0, 0.8);
        display: flex;
        align-items: center;
        justify-content: center;
        border-top: 4px solid var(--primary);

        .nav-btns {
            display: flex;
            align-items: center;
            gap: 4rem;

            .progress {
                font-family: var(--font-pixel);
                font-size: 0.8rem;
                color: var(--accent);
                strong {
                    font-size: 1.4rem;
                    color: white;
                }
            }

            .btn-3d.small {
                padding: 0.6rem 1.2rem;
                font-size: 0.6rem;
            }
        }
    }

    .loading-state,
    .error-state {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 2rem;
    }
    .spinner {
        width: 60px;
        height: 60px;
        border: 4px solid rgba(255, 255, 255, 0.1);
        border-top-color: white;
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }
    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
