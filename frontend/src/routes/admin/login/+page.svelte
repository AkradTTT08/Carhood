<script>
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { API_URL } from "$lib/api";

    let username = "";
    let password = "";
    let error = "";
    let loading = false;

    async function handleLogin() {
        loading = true;
        error = "";
        try {
            const res = await fetch(`${API_URL}/api/login`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
            });
            const data = await res.json();
            if (res.ok) {
                localStorage.setItem("admin_token", data.token);
                window.location.href = "/admin";
            } else {
                error = data.error || "Login failed";
            }
        } catch (e) {
            error = "Could not connect to server";
        } finally {
            loading = false;
        }
    }
</script>

<div class="login-bg">
    <img src="/image/BGUSER.svg" alt="Login Background" class="bg-image" />
</div>

<div class="login-container">
    <div class="glass-card login-box">
        <div class="logo-wrapper">
            <img
                src="/image/logotttspacequiz.png"
                alt="TTT SPACE QUIZ"
                class="system-logo"
            />
        </div>
        <h1>Admin Login</h1>
        <p class="subtitle">Welcome back, Captain!</p>

        {#if error}
            <div class="error-msg">{error}</div>
        {/if}

        <form on:submit|preventDefault={handleLogin}>
            <div class="form-group">
                <label for="username">Username</label>
                <input
                    type="text"
                    id="username"
                    class="input-field"
                    bind:value={username}
                    placeholder="ENTER USERNAME"
                    required
                />
            </div>

            <div class="form-group">
                <label for="password">Password</label>
                <input
                    type="password"
                    id="password"
                    class="input-field"
                    bind:value={password}
                    placeholder="ENTER PASSWORD"
                    required
                />
            </div>

            <button type="submit" class="btn-primary" disabled={loading}>
                {loading ? "AUTHENTICATING..." : "LOGIN"}
            </button>
        </form>
    </div>
</div>

<style lang="scss">
    .login-bg {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 1;
        background: #000;

        .bg-image {
            width: 100%;
            height: 100%;
            object-fit: cover;
            image-rendering: auto;
        }
    }

    .login-container {
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 1rem;
        position: relative;
        z-index: 2;
    }

    .login-box {
        width: 100%;
        max-width: 400px;
        text-align: center;
        animation: fadeIn 0.8s ease-out;

        .logo-wrapper {
            margin-bottom: 1.5rem;
            display: flex;
            justify-content: center;
            animation: float 4s ease-in-out infinite;
        }

        .system-logo {
            max-width: 150px;
            height: auto;
            filter: drop-shadow(0 0 10px rgba(0, 210, 255, 0.4));
        }

        h1 {
            margin: 0 0 1rem;
            color: var(--text-primary);
            font-size: 1.2rem;
            text-shadow: 0 0 10px rgba(0, 210, 255, 0.5);
        }

        .subtitle {
            color: var(--text-secondary);
            margin-bottom: 2rem;
            font-family: var(--font-pixel);
            font-size: 0.6rem;
            letter-spacing: 1px;
        }
    }

    .form-group {
        text-align: left;
        margin-bottom: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;

        label {
            font-family: var(--font-pixel);
            color: var(--text-secondary);
            font-size: 0.6rem;
            padding-left: 0.5rem;
        }
    }

    .error-msg {
        background: rgba(226, 27, 60, 0.2);
        color: #ff3e3e;
        padding: 0.8rem;
        border: 2px solid #ff3e3e;
        margin-bottom: 1.5rem;
        font-family: var(--font-pixel);
        font-size: 0.6rem;
        line-height: 1.4;
    }

    .btn-primary {
        width: 100%;
        margin-top: 1rem;
        font-size: 0.8rem;

        &:disabled {
            opacity: 0.5;
            cursor: not-allowed;
            top: 0 !important;
            box-shadow: 0 6px 0 #2a0e56 !important;
        }
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(20px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>
