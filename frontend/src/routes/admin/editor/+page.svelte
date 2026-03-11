<script>
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { alertStore } from "$lib/alertStore";
    import { goto } from "$app/navigation";
    import { API_URL, resolveImageUrl } from "$lib/api";
    import Modal from "$lib/components/Modal.svelte";

    let showExitModal = false;
    let showDeleteModal = false;

    export let data; // Suppress SvelteKit warning
    export let form; // Suppress SvelteKit warning

    // Data State
    let quiz = {
        title: "New Quiz",
        questions: [
            {
                text: "",
                image_url: "",
                options: ["", "", "", ""],
                correct_answer: -1,
                correct_answers: [],
                time_limit: 20,
                points: 1000,
                type: "Quiz",
                answer_type: "Single select",
            },
        ],
    };

    let selectedIndex = 0;
    let fileInput;

    onMount(async () => {
        const token = localStorage.getItem("admin_token");
        if (!token) {
            goto("/admin/login");
            return;
        }

        // Load existing quiz if ID is in URL
        const id = $page.url.searchParams.get("id");
        if (id) {
            try {
                const res = await fetch(`${API_URL}/api/games/${id}`, {
                    headers: { Authorization: `Bearer ${token}` },
                });
                if (res.ok) {
                    const data = await res.json();
                    // Migrate old data if necessary
                    if (data.questions) {
                        data.questions = data.questions.map((q) => ({
                            ...q,
                            correct_answers:
                                q.correct_answers ||
                                (q.correct_answer !== -1
                                    ? [q.correct_answer]
                                    : []),
                        }));
                    }
                    quiz = data;
                } else if (res.status === 401) {
                    goto("/admin/login");
                }
            } catch (e) {
                console.error("Load error:", e);
            }
        }
    });

    $: currentQuestion =
        quiz && quiz.questions ? quiz.questions[selectedIndex] : null;

    async function handleFileUpload(file) {
        if (!file) return;

        const token = localStorage.getItem("admin_token");
        const formData = new FormData();
        formData.append("file", file);

        try {
            const res = await fetch(`${API_URL}/api/upload`, {
                method: "POST",
                headers: { Authorization: `Bearer ${token}` },
                body: formData,
            });
            const data = await res.json();
            if (res.ok) {
                currentQuestion.image_url = data.image_url;
                quiz.questions = [...quiz.questions]; // Trigger reactivity
            }
        } catch (e) {
            console.error("Upload error:", e);
        }
    }

    function onFileChange(e) {
        const file = e.target.files[0];
        handleFileUpload(file);
    }

    function onDrop(e) {
        e.preventDefault();
        const file = e.dataTransfer.files[0];
        handleFileUpload(file);
    }

    function onDragOver(e) {
        e.preventDefault();
    }

    function addQuestion() {
        quiz.questions = [
            ...quiz.questions,
            {
                text: "",
                image_url: "",
                options: ["", "", "", ""],
                correct_answer: -1,
                correct_answers: [],
                time_limit: 20,
                points: 1000,
                type: "Quiz",
                answer_type: "Single select",
            },
        ];
        selectedIndex = quiz.questions.length - 1;
    }

    function deleteQuestion() {
        if (quiz.questions.length > 1) {
            quiz.questions = quiz.questions.filter(
                (_, i) => i !== selectedIndex,
            );
            selectedIndex = Math.max(0, selectedIndex - 1);
        }
    }

    function duplicateQuestion() {
        const clone = JSON.parse(JSON.stringify(currentQuestion));
        quiz.questions = [
            ...quiz.questions.slice(0, selectedIndex + 1),
            clone,
            ...quiz.questions.slice(selectedIndex + 1),
        ];
        selectedIndex++;
    }

    async function saveQuiz() {
        if (!quiz.questions || quiz.questions.length === 0) {
            alertStore.showAlert(
                "กรุณาสร้างคำถามอย่างน้อย 1 ข้อก่อนบันทึก (Please add at least 1 question before saving)",
            );
            return false;
        }

        const token = localStorage.getItem("admin_token");
        const isUpdate = quiz.id && quiz.id !== "";
        const url = isUpdate
            ? `${API_URL}/api/games/${quiz.id}`
            : `${API_URL}/api/games`;
        const method = isUpdate ? "PUT" : "POST";

        try {
            const res = await fetch(url, {
                method: method,
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(quiz),
            });
            if (res.ok) {
                return true;
            } else {
                const data = await res.json();
                alertStore.showAlert(
                    "Failed to save: " + (data.error || "Unknown error"),
                );
                return false;
            }
        } catch (e) {
            console.error(e);
            alertStore.showAlert(
                "Network error. Please check your connection.",
            );
            return false;
        }
    }

    function deleteQuiz() {
        showDeleteModal = true;
    }

    async function confirmDeleteQuiz() {
        showDeleteModal = false;
        // If it's a new quiz, just exit. If it has an ID, send DELETE.
        if (quiz.id) {
            const token = localStorage.getItem("admin_token");
            try {
                const res = await fetch(`${API_URL}/api/games/${quiz.id}`, {
                    method: "DELETE",
                    headers: { Authorization: `Bearer ${token}` },
                });
                if (res.ok) {
                    goto("/admin");
                }
            } catch (e) {
                console.error(e);
            }
        } else {
            goto("/admin");
        }
    }
    function exitEditor() {
        showExitModal = true;
    }

    async function confirmExitWithSave() {
        const success = await saveQuiz();
        if (success) {
            showExitModal = false;
            goto("/admin");
        }
    }

    function confirmExitNoSave() {
        showExitModal = false;
        goto("/admin");
    }

    function toggleAnswer(index) {
        if (!currentQuestion) return;

        if (currentQuestion.answer_type === "Multi-select") {
            const arr = currentQuestion.correct_answers || [];
            if (arr.includes(index)) {
                currentQuestion.correct_answers = arr.filter(
                    (i) => i !== index,
                );
            } else {
                currentQuestion.correct_answers = [...arr, index];
            }
        } else {
            // Single select
            currentQuestion.correct_answers = [index];
            currentQuestion.correct_answer = index;
        }
        quiz.questions = [...quiz.questions];
    }
