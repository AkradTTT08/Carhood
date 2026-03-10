<script>
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
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
                time_limit: 20,
                points: "Standard",
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
                const res = await fetch(
                    `http://127.0.0.1:8081/api/games/${id}`,
                    {
                        headers: { Authorization: `Bearer ${token}` },
                    },
                );
                if (res.ok) {
                    quiz = await res.json();
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
            const res = await fetch("http://127.0.0.1:8081/api/upload", {
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
                time_limit: 20,
                points: "Standard",
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
        const token = localStorage.getItem("admin_token");
        const isUpdate = quiz.id && quiz.id !== "";
        const url = isUpdate
            ? `http://127.0.0.1:8081/api/games/${quiz.id}`
            : "http://127.0.0.1:8081/api/games";
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
                alert("Failed to save: " + (data.error || "Unknown error"));
                return false;
            }
        } catch (e) {
            console.error(e);
            alert("Network error. Please check your connection.");
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
                const res = await fetch(
                    `http://localhost:8081/api/games/${quiz.id}`,
                    {
                        method: "DELETE",
                        headers: { Authorization: `Bearer ${token}` },
                    },
                );
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
</script>

<Modal
    bind:show={showExitModal}
    title="ออกจากระบบบันทึก (Exit Editor)"
    message="คุณต้องการบันทึกการเปลี่ยนแปลงก่อนออกจากหน้านี้หรือไม่?"
    confirmText="บันทึกและออก (Save & Exit)"
    cancelText="ออกโดยไม่บันทึก (Exit without saving)"
    type="info"
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
    on:confirm={confirmDeleteQuiz}
    on:cancel={() => (showDeleteModal = false)}
/>

<div class="editor-shell">
    <!-- Top Header -->
    <header class="top-nav">
        <div class="left">
            <button
                class="btn-3d purple btn-exit"
                on:click={exitEditor}
                style="padding: 0.6rem 1.2rem; min-width: 90px; font-weight: 800; border: none;"
            >
                EXIT
            </button>
            <input
                class="header-quiz-title"
                bind:value={quiz.title}
                placeholder="Enter quiz title..."
            />
        </div>
        <div class="right">
            <button
                class="btn-3d red btn-delete-quiz"
                on:click={deleteQuiz}
                style="padding: 0.5rem 1rem; font-size: 0.9rem; margin-right: 0.5rem;"
            >
                Delete Quiz
            </button>
            <button
                class="btn-3d blue btn-save"
                on:click={async () => {
                    const success = await saveQuiz();
                    if (success) goto("/admin");
                }}
                style="padding: 0.5rem 1.5rem; font-size: 1rem;"
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
                            src={currentQuestion.image_url}
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
                        <div class="ans-box blue">
                            <div class="shape diamond"></div>
                            <span class="tf-text">True</span>
                            <input
                                type="radio"
                                name="correct"
                                value={0}
                                bind:group={currentQuestion.correct_answer}
                            />
                        </div>
                        <div class="ans-box red">
                            <div class="shape triangle"></div>
                            <span class="tf-text">False</span>
                            <input
                                type="radio"
                                name="correct"
                                value={1}
                                bind:group={currentQuestion.correct_answer}
                            />
                        </div>
                    </div>
                {:else}
                    <div class="answer-grid">
                        <div class="ans-box red">
                            <div class="shape triangle"></div>
                            <input
                                type="text"
                                placeholder="Add answer 1"
                                bind:value={currentQuestion.options[0]}
                            />
                            <input
                                type="radio"
                                name="correct"
                                value={0}
                                bind:group={currentQuestion.correct_answer}
                            />
                        </div>
                        <div class="ans-box blue">
                            <div class="shape diamond"></div>
                            <input
                                type="text"
                                placeholder="Add answer 2"
                                bind:value={currentQuestion.options[1]}
                            />
                            <input
                                type="radio"
                                name="correct"
                                value={1}
                                bind:group={currentQuestion.correct_answer}
                            />
                        </div>
                        <div class="ans-box yellow">
                            <div class="shape circle"></div>
                            <input
                                type="text"
                                placeholder="Add answer 3 (optional)"
                                bind:value={currentQuestion.options[2]}
                            />
                            <input
                                type="radio"
                                name="correct"
                                value={2}
                                bind:group={currentQuestion.correct_answer}
                            />
                        </div>
                        <div class="ans-box green">
                            <div class="shape square"></div>
                            <input
                                type="text"
                                placeholder="Add answer 4 (optional)"
                                bind:value={currentQuestion.options[3]}
                            />
                            <input
                                type="radio"
                                name="correct"
                                value={3}
                                bind:group={currentQuestion.correct_answer}
                            />
                        </div>
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
                <label>Points</label>
                <select bind:value={currentQuestion.points}>
                    <option>Standard</option>
                    <option>Double points</option>
                    <option>No points</option>
                </select>
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
            padding: 4px 8px;
            border-radius: 4px;
            background: #f8f9fa;
            width: 300px;
            &:focus {
                background: white;
                border-color: #ccc;
                outline: none;
            }
        }

        .btn-delete-quiz {
            background: #e21b3c;
            color: white;
            border: none;
            margin-right: 0.5rem;
            &:hover {
                background: #d01937;
            }
        }

        .btn-exit {
            // Specific overrides if needed, but removing generic button style is better
        }

        .btn-save {
            background: #1368ce;
            color: white;
            border: none;
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

                input[type="radio"] {
                    width: 40px;
                    height: 40px;
                    flex-shrink: 0;
                    aspect-ratio: 1 / 1;
                    cursor: pointer;
                    accent-color: #fff;
                    appearance: none;
                    background: rgba(255, 255, 255, 0.2);
                    border: 3px solid rgba(255, 255, 255, 0.5);
                    border-radius: 50%;
                    display: grid;
                    place-items: center;
                    transition: all 0.2s;

                    &:checked {
                        background: white;
                        border-color: white;
                        &::after {
                            content: "✔";
                            color: #26890c;
                            font-size: 1.5rem;
                            font-weight: bold;
                        }
                    }
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
            select {
                padding: 0.75rem;
                border: 1px solid #ccc;
                border-radius: 4px;
                font-size: 1rem;
                background: #f9f9f9;
                cursor: pointer;
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
