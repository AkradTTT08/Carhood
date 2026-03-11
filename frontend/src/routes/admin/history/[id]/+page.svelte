<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { fade, fly, scale } from "svelte/transition";
    import { API_URL } from "$lib/api";

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

    let record: QuizHistory | null = null;
    let loading = true;
    let error: string | null = null;

    async function fetchDetail() {
        const id = $page.params.id;
        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch(`${API_URL}/api/history/${id}`, {
                headers: { Authorization: `Bearer ${token}` },
            });
            if (res.status === 401) goto("/admin/login");
            if (!res.ok) throw new Error("Could not find record");
            record = await res.json();
        } catch (e: any) {
            error = e.message;
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

    onMount(fetchDetail);
</script>

<div class="page-bg">
    <img src="/image/BGUSER.svg" alt="Background" class="bg-image" />
</div>

<div class="detail-container" in:fade>
    <header class="admin-header">
        <div class="header-left">
            <button class="back-link" on:click={() => goto("/admin/history")}>
                ⬅️ Back to History
            </button>
            <h1>Game Details</h1>
            <p>SESSION RESULTS & LEADERBOARD</p>
        </div>
    </header>

    <div class="detail-content">
        {#if loading}
            <div class="center">
                <div class="spinner"></div>
                <p>Loading session details...</p>
            </div>
        {:else if error}
            <div class="center">
                <p class="error">{error}</p>
                <button
                    class="btn-3d purple"
                    on:click={() => goto("/admin/history")}>Go Back</button
                >
            </div>
        {:else if record}
            <div class="summary-card glass-card" in:scale>
                <div class="title-section">
                    <span class="label">Quiz Title</span>
                    <h1>{record.quiz_title}</h1>
                </div>

                <div class="stats-grid">
                    <div class="stat-box">
                        <span class="label">Played On</span>
                        <h3>{formatDate(record.played_at)}</h3>
                    </div>
                    <div class="stat-box">
                        <span class="label">Room PIN</span>
                        <h3 class="pin">{record.room_id}</h3>
                    </div>
                    <div class="stat-box">
                        <span class="label">Total Players</span>
                        <h3>{record.players_count}</h3>
                    </div>
                </div>
            </div>

            <div class="leaderboard-section">
                <h2>Final Leaderboard 🏆</h2>
                <div class="leaderboard-list">
                    {#each record.leaderboard as player, i}
                        <div
                            class="player-rank glass-card"
                            in:fly={{ x: -20, delay: i * 50 }}
                        >
                            <div class="rank-info">
                                <span class="rank-num">#{i + 1}</span>
                                <span class="player-name"
                                    >{player.username}</span
                                >
                            </div>
                            <div class="score">
                                {player.score} <span class="pts">pts</span>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
    </div>
</div>

<style lang="scss">
    .page-bg {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 0;
        background: #000;

        .bg-image {
            width: 100%;
            height: 100%;
            object-fit: cover;
            image-rendering: auto;
        }
    }

    .detail-container {
        width: 100%;
        box-sizing: border-box;
        position: relative;
        z-index: 1;
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

        .header-left {
            display: flex;
            align-items: center;
            gap: 2rem;
        }
        .back-link {
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.1);
            color: white;
            padding: 0.8rem 1.2rem;
            border-radius: 12px;
            font-weight: 700;
            cursor: pointer;
            font-size: 0.9rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
            transition: all 0.2s;
            &:hover {
                background: rgba(255, 255, 255, 0.1);
                border-color: var(--accent);
                transform: translateX(-5px);
            }
        }
        h1 {
            font-size: 1.2rem;
            margin: 0;
            color: var(--accent);
            line-height: 1.6;
        }
    }

    .detail-content {
        padding: 0 2rem 2rem 2rem;
        max-width: 900px;
        margin: 0 auto;
    }

    .summary-card {
        padding: 3rem;
        border-radius: 30px;
        margin-bottom: 3rem;
        text-align: center;

        .label {
            display: block;
            font-size: 0.9rem;
            color: var(--text-secondary);
            text-transform: uppercase;
            font-weight: 700;
            letter-spacing: 2px;
            margin-bottom: 0.5rem;
        }
        h1 {
            font-size: 3rem;
            margin: 0 0 2rem;
            color: white;
        }

        .stats-grid {
            display: grid;
            grid-template-columns: 1fr 1fr 1fr;
            gap: 2rem;
            border-top: 1px solid rgba(255, 255, 255, 0.1);
            padding-top: 2rem;
            .stat-box {
                h3 {
                    margin: 0;
                    font-size: 1.4rem;
                    color: var(--accent);
                }
                .pin {
                    font-family: monospace;
                    letter-spacing: 2px;
                }
            }
        }
    }

    .leaderboard-section {
        h2 {
            margin-bottom: 1.5rem;
            font-size: 1.8rem;
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
        padding: 1.2rem 2rem;
        border-radius: 15px;

        .rank-info {
            display: flex;
            align-items: center;
            gap: 1.5rem;
            .rank-num {
                font-size: 1.5rem;
                font-weight: 900;
                opacity: 0.3;
                width: 40px;
            }
            .player-name {
                font-size: 1.4rem;
                font-weight: 700;
            }
        }

        .score {
            font-size: 1.8rem;
            font-weight: 900;
            color: var(--accent);
            .pts {
                font-size: 1rem;
                opacity: 0.5;
            }
        }
    }

    .center {
        text-align: center;
        padding: 5rem 0;
        .error {
            color: #ff4d4d;
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

    @media (max-width: 768px) {
        .summary-card {
            padding: 2rem;
            .stats-grid {
                grid-template-columns: 1fr;
            }
        }
    }
</style>
