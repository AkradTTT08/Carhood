<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { fade, fly } from "svelte/transition";

    interface Quiz {
        id: string;
        title: string;
        questions: any[];
        room_id: string;
    }

    let quizzes: Quiz[] = [];
    let loading = true;
    let searchQuery = "";

    $: filteredQuizzes = quizzes.filter((q) =>
        q.title.toLowerCase().includes(searchQuery.toLowerCase()),
    );

    async function fetchQuizzes() {
        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch("http://localhost:8081/api/games", {
                headers: { Authorization: `Bearer ${token}` },
            });
            if (res.status === 401) goto("/admin/login");
            const result = await res.json();
            quizzes = Array.isArray(result) ? result : [];
        } catch (e) {
            console.error("Fetch error:", e);
        } finally {
            loading = false;
        }
    }

    onMount(fetchQuizzes);

    function startQuiz(quizId: string) {
        goto(`/admin/host/${quizId}`);
    }
</script>

<div class="start-quiz-container" in:fade>
    <header class="admin-header">
        <div class="header-left">
            <h1>START MISSION</h1>
            <p>SELECT A QUIZ MISSION TO INITIATE</p>
        </div>

        <div class="header-actions">
            <div class="search-wrapper">
                <span class="search-icon">🔍</span>
                <input
                    type="text"
                    placeholder="Search quizzes..."
                    bind:value={searchQuery}
                    class="search-input"
                />
            </div>
        </div>
    </header>

    <div class="quiz-content">
        {#if loading}
            <div class="loading-state">
                <div class="spinner"></div>
                <p>Loading your quizzes...</p>
            </div>
        {:else if filteredQuizzes.length > 0}
            <div class="quiz-list">
                {#each filteredQuizzes as quiz, i}
                    <div
                        class="quiz-row glass-card"
                        in:fly={{ y: 20, delay: i * 50 }}
                    >
                        <div class="quiz-info">
                            <span class="index">{i + 1}</span>
                            <div class="details">
                                <h3>{quiz.title}</h3>
                                <p>{quiz.questions?.length || 0} Questions</p>
                            </div>
                        </div>
                        <button
                            class="btn-3d blue start-btn"
                            on:click={() => startQuiz(quiz.id)}
                        >
                            START GAME 🚀
                        </button>
                    </div>
                {/each}
            </div>
        {:else}
            <div class="empty-state">
                <button
                    class="btn-3d purple"
                    on:click={() => goto("/admin/editor")}
                >
                    Create New Quiz
                </button>
                <p>No quizzes found matching "{searchQuery}"</p>
            </div>
        {/if}
    </div>
</div>

<style lang="scss">
    .start-quiz-container {
        width: 100%;
        box-sizing: border-box;
    }

    .quiz-content {
        padding: 0 2rem 2rem 2rem;
    }

    .admin-header {
        position: sticky;
        top: 0;
        z-index: 100;
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin: -2rem -2rem 3rem -2rem;
        padding: 2rem;
        background: rgba(10, 10, 30, 0.8);
        backdrop-filter: blur(10px);
        border-bottom: 4px solid var(--primary);
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);

        h1 {
            font-size: 1.2rem;
            margin: 0;
            color: var(--accent);
        }
        p {
            color: var(--text-secondary);
            font-size: 0.7rem;
            margin-top: 0.5rem;
            font-family: var(--font-pixel);
        }
    }

    .header-actions {
        display: flex;
        align-items: center;
        gap: 1.5rem;
    }

    .search-wrapper {
        position: relative;
        display: flex;
        align-items: center;
        background: rgba(255, 255, 255, 0.05);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        padding: 0 1rem;
        width: 300px;
        transition: all 0.3s ease;

        &:focus-within {
            background: rgba(255, 255, 255, 0.1);
            border-color: var(--accent);
        }

        .search-icon {
            opacity: 0.5;
            margin-right: 0.5rem;
        }
        .search-input {
            background: none;
            border: none;
            color: white;
            padding: 0.8rem 0;
            font-size: 1rem;
            width: 100%;
            outline: none;
            &::placeholder {
                color: rgba(255, 255, 255, 0.3);
            }
        }
    }

    .quiz-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .quiz-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 2rem;
        border: 4px solid rgba(255, 255, 255, 0.1);
        box-shadow: 0 0 0 4px #000;
        transition: transform 0.2s;

        &:hover {
            transform: scale(1.02);
            border-color: var(--accent);
            background: rgba(255, 255, 255, 0.05);
            box-shadow:
                0 0 0 4px #000,
                0 0 20px rgba(0, 210, 255, 0.3);
        }

        .quiz-info {
            display: flex;
            align-items: center;
            gap: 2rem;

            .index {
                font-size: 1.5rem;
                font-weight: 800;
                opacity: 0.3;
                width: 30px;
            }

            .details {
                h3 {
                    margin: 0;
                    font-size: 1.4rem;
                }
                p {
                    margin: 0.2rem 0 0;
                    color: var(--text-secondary);
                    font-weight: 600;
                }
            }
        }

        .start-btn {
            padding: 0.8rem 2rem;
            font-size: 1rem;
        }
    }

    .loading-state,
    .empty-state {
        text-align: center;
        padding: 5rem 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 1.5rem;
        p {
            font-size: 1.2rem;
            color: var(--text-secondary);
            margin: 0;
        }
    }

    .spinner {
        width: 50px;
        height: 50px;
        border: 4px solid rgba(255, 255, 255, 0.1);
        border-top: 4px solid var(--accent);
        border-radius: 50%;
        margin: 0 auto 1.5rem;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
