<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import Modal from "$lib/components/Modal.svelte";

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

    let showCreateModal = false;
    let newQuiz: Quiz = {
        title: "",
        questions: [],
        room_id: Math.floor(100000 + Math.random() * 900000).toString(),
    };

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
        }
    }

    function addQuestion() {
        newQuiz.questions = [
            ...newQuiz.questions,
            {
                text: "",
                options: ["", "", "", ""],
                correct_answer: 0,
                time_limit: 20,
            },
        ];
    }

    async function saveQuiz() {
        if (!newQuiz.title || newQuiz.questions.length === 0) return;

        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch("http://localhost:8081/api/games", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(newQuiz),
            });
            if (res.ok) {
                showCreateModal = false;
                newQuiz = {
                    title: "",
                    questions: [],
                    room_id: Math.floor(
                        100000 + Math.random() * 900000,
                    ).toString(),
                };
                await fetchQuizzes();
            }
        } catch (e) {
            console.error("Save error:", e);
        }
    }

    async function duplicateQuiz(id: string) {
        const token = localStorage.getItem("admin_token");
        try {
            const res = await fetch(
                `http://localhost:8081/api/games/${id}/duplicate`,
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

            <button
                on:click={() => (showCreateModal = true)}
                class="btn-3d blue"
            >
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

{#if showCreateModal}
    <div class="modal-overlay">
        <div class="glass-card modal-content">
            <header class="modal-header">
                <h2>Create New Quiz</h2>
                <button
                    class="btn-close"
                    on:click={() => (showCreateModal = false)}>&times;</button
                >
            </header>

            <section class="modal-body">
                <div class="input-group">
                    <label>Quiz Title</label>
                    <input
                        type="text"
                        placeholder="e.g. Science Trivia"
                        bind:value={newQuiz.title}
                        class="input-field"
                    />
                </div>

                <div class="questions-section">
                    <h3>Questions ({newQuiz.questions.length})</h3>
                    <div class="questions-list">
                        {#each newQuiz.questions as question, i}
                            <div class="question-builder">
                                <div class="q-header">
                                    <span>Question {i + 1}</span>
                                    <button
                                        class="btn-remove"
                                        on:click={() =>
                                            (newQuiz.questions =
                                                newQuiz.questions.filter(
                                                    (_, idx) => idx !== i,
                                                ))}>&times;</button
                                    >
                                </div>
                                <input
                                    type="text"
                                    placeholder="Your question here..."
                                    bind:value={question.text}
                                    class="input-field"
                                />

                                <div class="options-inputs">
                                    {#each question.options as option, optIdx}
                                        <div class="option-row">
                                            <label class="radio-label">
                                                <input
                                                    type="radio"
                                                    name="correct-{i}"
                                                    value={optIdx}
                                                    bind:group={
                                                        question.correct_answer
                                                    }
                                                />
                                            </label>
                                            <input
                                                type="text"
                                                placeholder="Option {optIdx +
                                                    1}"
                                                bind:value={
                                                    question.options[optIdx]
                                                }
                                                class="input-field"
                                            />
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        {/each}
                    </div>
                    <button class="btn-secondary" on:click={addQuestion}
                        >+ Add Question</button
                    >
                </div>
            </section>

            <footer class="modal-footer">
                <button
                    class="btn-ghost"
                    on:click={() => (showCreateModal = false)}>Cancel</button
                >
                <button class="btn-primary" on:click={saveQuiz}
                    >Save Quiz</button
                >
            </footer>
        </div>
    </div>
{/if}

<style lang="scss">
    .master-quiz-container {
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

    .modal-overlay {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.85);
        backdrop-filter: blur(5px);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1000;
        padding: 1rem;
    }

    .modal-content {
        width: 100%;
        max-width: 900px;
        height: 90vh;
        display: flex;
        flex-direction: column;
        padding: 0;
        overflow: hidden;
    }

    .modal-header,
    .modal-footer {
        padding: 1.5rem 2rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        background: rgba(255, 255, 255, 0.05);
    }

    .modal-body {
        flex: 1;
        overflow-y: auto;
        padding: 2rem;
        display: flex;
        flex-direction: column;
        gap: 2rem;
    }

    .input-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        label {
            font-weight: 600;
            color: var(--text-secondary);
        }
    }

    .question-builder {
        background: rgba(255, 255, 255, 0.03);
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 12px;
        padding: 1.5rem;
        margin-bottom: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .q-header {
        display: flex;
        justify-content: space-between;
        font-weight: 700;
        color: var(--accent);
    }

    .options-inputs {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }

    .option-row {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        input[type="radio"] {
            width: 20px;
            height: 20px;
            accent-color: var(--success);
        }
    }

    .btn-remove {
        background: none;
        border: none;
        color: var(--secondary);
        font-size: 1.5rem;
        cursor: pointer;
    }

    .btn-secondary {
        width: 100%;
        padding: 1rem;
        background: rgba(255, 255, 255, 0.05);
        border: 2px dashed rgba(255, 255, 255, 0.2);
        color: white;
        border-radius: 12px;
        cursor: pointer;
        font-weight: 600;
        &:hover {
            background: rgba(255, 255, 255, 0.1);
        }
    }

    .btn-close {
        background: none;
        border: none;
        color: white;
        font-size: 2rem;
        cursor: pointer;
    }

    .modal-footer {
        justify-content: flex-end;
        gap: 1rem;
    }

    .btn-ghost {
        background: none;
        border: none;
        color: var(--text-secondary);
        cursor: pointer;
        font-weight: 600;
    }

    @media (max-width: 768px) {
        .options-inputs {
            grid-template-columns: 1fr;
        }
    }
</style>
