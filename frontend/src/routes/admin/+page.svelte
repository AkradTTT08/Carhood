<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import Modal from "$lib/components/Modal.svelte";
    import { API_URL } from "$lib/api";

    interface Question {
        text: string;
        options: string[];
        correct_answer: number;
        time_limit?: number;
    }

    interface Quiz {
        id?: string;
        title: string;
        questions: Question[];
        room_id: string;
    }

    let showDeleteModal = false;
    let quizToDelete: string | null = null;

    let quizzes: Quiz[] = [];
    let searchQuery = "";

    $: filteredQuizzes = quizzes.filter((q) =>
        q.title.toLowerCase().includes(searchQuery.toLowerCase()),
    );

    async function fetchQuizzes() {
        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch(`${API_URL}/api/games`, {
                headers: { Authorization: `Bearer ${token}` },
            });
            if (res.status === 401) goto("/admin/login");
            const result = await res.json();
            quizzes = Array.isArray(result) ? result : [];
        } catch (e) {
            console.error("Fetch error:", e);
        }
    }

    async function duplicateQuiz(id: string) {
        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch(
                `${API_URL}/api/games/${id}/duplicate`,
                {
                    method: "POST",
                    headers: { Authorization: `Bearer ${token}` },
                },
            );
            if (res.ok) await fetchQuizzes();
        } catch (e) {
            console.error("Duplicate error:", e);
        }
    }

    onMount(fetchQuizzes);
    function deleteQuiz(id: string) {
        quizToDelete = id;
        showDeleteModal = true;
    }

    async function confirmDelete() {
        if (!quizToDelete) return;
        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch(
                `http://localhost:8081/api/games/${quizToDelete}`,
                {
                    method: "DELETE",
                    headers: { Authorization: `Bearer ${token}` },
                },
            );
            if (res.ok) {
                quizzes = quizzes.filter((q) => q.id !== quizToDelete);
                showDeleteModal = false;
                quizToDelete = null;
            }
        } catch (e) {
            console.error(e);
        }
    }
</script>

<div class="page-bg">
    <img src="/image/BGUSER.svg" alt="Background" class="bg-image" />
</div>

<div class="master-quiz-container">
    <Modal
        bind:show={showDeleteModal}
        title="ลบ Quiz (Delete Quiz)"
        message="คุณแน่ใจหรือไม่ว่าต้องการลบ Quiz นี้? การกระทำนี้ไม่สามารถย้อนกลับได้"
        confirmText="ลบทันที (Delete Now)"
        cancelText="ยกเลิก (Cancel)"
        type="danger"
        on:confirm={confirmDelete}
        on:cancel={() => (showDeleteModal = false)}
    />

    <header class="admin-header">
        <div class="header-left">
            <h1>MASTER QUIZ</h1>
            <p>COMMAND CENTER: MANAGE YOUR MISSIONS</p>
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

            <button on:click={() => goto("/admin/editor")} class="btn-3d blue">
                + NEW QUIZ
            </button>
        </div>
    </header>

    <div class="quiz-content">
        {#if filteredQuizzes.length > 0}
            <div class="quiz-grid">
                {#each filteredQuizzes as quiz}
                    <div class="glass-card quiz-card">
                        <div class="quiz-info">
                            <h3>{quiz.title}</h3>
                            <span class="badge"
                                >{quiz.questions?.length || 0} Qs</span
                            >
                        </div>
                        <div class="card-actions">
                            <button
                                class="btn-3d teal btn-preview"
                                on:click={() =>
                                    (window.location.href = `/play/preview/${
                                        quiz.id || ""
                                    }`)}>Preview</button
                            >
                            <a
                                href="/admin/editor?id={quiz.id || ''}"
                                class="btn-3d purple btn-edit">Edit</a
                            >
                            <button
                                class="btn-3d red btn-icon"
                                on:click={() => quiz.id && deleteQuiz(quiz.id)}
                            >
                                🗑️
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
        {:else}
            <div class="empty-state">
                <p>No quizzes found. Create your first one!</p>
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

    .master-quiz-container {
        width: 100%;
        box-sizing: border-box;
        position: relative;
        z-index: 1;
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
            line-height: 1.6;
        }

        p {
            color: var(--text-secondary);
            font-size: 0.7rem;
            margin-top: 0.8rem;
            font-family: var(--font-pixel);
            line-height: 1.6;
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
            box-shadow: 0 0 20px rgba(var(--accent-rgb), 0.2);
        }

        .search-icon {
            font-size: 1.2rem;
            margin-right: 0.5rem;
            opacity: 0.5;
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

    .quiz-grid {
        display: grid;
        grid-template-columns: 1fr; // Full width cards
        gap: 1.5rem;
    }

    .quiz-card {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 2rem;
        padding: 2rem;
        border: 4px solid rgba(255, 255, 255, 0.1);
        box-shadow: 0 0 0 4px #000;
        background: var(--bg-card);
        transition: transform 0.2s;

        &:hover {
            transform: scale(1.02);
            border-color: var(--accent);
            box-shadow:
                0 0 0 4px #000,
                0 0 20px rgba(0, 210, 255, 0.3);
        }
    }

    .quiz-info {
        flex: 1; // Take remaining space
        display: flex;
        justify-content: flex-start;
        align-items: center;
        gap: 1.5rem;

        h3 {
            font-size: 1.6rem;
            margin: 0;
            font-weight: 800;
        }
    }

    .badge {
        background: var(--primary);
        padding: 0.5rem 1rem;
        border-radius: 12px;
        font-size: 1rem;
        font-weight: 800;
        box-shadow: 0 4px 0 rgba(0, 0, 0, 0.2);
        white-space: nowrap; // Ensure content is on one line
        display: flex;
        align-items: center;
        justify-content: center;
        height: fit-content;
    }

    .card-actions {
        display: flex;
        gap: 1rem;
        margin-top: 0; // Remove top margin for row layout

        .btn-preview {
            width: 140px;
        }

        .btn-edit {
            width: 100px;
            text-align: center;
            text-decoration: none;
        }
    }

    .btn-edit {
        background: rgba(255, 255, 255, 0.1);
    }

    .btn-icon {
        width: 45px;
        padding: 0.5rem;
        font-size: 1.2rem;
    }
</style>
