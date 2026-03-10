<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { fade, fly } from "svelte/transition";

    interface PlayerResult {
        username: string;
        score: number;
    }

    interface QuizHistory {
        id: string;
        quiz_id: string;
        quiz_title: string;
        room_id: string;
        played_at: string;
        players_count: number;
        leaderboard: PlayerResult[];
    }

    let history: QuizHistory[] = [];
    let loading = true;
    let searchQuery = "";

    $: filteredHistory = history.filter(
        (h) =>
            h.quiz_title.toLowerCase().includes(searchQuery.toLowerCase()) ||
            h.room_id.toLowerCase().includes(searchQuery.toLowerCase()),
    );

    async function fetchHistory() {
        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch("http://localhost:8081/api/history", {
                headers: { Authorization: `Bearer ${token}` },
            });
            if (res.status === 401) goto("/admin/login");
            const result = await res.json();
            history = Array.isArray(result) ? result : [];
        } catch (e) {
            console.error("Fetch error:", e);
        } finally {
            loading = false;
        }
    }

    function formatDate(dateStr: string) {
        const date = new Date(dateStr);
        return date.toLocaleString("th-TH", {
            year: "numeric",
            month: "long",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    onMount(fetchHistory);
</script>

<div class="history-container" in:fade>
    <header class="admin-header">
        <div class="header-left">
            <h1>Quiz History</h1>
            <p>เรียกดูประวัติการเล่นย้อนหลัง</p>
        </div>

        <div class="header-actions">
            <div class="search-wrapper">
                <span class="search-icon">🔍</span>
                <input
                    type="text"
                    placeholder="Search history..."
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
                <p>Loading history...</p>
            </div>
        {:else if filteredHistory.length > 0}
            <div class="history-list">
                {#each filteredHistory as item, i}
                    <div
                        class="history-row glass-card"
                        in:fly={{ y: 20, delay: i * 50 }}
                        on:click={() => goto(`/admin/history/${item.id}`)}
                    >
                        <div class="history-info">
                            <div class="details">
                                <h3>{item.quiz_title}</h3>
                                <p class="timestamp">
                                    📅 {formatDate(item.played_at)}
                                </p>
                            </div>
                        </div>
                        <div class="meta">
                            <div class="stat">
                                <span class="label">PIN</span>
                                <span class="value">{item.room_id}</span>
                            </div>
                            <div class="stat">
                                <span class="label">Players</span>
                                <span class="value">{item.players_count}</span>
                            </div>
                            <span class="arrow">➡️</span>
                        </div>
                    </div>
                {/each}
            </div>
        {:else}
            <div class="empty-state">
                <p>
                    {searchQuery
                        ? `No results for "${searchQuery}"`
                        : "No game history yet."}
                </p>
            </div>
        {/if}
    </div>
</div>

<style lang="scss">
    .history-container {
        width: 100%;
        box-sizing: border-box;
    }

    .admin-header {
        position: sticky;
        top: 0;
        z-index: 100;
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin: -2rem -2rem 3rem -2rem;
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

    .quiz-content {
        padding: 0 2rem 2rem 2rem;
    }

    .history-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .history-row {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.5rem 2rem;
        border-radius: 20px;
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
            transform: scale(1.01);
            background: rgba(255, 255, 255, 0.1);
            .arrow {
                transform: translateX(5px);
                opacity: 1;
            }
        }

        .details {
            h3 {
                margin: 0;
                font-size: 1.4rem;
                color: white;
            }
            .timestamp {
                margin: 0.3rem 0 0;
                color: var(--text-secondary);
                font-size: 0.9rem;
            }
        }

        .meta {
            display: flex;
            align-items: center;
            gap: 2.5rem;

            .stat {
                text-align: right;
                .label {
                    display: block;
                    font-size: 0.7rem;
                    color: var(--text-secondary);
                    text-transform: uppercase;
                    font-weight: 700;
                    letter-spacing: 1px;
                }
                .value {
                    font-size: 1.2rem;
                    font-weight: 800;
                    color: var(--accent);
                }
            }

            .arrow {
                font-size: 1.2rem;
                opacity: 0.3;
                transition: all 0.2s;
            }
        }
    }

    .loading-state,
    .empty-state {
        text-align: center;
        padding: 5rem 0;
        p {
            font-size: 1.2rem;
            color: var(--text-secondary);
            margin-bottom: 2rem;
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
