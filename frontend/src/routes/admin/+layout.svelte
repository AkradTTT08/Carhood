<script>
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import AdminSidebar from "$lib/components/AdminSidebar.svelte";

    let authenticated = false;
    let checking = true;

    // Routes that should NOT show the sidebar
    $: hideSidebar =
        $page.url.pathname === "/admin/login" ||
        $page.url.pathname === "/admin/editor" ||
        $page.url.pathname.startsWith("/play/preview");

    onMount(() => {
        const token = localStorage.getItem("admin_token");
        const isLoginPage = $page.url.pathname === "/admin/login";

        if (!token && !isLoginPage) {
            window.location.href = "/admin/login";
        } else if (token && isLoginPage) {
            window.location.href = "/admin";
        } else {
            authenticated = true;
            checking = false;
        }
    });
</script>

{#if !checking}
    {#if hideSidebar}
        <slot />
    {:else}
        <div class="admin-layout">
            <AdminSidebar />
            <main class="content">
                <slot />
            </main>
        </div>
    {/if}
{/if}

<style lang="scss">
    :global(body) {
        margin: 0;
        background: radial-gradient(circle at top left, #46178f, #250818);
        color: white;
        min-height: 100vh;
        overflow-x: hidden;
    }

    .admin-layout {
        display: flex;
        min-height: 100vh;
    }

    .content {
        flex: 1;
        padding: 2rem;
        overflow-y: auto;
        margin-left: 280px; // Match sidebar width
    }
</style>