</script>

<Modal
    bind:show={showExitModal}
    title="ออกจากระบบบันทึก (Exit Editor)"
    message="คุณต้องการบันทึกการเปลี่ยนแปลงก่อนออกจากหน้านี้หรือไม่?"
    confirmText="บันทึกและออก (Save & Exit)"
    cancelText="ออกโดยไม่บันทึก (Exit without saving)"
    type="info"
    classic={true}
    on:confirm={confirmExitWithSave}
    on:cancel={confirmExitNoSave}
/>

<Modal
    bind:show={showDeleteModal}
    title="ลบ Quiz (Delete Quiz)"
    message="คุณแน่ใจหรือไม่ว่าต้องการลบ Quiz นี้? การกระทำนี้ไม่สามารถย้อนกลับได้"
    confirmText="ลบทันที (Delete Now)"
    cancelText="ยกเลิก (Cancel)"
    type="danger"
    classic={true}
    on:confirm={confirmDeleteQuiz}
    on:cancel={() => (showDeleteModal = false)}
/>

<div class="editor-shell">
    <!-- Top Header -->
    <header class="top-nav">
        <div class="left">
            <button class="btn-nav btn-exit" on:click={exitEditor}>
                EXIT
            </button>
            <input
                class="header-quiz-title"
                bind:value={quiz.title}
                placeholder="Enter quiz title..."
            />
        </div>
        <div class="right">
            <button class="btn-nav btn-delete-quiz" on:click={deleteQuiz}>
                Delete Quiz
            </button>
            <button
                class="btn-nav btn-save"
                class:disabled={!quiz.questions || quiz.questions.length === 0}
                on:click={async () => {
                    const success = await saveQuiz();
                    if (success) goto("/admin");
                }}
            >
                Save
            </button>
        </div>
    </header>

    <div class="editor-main">
        <!-- Left Sidebar: Question List -->
        <aside class="question-list">
            <div class="scroll-area">
                {#each (quiz && quiz.questions) || [] as q, i}
                    <div
                        class="q-thumb {selectedIndex === i ? 'active' : ''}"
                        on:click={() => (selectedIndex = i)}
                    >
                        <span class="q-num">{i + 1} {q.type}</span>
                        <div class="q-preview">
                            <div class="p-text">{q.text || "Question"}</div>
                            <div class="p-grid">
                                <span></span><span></span><span></span><span
                                ></span>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
            <div class="list-actions">
                <button class="btn-add" on:click={addQuestion}
                    >+ Add question</button
                >
            </div>
        </aside>

        <!-- Center: Workspace -->
        <main class="workspace">
            <div class="editor-canvas">
                <div class="question-input-wrap">
                    <input
                        type="text"
                        placeholder="Start typing your question"
                        bind:value={currentQuestion.text}
                        class="q-input"
                    />
                </div>

                <div
                    class="media-placeholder"
                    on:drop={onDrop}
                    on:dragover={onDragOver}
                >
                    {#if currentQuestion.image_url}
                        <img
                            src={resolveImageUrl(currentQuestion.image_url)}
                            alt="Question media"
                            class="uploaded-img"
                        />
                        <button
                            class="btn-remove-media"
                            on:click={() => {
                                currentQuestion.image_url = "";
                                quiz.questions = [...quiz.questions];
                            }}>×</button
                        >
                    {:else}
                        <div
                            class="plus-box"
                            on:click={() => fileInput.click()}
                        >
                            +
                        </div>
                        <p>Find and insert media</p>
                        <div class="upload-links">
                            <span
                                class="link"
                                on:click={() => fileInput.click()}
                                >Upload file</span
                            > or drag here to upload
                        </div>
                    {/if}
                    <input
                        type="file"
                        accept="image/*"
                        hidden
                        bind:this={fileInput}
                        on:change={onFileChange}
                    />
                </div>

                {#if currentQuestion.type === "True or False"}
                    <div class="answer-grid tf-mode">
                        <div
                            class="ans-box blue {currentQuestion.correct_answers.includes(
                                0,
                            )
                                ? 'selected'
                                : ''}"
                            on:click={() => toggleAnswer(0)}
                        >
                            <div class="shape diamond"></div>
                            <span class="tf-text">True</span>
                            <div class="selection-indicator">
                                {#if currentQuestion.correct_answers.includes(0)}
                                    <span class="check">✓</span>
                                {/if}
                            </div>
                        </div>
                        <div
                            class="ans-box red {currentQuestion.correct_answers.includes(
                                1,
                            )
                                ? 'selected'
                                : ''}"
                            on:click={() => toggleAnswer(1)}
                        >
                            <div class="shape triangle"></div>
                            <span class="tf-text">False</span>
                            <div class="selection-indicator">
                                {#if currentQuestion.correct_answers.includes(1)}
                                    <span class="check">✓</span>
                                {/if}
                            </div>
                        </div>
                    </div>
                {:else}
                    <div class="answer-grid">
                        {#each [0, 1, 2, 3] as i}
                            <div
                                class="ans-box {i === 0
                                    ? 'red'
                                    : i === 1
                                      ? 'blue'
                                      : i === 2
                                        ? 'yellow'
                                        : 'green'} {currentQuestion.correct_answers.includes(
                                    i,
                                )
                                    ? 'selected'
                                    : ''}"
                                on:click={() => toggleAnswer(i)}
                            >
                                <div
                                    class="shape {i === 0
                                        ? 'triangle'
                                        : i === 1
                                          ? 'diamond'
                                          : i === 2
                                            ? 'circle'
                                            : 'square'}"
                                ></div>
                                <input
                                    type="text"
                                    placeholder={i < 2
                                        ? `Add answer ${i + 1}`
                                        : `Add answer ${i + 1} (optional)`}
                                    bind:value={currentQuestion.options[i]}
                                    on:click|stopPropagation
                                />
                                <div class="selection-indicator">
                                    {#if currentQuestion.correct_answers.includes(i)}
                                        <span class="check">✓</span>
                                    {/if}
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>
        </main>

        <!-- Right Sidebar: Properties -->
        <aside class="properties">
            <div class="prop-group">
                <label>Question type</label>
                <select bind:value={currentQuestion.type}>
                    <option>Quiz</option>
                    <option>True or False</option>
                </select>
            </div>

            <div class="prop-group">
                <label>Time limit</label>
                <select bind:value={currentQuestion.time_limit}>
                    <option value={5}>5 seconds</option>
                    <option value={10}>10 seconds</option>
                    <option value={20}>20 seconds</option>
                    <option value={30}>30 seconds</option>
                </select>
            </div>

            <div class="prop-group">
                <label
                    >Points {currentQuestion.points === 0
                        ? "(No points)"
                        : ""}</label
                >
                <input
                    type="number"
                    bind:value={currentQuestion.points}
                    min="0"
                    step="100"
                    class="prop-input"
                    placeholder="Points (0 = No points)"
                />
            </div>

            <div class="prop-group">
                <label>Answer options</label>
                <select bind:value={currentQuestion.answer_type}>
                    <option>Single select</option>
                    <option>Multi-select</option>
                </select>
            </div>

            <div class="prop-footer">
                <button class="btn-del" on:click={deleteQuestion}>Delete</button
                >
                <button class="btn-dup" on:click={duplicateQuestion}
                    >Duplicate</button
                >
            </div>
        </aside>
    </div>
</div>

<style lang="scss">
    .editor-shell {
        height: 100vh;
        display: flex;
        flex-direction: column;
        background: #f2f2f2;
        color: #333;
    }

    .top-nav {
        height: 60px;
        background: white;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0 1rem;
        box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        z-index: 100;

        .left,
        .right {
            display: flex;
            align-items: center;
            gap: 1rem;
        }

        .header-quiz-title {
            font-weight: 700;
            font-size: 1.1rem;
            border: 1px solid transparent;
            padding: 4px 12px;
            border-radius: 6px;
            background: #f8f9fa;
            width: 320px;
            transition: all 0.2s;
            &:focus {
                background: white;
                border: 1px solid #00d2ff;
                outline: none;
                box-shadow: 0 0 0 3px rgba(0, 210, 255, 0.1);
            }
        }

        .btn-nav {
            padding: 0.6rem 1.2rem;
            border-radius: 6px;
            font-weight: 600;
            cursor: pointer;
            border: none;
            transition: all 0.2s;
            font-family: var(--font-main);
            font-size: 0.9rem;
            display: flex;
            align-items: center;
            justify-content: center;

            &:active {
                transform: translateY(1px);
            }
        }

        .btn-exit {
            background: #6e41e2;
            color: white;
            &:hover {
                background: #5d35c2;
            }
        }

        .btn-delete-quiz {
            background: #e21b3c;
            color: white;
            &:hover {
                background: #d01937;
            }
        }

        .btn-save {
            background: #1368ce;
            color: white;
            &:hover:not(.disabled) {
                background: #0f56aa;
            }
            &.disabled {
                opacity: 0.5;
                cursor: not-allowed;
                filter: grayscale(0.5);
            }
        }
    }

    .editor-main {
        flex: 1;
        display: flex;
        overflow: hidden;
    }

    // Left Sidebar
    .question-list {
        width: 200px;
        background: white;
        border-right: 1px solid #ddd;
        display: flex;
        flex-direction: column;

        .scroll-area {
            flex: 1;
            overflow-y: auto;
            padding: 1rem;
            display: flex;
            flex-direction: column;
            gap: 1rem;
        }

        .q-thumb {
            padding: 0.5rem;
            border: 2px solid transparent;
            border-radius: 8px;
            cursor: pointer;
            position: relative;

            &.active {
                border-color: #1368ce;
                background: #f0f7ff;
            }

            .q-num {
                font-size: 0.7rem;
                font-weight: 700;
                color: #666;
            }

            .q-preview {
                background: white;
                border: 1px solid #ddd;
                border-radius: 4px;
                height: 80px;
                margin-top: 0.2rem;
                display: flex;
                flex-direction: column;
                padding: 4px;

                .p-text {
                    font-size: 0.6rem;
                    text-align: center;
                    flex: 1;
                }
                .p-grid {
                    display: grid;
                    grid-template-columns: 1fr 1fr;
                    gap: 2px;
                    height: 20px;
                    span {
                        background: #eee;
                        border-radius: 1px;
                    }
                }
            }
        }

        .list-actions {
            padding: 1rem;
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
            border-top: 1px solid #ddd;

            button {
                width: 100%;
                padding: 0.75rem;
                border-radius: 8px;
                font-weight: 700;
                cursor: pointer;
            }

            .btn-add {
                background: #1368ce;
                color: white;
                border: none;
            }
            .btn-generate {
                background: white;
                color: #46178f;
                border: 2px solid #46178f;
            }
        }
    }

    // Center Workspace
    .workspace {
        flex: 1;
        background: #46178f; // Kahoot purple
        padding: 2rem;
        display: flex;
        justify-content: center;
        align-items: center;
        overflow-y: auto;
        background-image: radial-gradient(
                circle at 20% 30%,
                rgba(255, 255, 255, 0.05) 0%,
                transparent 40%
            ),
            radial-gradient(
                circle at 80% 70%,
                rgba(255, 255, 255, 0.05) 0%,
                transparent 40%
            );

        .editor-canvas {
            width: 95%;
            max-width: 1200px; // Slightly larger but not infinite to maintain readability
            display: flex;
            flex-direction: column;
            gap: 1.5rem;
            animation: slideUp 0.6s cubic-bezier(0.23, 1, 0.32, 1);
        }

        .question-input-wrap {
            background: white;
            border-radius: 12px;
            padding: 1.5rem;
            box-shadow:
                0 8px 0 rgba(0, 0, 0, 0.2),
                0 15px 30px rgba(0, 0, 0, 0.1);

            .q-input {
                width: 100%;
                border: none;
                font-size: 1.8rem;
                text-align: center;
                font-weight: 700;
                color: #333;
                background: transparent;
                &:focus {
                    outline: none;
                }
            }
        }

        .media-placeholder {
            height: 320px;
            background: rgba(255, 255, 255, 0.9);
            border-radius: 16px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            gap: 1.5rem;
            box-shadow: 0 8px 0 rgba(0, 0, 0, 0.1);
            border: 2px dashed rgba(0, 0, 0, 0.1);
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;

            &:hover {
                background: white;
                transform: scale(1.01);
                box-shadow: 0 12px 0 rgba(0, 0, 0, 0.1);
            }

            .plus-box {
                width: 60px;
                height: 60px;
                background: #f2f2f2;
                border: 3px solid #ddd;
                display: flex;
                align-items: center;
                justify-content: center;
                font-size: 2rem;
                color: #999;
                border-radius: 12px;
                transition: transform 0.2s;
            }

            p {
                font-weight: 700;
                color: #333;
                font-size: 1.2rem;
            }

            .upload-links {
                font-size: 0.95rem;
                color: #666;
                .link {
                    color: #1368ce;
                    font-weight: 800;
                    cursor: pointer;
                    text-decoration: none;
                    &:hover {
                        text-decoration: underline;
                    }
                }
            }

            .uploaded-img {
                width: 100%;
                height: 100%;
                object-fit: contain;
                border-radius: 12px;
            }

            .btn-remove-media {
                position: absolute;
                top: 10px;
                right: 10px;
                background: rgba(0, 0, 0, 0.5);
                color: white;
                border: none;
                border-radius: 50%;
                width: 30px;
                height: 30px;
                cursor: pointer;
                font-size: 1.2rem;
                display: flex;
                align-items: center;
                justify-content: center;
                &:hover {
                    background: rgba(0, 0, 0, 0.8);
                }
            }
        }

        .answer-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 1.5rem;

            &.tf-mode {
                grid-template-columns: 1fr 1fr;

                .ans-box {
                    height: 140px;
                    justify-content: space-between;
                    padding: 0 2rem;
                    box-shadow: 0 8px 0 rgba(0, 0, 0, 0.2);

                    .tf-text {
                        font-size: 1.8rem;
                        font-weight: 800;
                        color: white;
                        flex: 1;
                        margin-left: 1.5rem;
                        text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
                    }

                    input[type="radio"] {
                        width: 50px;
                        height: 50px;
                        flex-shrink: 0;
                        aspect-ratio: 1 / 1;
                        border: 4px solid white;
                        background: rgba(255, 255, 255, 0.2);
                        border-radius: 50%;
                        cursor: pointer;
                        appearance: none;
                        display: grid;
                        place-items: center;
                        transition: all 0.2s;

                        &:checked {
                            background: white;
                            &::after {
                                content: "✔";
                                color: #26890c;
                                font-size: 1.8rem;
                                font-weight: bold;
                            }
                        }
                    }
                }
            }

            .ans-box {
                display: flex;
                align-items: center;
                padding: 1.8rem;
                background: white;
                border-radius: 12px;
                gap: 1.5rem;
                position: relative;
                overflow: hidden;
                transition: all 0.1s ease;
                top: 0;
                box-shadow: 0 8px 0 rgba(0, 0, 0, 0.2);

                &:active {
                    top: 4px;
                    box-shadow: 0 4px 0 rgba(0, 0, 0, 0.2);
                }

                &.red {
                    background: #e21b3c;
                    &:hover {
                        background: #f12c4d;
                    }
                }
                &.blue {
                    background: #1368ce;
                    &:hover {
                        background: #1a7df5;
                    }
                }
                &.yellow {
                    background: #d89e00;
                    &:hover {
                        background: #f0af00;
                    }
                }
                &.green {
                    background: #26890c;
                    &:hover {
                        background: #2faf0f;
                    }
                }

                input[type="text"] {
                    flex: 1;
                    border: none;
                    font-size: 1.3rem;
                    font-weight: 700;
                    color: white;
                    background: transparent;
                    &::placeholder {
                        color: rgba(255, 255, 255, 0.8);
                    }
                    &:focus {
                        outline: none;
                    }
                }

                &.selected {
                    box-shadow:
                        0 0 0 4px rgba(255, 255, 255, 0.4),
                        0 8px 15px rgba(0, 0, 0, 0.3);
                    transform: scale(1.02);
                }

                .selection-indicator {
                    width: 48px;
                    height: 48px;
                    border: 3px solid rgba(255, 255, 255, 0.4);
                    border-radius: 50%;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    background: rgba(255, 255, 255, 0.1);
                    transition: all 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);

                    .check {
                        color: white;
                        font-size: 2rem;
                        font-weight: 900;
                        animation: popIn 0.3s ease;
                    }

                    :global(.ans-box.selected) & {
                        background: #fff;
                        border-color: #fff;
                        .check {
                            color: #26890c;
                        }
                    }
                }

                input[type="radio"],
                input[type="checkbox"] {
                    display: none;
                }

                .shape {
                    width: 45px;
                    height: 45px;
                    flex-shrink: 0;
                    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
                }

                &.red {
                    .shape {
                        background: white;
                        clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
                    }
                }
                &.blue {
                    .shape {
                        background: white;
                        transform: rotate(45deg);
                        width: 35px;
                        height: 35px;
                        margin-left: 5px;
                        margin-right: 5px;
                    }
                }
                &.yellow {
                    .shape {
                        background: white;
                        border-radius: 50%;
                    }
                }
                &.green {
                    .shape {
                        background: white;
                    }
                }
            }
        }
    }

    @keyframes slideUp {
        from {
            opacity: 0;
            transform: translateY(30px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    // Right Sidebar
    .properties {
        width: 300px;
        background: white;
        border-left: 1px solid #ddd;
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 2rem;

        .prop-group {
            display: flex;
            flex-direction: column;
            gap: 0.5rem;

            label {
                font-weight: 700;
                font-size: 0.9rem;
                color: #333;
            }
            select,
            .prop-input {
                padding: 0.75rem;
                border: 1px solid #ccc;
                border-radius: 4px;
                font-size: 1rem;
                background: #f9f9f9;
                cursor: pointer;
                font-family: var(--font-main);
            }

            .prop-input {
                cursor: text;
                &:focus {
                    outline: none;
                    border-color: #1368ce;
                    background: white;
                }
            }
        }

        .prop-footer {
            margin-top: auto;
            display: flex;
            gap: 1rem;

            button {
                flex: 1;
                padding: 0.75rem;
                border-radius: 4px;
                font-weight: 700;
                cursor: pointer;
                border: none;
            }

            .btn-del {
                background: #ddd;
                color: #333;
            }
            .btn-dup {
                background: white;
                color: #333;
                border: 1px solid #ccc;
            }
        }
    }
</style>
