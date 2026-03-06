<script>
    import { onMount } from "svelte";
    import { getCssText } from "./lib/styles.js";
    import { ChatService } from "./lib/api-client.js";
    import { DEFAULT_CONFIG, STRINGS } from "./lib/constants.js";
    import Launcher from "./components/Launcher.svelte";
    import ChatWindow from "./components/ChatWindow.svelte";

    // Props from main.js
    let { config = {}, shadowRoot = null } = $props();

    // Resolved config using $derived to avoid state tracking issues in Svelte 5
    const apiBase = $derived(
        (config.apiBase || DEFAULT_CONFIG.apiBase).replace(/\/$/, ""),
    );
    const bu = $derived(config.bu);
    const username = $derived(config.user);
    const theme = $derived(config.theme || DEFAULT_CONFIG.theme);
    const title = $derived(config.title || DEFAULT_CONFIG.title);
    const promptExtend = $derived(config.promptExtend || null);

    // State
    let isOpen = $state(false);
    let isExpanded = $state(false);
    let currentRoomId = $state(null);
    let requestClose = $state(false);

    // Services
    const api = $derived(apiBase ? new ChatService(apiBase) : null);

    // Inject CSS into Shadow DOM
    onMount(() => {
        if (shadowRoot) {
            const style = document.createElement("style");
            style.textContent = getCssText(theme);
            shadowRoot.prepend(style);

            // Load Google Fonts into main document (only once)
            if (
                !document.querySelector(
                    'link[href*="fonts.googleapis.com/css2?family=Inter"]',
                )
            ) {
                const link = document.createElement("link");
                link.rel = "stylesheet";
                link.href =
                    "https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Sarabun:wght@400;500;600;700&display=swap";
                document.head.appendChild(link);
            }
        }
    });

    function toggleChat() {
        if (isOpen) {
            // Signal ChatWindow to play close animation
            requestClose = true;
        } else {
            isOpen = true;
            requestClose = false;
        }
    }

    function toggleExpand() {
        isExpanded = !isExpanded;
    }

    function closeChat() {
        // Called by ChatWindow after close animation finishes
        isOpen = false;
        isExpanded = false;
        requestClose = false;
    }
</script>

<Launcher {isOpen} {theme} on_toggle={toggleChat} />

{#if isOpen}
    <ChatWindow
        {api}
        {bu}
        {username}
        {apiBase}
        {theme}
        {title}
        {promptExtend}
        {isExpanded}
        {requestClose}
        bind:currentRoomId
        on_close={closeChat}
        on_toggle_expand={toggleExpand}
    />
{/if}
