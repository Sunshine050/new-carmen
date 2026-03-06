<script>
    import { onMount } from "svelte";
    import { ICONS } from "../lib/icons.js";
    import { STRINGS, SUGGESTED_QUESTIONS } from "../lib/constants.js";

    let { on_suggestion_click } = $props();

    let typedTitle = $state("");
    let showDesc = $state(false);
    let showSuggestions = $state(false);

    const fullTitle = STRINGS.welcome_title;

    onMount(() => {
        // Typing effect for title
        let i = 0;
        const typeTitle = () => {
            if (i < fullTitle.length) {
                typedTitle = fullTitle.substring(0, i + 1);
                i++;
                setTimeout(typeTitle, 40);
            } else {
                // Show description after title is done
                showDesc = true;
            }
        };

        setTimeout(typeTitle, 300);
        setTimeout(() => {
            showSuggestions = true;
        }, 1200);
    });
</script>

<div class="welcome-hero">
    <div class="hero-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5z"></path>
            <path d="M2 17l10 5 10-5"></path>
            <path d="M2 12l10 5 10-5"></path>
        </svg>
    </div>
    <h2>{typedTitle}</h2>
    <p
        style="opacity: {showDesc ? 1 : 0}; transform: translateY({showDesc
            ? 0
            : 10}px); transition: all 0.8s ease;"
    >
        {STRINGS.welcome_desc}
    </p>
</div>

{#if showSuggestions}
    <div class="suggestions-container">
        {#each SUGGESTED_QUESTIONS as q, i}
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div
                class="suggestion-chip"
                style="animation-delay: {i * 0.1}s"
                onclick={() => on_suggestion_click(q)}
            >
                {q}
            </div>
        {/each}
    </div>
{/if}

<style>
    /* Welcome Hero */
    .welcome-hero {
        padding: 32px 20px 16px 20px !important;
        text-align: center !important;
        animation: heroFade 0.8s ease-out !important;
    }
    .welcome-hero .hero-icon {
        width: 64px !important;
        height: 64px !important;
        background: var(--primary-gradient) !important;
        border-radius: 22px !important;
        margin: 0 auto 16px auto !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        box-shadow: 0 16px 32px rgba(52, 85, 139, 0.2) !important;
        animation: heroIconPop 1s cubic-bezier(0.16, 1, 0.3, 1) !important;
    }
    :global(.welcome-hero .hero-icon svg) {
        width: 32px !important;
        height: 32px !important;
    }
    .welcome-hero h2 {
        font-size: 20px !important;
        font-weight: 700 !important;
        color: var(--text-dark) !important;
        margin-bottom: 8px !important;
    }
    .welcome-hero p {
        color: var(--text-gray) !important;
        font-size: 14px !important;
        line-height: 1.5 !important;
        max-width: 260px !important;
        margin: 0 auto !important;
    }
    @keyframes heroFade {
        from {
            opacity: 0;
            transform: translateY(20px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
    @keyframes heroIconPop {
        from {
            transform: scale(0.5) rotate(-20deg);
            opacity: 0;
        }
        to {
            transform: scale(1) rotate(0);
            opacity: 1;
        }
    }

    /* Suggestions */
    .suggestions-container {
        display: flex !important;
        flex-wrap: wrap !important;
        align-items: flex-start !important;
        gap: 6px !important;
        margin: 8px 0 !important;
        padding: 0 4px !important;
    }
    .suggestion-chip {
        background: rgba(255, 255, 255, 0.8) !important;
        backdrop-filter: blur(4px) !important;
        border: 1px solid #e2e8f0 !important;
        border-radius: 10px !important;
        padding: 6px 12px !important;
        font-size: 13px !important;
        color: var(--text-dark) !important;
        cursor: pointer !important;
        transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1) !important;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02) !important;
        max-width: 100% !important;
        animation: chipEnter 0.5s cubic-bezier(0.16, 1, 0.3, 1) backwards !important;
    }
    .suggestion-chip:hover {
        background: white !important;
        border-color: var(--primary-color) !important;
        color: var(--primary-color) !important;
        transform: translateY(-2px) !important;
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08) !important;
    }
    @keyframes chipEnter {
        from {
            opacity: 0;
            transform: translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    @media (max-width: 768px) {
        .suggestions-container {
            display: flex !important;
            flex-wrap: wrap !important;
            gap: 8px !important;
            padding: 8px 16px 12px 16px !important;
            margin: 0 !important;
            justify-content: flex-start !important;
        }
        .suggestion-chip {
            flex: 0 1 auto !important;
            font-size: 13px !important;
            padding: 6px 14px !important;
            border-radius: 12px !important;
            background: rgba(255, 255, 255, 0.95) !important;
            border: 1px solid rgba(52, 85, 139, 0.15) !important;
            color: var(--text-dark) !important;
            box-shadow: 0 2px 6px rgba(0, 0, 0, 0.04) !important;
            text-align: left !important;
            line-height: 1.4 !important;
        }
    }
</style>
