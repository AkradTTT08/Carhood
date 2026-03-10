<script>
    import { onMount } from "svelte";

    let username = "";
    let password = "";
    let error = "";
    let loading = false;

    async function handleLogin() {
        loading = true;
        error = "";
        try {
            const res = await fetch("http://localhost:8081/api/login", {
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

<div class="login-container">
    <div class="login-box">
        <h1>Admin Login</h1>
        <p class="subtitle">Welcome back, Captain!</p>

        {#if error}
            <div class="error-msg">{error}</div>
        {/if}

        <form on:submit|preventDefault={handleLogin}>
            <div class="input-group">
                <label for="username">Username</label>
                <input
                    type="text"
                    id="username"
                    bind:value={username}
                    placeholder="Enter username"
                    required
                />
            </div>

            <div class="input-group">
                <label for="password">Password</label>
                <input
                    type="password"
                    id="password"
                    bind:value={password}
                    placeholder="Enter password"
                    required
                />
            </div>

            <button
                type="submit"
                class="btn-3d {loading ? 'loading' : ''}"
                disabled={loading}
            >
                {loading ? "Authenticating..." : "Login"}
            </button>
        </form>
    </div>
</div>

<style lang="scss">
    .login-container {
        min-height: 100vh;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, #46178f 0%, #25074d 100%);
        font-family: "Inter", sans-serif;
    }

    .login-box {
        background: white;
        padding: 3rem;
        border-radius: 20px;
        box-shadow:
            0 20px 40px rgba(0, 0, 0, 0.3),
            0 10px 10px rgba(0, 0, 0, 0.2);
        width: 100%;
        max-width: 400px;
        text-align: center;
        animation: slideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1);

        h1 {
            margin: 0 0 0.5rem;
            color: #333;
            font-size: 2rem;
            font-weight: 800;
        }

        .subtitle {
            color: #666;
            margin-bottom: 2rem;
        }
    }

    .input-group {
        text-align: left;
        margin-bottom: 1.5rem;

        label {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 600;
            color: #444;
            font-size: 0.9rem;
        }

        input {
            width: 100%;
            padding: 0.8rem 1rem;
            border: 2px solid #eee;
            border-radius: 10px;
            font-size: 1rem;
            transition: all 0.3s ease;
            box-sizing: border-box;

            &:focus {
                outline: none;
                border-color: #46178f;
                box-shadow: 0 0 0 4px rgba(70, 23, 143, 0.1);
            }
        }
    }

    .error-msg {
        background: #ffebeb;
        color: #e21b3c;
        padding: 0.8rem;
        border-radius: 10px;
        margin-bottom: 1.5rem;
        font-weight: 600;
        font-size: 0.9rem;
        border: 1px solid #ffd1d1;
    }

    .btn-3d {
        width: 100%;
        padding: 1rem;
        font-size: 1.1rem;
        font-weight: 700;
        color: white;
        background: #46178f;
        border: none;
        border-radius: 12px;
        cursor: pointer;
        position: relative;
        transition:
            transform 0.1s,
            background 0.3s;
        box-shadow: 0 6px 0 #2d0f5d;
        margin-top: 1rem;

        &:hover {
            background: #551db1;
        }

        &:active {
            transform: translateY(4px);
            box-shadow: 0 2px 0 #2d0f5d;
        }

        &:disabled {
            background: #ccc;
            box-shadow: 0 6px 0 #999;
            cursor: not-allowed;
            transform: none;
        }

        &.loading {
            opacity: 0.8;
        }
    }

    @keyframes slideUp {
        from {
            opacity: 0;
            transform: translateY(40px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>
