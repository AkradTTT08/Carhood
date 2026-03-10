<script>
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";

    const menuItems = [
        { name: "Master Quiz", path: "/admin", icon: "📝" },
        { name: "Start Quiz", path: "/admin/start", icon: "🎮" },
        { name: "History", path: "/admin/history", icon: "🕰️" },
    ];

    function logout() {
        localStorage.removeItem("admin_token");
        goto("/admin/login");
    }
</script>

<aside class="sidebar glass-card">
    <div class="logo-container">
        <img
            src="/image/logotttspacequiz.png"
            alt="SPACE QUIZ"
            class="sidebar-logo"
        />
    </div>

    <nav class="menu">
        {#each menuItems as item}
            <a
                href={item.path}
                class="menu-item {$page.url.pathname === item.path
                    ? 'active'
                    : ''}"
            >
                <span class="item-icon">{item.icon}</span>
                <span class="item-text">{item.name}</span>
            </a>
        {/each}
    </nav>

    <div class="footer">
        <button class="btn-3d red logout-btn" on:click={logout}>
            LOGOUT 🚪
        </button>
    </div>
</aside>

<style lang="scss">
    .sidebar {
        width: 280px;
        height: 100vh;
        margin: 0;
        display: flex;
        flex-direction: column;
        padding: 2rem 1.5rem;
        border-radius: 0 30px 30px 0;
        position: fixed;
        left: 0;
        top: 0;
        z-index: 100;

        .logo-container {
            margin-bottom: 3rem;
            display: flex;
            justify-content: center;
            padding: 0 0.5rem;

            .sidebar-logo {
                max-width: 130px;
                height: auto;
                filter: drop-shadow(0 0 5px rgba(0, 210, 255, 0.2));
                animation: pulse 3s ease-in-out infinite;
            }
        }

        @keyframes pulse {
            0%,
            100% {
                transform: scale(1);
                opacity: 0.9;
            }
            50% {
                transform: scale(1.03);
                opacity: 1;
            }
        }

        .menu {
            flex: 1;
            display: flex;
            flex-direction: column;
            gap: 1.5rem;
            padding: 0 0.5rem;

            .menu-item {
                display: flex;
                align-items: center;
                gap: 1.2rem;
                padding: 1.2rem;
                text-decoration: none;
                transition: all 0.1s ease;
                position: relative;
                top: 0;
                border: 4px solid rgba(255, 255, 255, 0.2);
                font-family: var(--font-pixel);
                font-size: 0.7rem;
                justify-content: flex-start;
                background: rgba(255, 255, 255, 0.05);
                box-shadow: 0 6px 0 rgba(0, 0, 0, 0.5);
                color: var(--text-secondary);

                .item-icon {
                    font-size: 1.4rem;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    width: 30px;
                }

                .item-text {
                    font-weight: 700;
                    letter-spacing: 1px;
                    line-height: 1.5;
                }

                &:hover {
                    top: -2px;
                    background: rgba(255, 255, 255, 0.1);
                    box-shadow: 0 8px 0 rgba(0, 0, 0, 0.5);
                    color: white;
                    border-color: rgba(255, 255, 255, 0.4);
                }

                &:active {
                    top: 4px;
                    box-shadow: none;
                }

                &.active {
                    background: var(--primary);
                    color: white;
                    border-color: #fff;
                    box-shadow: 0 6px 0 #2a0e56;

                    &:hover {
                        box-shadow: 0 8px 0 #2a0e56;
                    }

                    &:active {
                        box-shadow: none;
                    }
                }
            }
        }

        .footer {
            padding-top: 2rem;
            border-top: 4px solid rgba(255, 255, 255, 0.05);

            .logout-btn {
                width: 100%;
                padding: 1.2rem;
                font-size: 0.8rem;
                letter-spacing: 2px;
            }
        }
    }
</style>
