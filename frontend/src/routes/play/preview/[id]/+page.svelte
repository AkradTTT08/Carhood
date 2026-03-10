<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { fade, fly, scale } from "svelte/transition";

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
                const res = await fetch(
                    `http://localhost:8081/api/games/${id}`,
                    {
                        headers: {
                            Authorization: `Bearer ${token}`,
                        },
                    },
                );
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
        if (currentQuestion) {
            timeLeft = currentQuestion.time_limit || 20;
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
                    Exit Preview
                </button>
            </div>
        </header>

        <main class="v1-main">
            <div class="canvas-panel">
                <div class="q-card" in:fly={{ y: -30 }}>
                    <h2>{currentQuestion.text}</h2>
                </div>

                <div class="visual-area">
                    <div class="stat-circle timer" class:danger={timeLeft <= 5}>
                        <span class="num">{timeLeft}</span>
                    </div>

                    <div class="media-center">
                        {#if currentQuestion.image_url}
                            <img
                                src={currentQuestion.image_url}
                                alt="Question media"
                                in:scale
                            />
                        {:else}
                            <div class="placeholder" in:scale>🏎️</div>
                        {/if}
                    </div>

                    <div class="stat-circle answers">
                        <span class="num">{answersCount}</span>
                        <small>Answers</small>
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
                    class="nav-link"
                    on:click={prevQuestion}
                    disabled={currentQuestionIndex === 0}>❮ Prev</button
                >
                <div class="progress">
                    <strong>{currentQuestionIndex + 1}</strong> / {quiz
                        .questions.length}
                </div>
                <button
                    class="nav-link"
                    on:click={nextQuestion}
                    disabled={currentQuestionIndex ===
                        quiz.questions.length - 1}>Next ❯</button
                >
            </div>
        </footer>
    {/if}
</div>

<style lang="scss">
    .v1-preview-root {
        position: fixed;
        inset: 0;
        background: radial-gradient(circle at top left, #46178f, #250818);
        color: white;
        z-index: 10000;
        display: flex;
        flex-direction: column;
        font-family: "Inter", sans-serif;
    }

    .v1-header {
        background: rgba(0, 0, 0, 0.3);
        backdrop-filter: blur(10px);
        padding: 1rem 2rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);

        .left {
            display: flex;
            align-items: center;
            gap: 1.5rem;
            h1 {
                font-size: 1.2rem;
                margin: 0;
            }
        }
    }

    .v1-main {
        flex: 1;
        display: flex;
        justify-content: center;
        align-items: flex-start; // Start from top if overflowing
        padding: 1.5rem;
        overflow-y: auto; // Allow scroll for the content

        // Hide scrollbar for clean look but allow scrolling
        &::-webkit-scrollbar {
            width: 0;
            background: transparent;
        }
        scrollbar-width: none;
    }

    .canvas-panel {
        width: 100%;
        max-width: 1100px;
        display: flex;
        flex-direction: column;
        gap: 1.5rem; // Reduced gap
        padding-bottom: 2rem; // Extra bottom space for scrolling
    }

    .q-card {
        background: white;
        padding: 2rem;
        border-radius: 12px;
        box-shadow: 0 10px 0 rgba(0, 0, 0, 0.2);
        h2 {
            color: #333;
            text-align: center;
            margin: 0;
            font-size: 2.4rem;
            font-weight: 800;
        }
    }

    .visual-area {
        height: 280px;
        display: flex;
        align-items: center;
        justify-content: space-around;
        position: relative;
    }

    .stat-circle {
        width: 100px;
        height: 100px;
        border-radius: 50%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);

        &.timer {
            background: #e21b3c;
            .num {
                font-size: 2.5rem;
                font-weight: 900;
            }
            &.danger {
                animation: pulse 1s infinite;
            }
        }

        &.answers {
            background: #1368ce;
            .num {
                font-size: 2.2rem;
                font-weight: 900;
                line-height: 1;
            }
            small {
                font-size: 0.6rem;
                text-transform: uppercase;
                font-weight: 700;
                opacity: 0.8;
            }
        }
    }

    @keyframes pulse {
        0%,
        100% {
            transform: scale(1);
        }
        50% {
            transform: scale(1.1);
        }
    }

    .media-center {
        flex: 1;
        max-width: 500px;
        height: 100%;
        margin: 0 2rem;
        background: white;
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        overflow: hidden;
        border: 2px solid rgba(255, 255, 255, 0.1);

        img {
            width: 100%;
            height: 100%;
            object-fit: contain;
        }
        .placeholder {
            font-size: 8rem;
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
        padding: 1.5rem 2rem;
        border-radius: 12px;
        border: none;
        cursor: pointer;
        transition: all 0.2s;
        box-shadow: 0 8px 0 rgba(0, 0, 0, 0.2);
        top: 0;

        &.red {
            background: #e21b3c;
        }
        &.blue {
            background: #1368ce;
        }
        &.yellow {
            background: #d89e00;
        }
        &.green {
            background: #26890c;
        }

        &:hover:not(:disabled) {
            top: -4px;
            box-shadow: 0 12px 0 rgba(0, 0, 0, 0.2);
        }

        &:active:not(:disabled) {
            top: 4px;
            box-shadow: 0 4px 0 rgba(0, 0, 0, 0.2);
        }

        &.selected {
            border: 4px solid white;
        }

        &.dimmed {
            opacity: 0.3;
        }

        &.correct {
            transform: scale(1.05);
            z-index: 10;
        }

        .shape {
            width: 35px;
            height: 35px;
            background: white;
            flex-shrink: 0;
            &.triangle {
                clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
            }
            &.diamond {
                transform: rotate(45deg);
                width: 28px;
                height: 28px;
            }
            &.circle {
                border-radius: 50%;
            }
        }

        .text {
            color: white;
            font-size: 1.5rem;
            font-weight: 700;
            margin-left: 2rem;
            text-align: left;
            flex: 1;
        }

        .check {
            position: absolute;
            right: 1.5rem;
            background: white;
            color: #26890c;
            width: 30px;
            height: 30px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 900;
        }
    }

    .v1-footer {
        height: 80px;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        border-top: 1px solid rgba(255, 255, 255, 0.1);

        .nav-btns {
            display: flex;
            align-items: center;
            gap: 3rem;

            .nav-link {
                background: none;
                border: none;
                color: rgba(255, 255, 255, 0.9); // Increased contrast
                font-size: 1.1rem;
                font-weight: 600;
                cursor: pointer;
                transition: color 0.2s;
                &:hover:not(:disabled) {
                    color: white;
                }
                &:disabled {
                    opacity: 0.2;
                    cursor: not-allowed;
                }
            }

            .progress {
                font-size: 1.1rem;
                strong {
                    font-size: 1.5rem;
                }
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
