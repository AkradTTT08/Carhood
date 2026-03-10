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
            <h1>Start Game</h1>
            <p>เลือกชุดคำถามที่คุณต้องการเริ่มเล่น</p>
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
                <p>No quizzes found matching "{searchQuery}"</p>
                <button
                    class="btn-3d purple"
                    on:click={() => goto("/admin/editor")}
                >
                    Create New Quiz
                </button>
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
        margin: -2rem -2rem 3rem -2rem; // Offset parent padding
        padding: 1.5rem 2rem;
        background: rgba(70, 23, 143, 0.4);
        backdrop-filter: blur(15px);
        -webkit-backdrop-filter: blur(15px);
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);

        h1 {
            font-size: 2.2rem;
            margin: 0;
        }
        p {
            color: var(--text-secondary);
            font-size: 0.95rem;
            margin: 0;
            opacity: 0.8;
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
        padding: 1.5rem 2rem;
        border-radius: 20px;
        transition: transform 0.2s;

        &:hover {
            transform: scale(1.01);
            background: rgba(255, 255, 255, 0.1);
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
        p {
            font-size: 1.2rem;
            color: var(--text-secondary);
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
