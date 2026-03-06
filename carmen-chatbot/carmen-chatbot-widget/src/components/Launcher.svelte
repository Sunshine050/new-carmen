<script>
    import { onMount } from "svelte";
    import { ICONS } from "../lib/icons.js";

    let { isOpen = false, theme = "#34558b", on_toggle } = $props();

    let showTooltip = $state(false);

    onMount(() => {
        // Show tooltip after 2 seconds if chat is not open
        const timer = setTimeout(() => {
            if (!isOpen) showTooltip = true;
        }, 2000);

        // Auto-hide tooltip after 8 seconds
        const hideTimer = setTimeout(() => {
            showTooltip = false;
        }, 10000);

        return () => {
            clearTimeout(timer);
            clearTimeout(hideTimer);
        };
    });

    function handleClick() {
        showTooltip = false;
        on_toggle();
    }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="chat-btn" onclick={handleClick}>
    {@html ICONS.launcher}
</div>

{#if showTooltip && !isOpen}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="chat-tooltip show" onclick={handleClick}>
        <div class="tooltip-avatar">
            {@html ICONS.botAvatar}
        </div>
        <div class="tooltip-content">
            <span class="tooltip-greet">ผู้ช่วย AI พร้อมให้คำแนะนำ</span>
            <span class="tooltip-text"
                >สอบถามข้อมูลคู่มือหรือการใช้งานได้ทันที!</span
            >
        </div>
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
            class="tooltip-close"
            onclick={(e) => {
                e.stopPropagation();
                showTooltip = false;
            }}
        >
            <svg
                viewBox="0 0 24 24"
                width="16"
                height="16"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                ><line x1="18" y1="6" x2="6" y2="18"></line><line
                    x1="6"
                    y1="6"
                    x2="18"
                    y2="18"
                ></line></svg
            >
        </div>
    </div>
{/if}

<style>
    /* Launcher Button */
    .chat-btn {
        width: 70px !important;
        height: 70px !important;
        background: var(--primary-gradient) !important;
        border-radius: 50% !important;
        box-shadow:
            0 12px 32px rgba(0, 0, 0, 0.2),
            inset 0 -4px 0 rgba(0, 0, 0, 0.1),
            0 0 0 2px rgba(255, 255, 255, 0.1) !important;
        cursor: pointer !important;
        display: flex !important;
        justify-content: center !important;
        align-items: center !important;
        transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) !important;
        position: absolute !important;
        bottom: 0 !important;
        right: 0 !important;
        z-index: 2000002 !important;
        overflow: hidden !important;
    }
    .chat-btn::before {
        content: "" !important;
        position: absolute !important;
        inset: 0 !important;
        background: linear-gradient(
            45deg,
            transparent,
            rgba(255, 255, 255, 0.3),
            transparent
        ) !important;
        transform: translateX(-100%) !important;
        transition: 0.8s !important;
    }
    .chat-btn:hover {
        transform: translateY(-10px) scale(1.05) !important;
        box-shadow:
            0 24px 48px rgba(0, 0, 0, 0.35),
            0 0 0 4px rgba(52, 85, 139, 0.2) !important;
        border-radius: 45% !important;
    }
    .chat-btn:hover::before {
        transform: translateX(100%) !important;
    }
    :global(.chat-btn svg) {
        width: 28px !important;
        height: 28px !important;
        transition: transform 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) !important;
        filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2)) !important;
    }
    :global(.chat-btn:hover svg) {
        transform: scale(1.15) rotate(5deg) !important;
    }
    .chat-btn:not(:hover) {
        animation: launcherPulseCircular 4s infinite
            cubic-bezier(0.4, 0, 0.2, 1) !important;
    }
    @keyframes launcherPulseCircular {
        0% {
            box-shadow:
                0 12px 32px rgba(0, 0, 0, 0.2),
                0 0 0 0 rgba(52, 85, 139, 0.4);
        }
        50% {
            box-shadow:
                0 16px 40px rgba(0, 0, 0, 0.25),
                0 0 0 20px rgba(52, 85, 139, 0);
        }
        100% {
            box-shadow:
                0 12px 32px rgba(0, 0, 0, 0.2),
                0 0 0 0 rgba(52, 85, 139, 0);
        }
    }

    /* Tooltip */
    .chat-tooltip {
        position: absolute !important;
        bottom: 88px !important;
        right: 0 !important;
        background: rgba(255, 255, 255, 0.9) !important;
        backdrop-filter: blur(12px) !important;
        padding: 12px 18px !important;
        border-radius: 20px !important;
        box-shadow:
            0 20px 40px rgba(0, 0, 0, 0.12),
            0 0 0 1px rgba(0, 0, 0, 0.04) !important;
        display: none;
        align-items: center !important;
        gap: 14px !important;
        animation: tooltipEnter 0.6s cubic-bezier(0.16, 1, 0.3, 1) backwards !important;
        z-index: 2000001 !important;
        width: max-content !important;
        cursor: pointer !important;
        max-width: 320px !important;
    }
    .chat-tooltip.show {
        display: flex !important;
    }
    @keyframes tooltipEnter {
        0% {
            opacity: 0;
            transform: translateY(20px) scale(0.9);
            filter: blur(4px);
        }
        100% {
            opacity: 1;
            transform: translateY(0) scale(1);
            filter: blur(0);
        }
    }
    .tooltip-avatar {
        width: 38px !important;
        height: 38px !important;
        flex-shrink: 0 !important;
    }
    :global(.tooltip-avatar .avatar-inner) {
        width: 100% !important;
        height: 100% !important;
        background: var(--primary-gradient) !important;
        border-radius: 12px !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        box-shadow: 0 4px 8px rgba(52, 85, 139, 0.2) !important;
    }
    :global(.tooltip-avatar .avatar-inner svg) {
        width: 22px !important;
        height: 22px !important;
    }
    .tooltip-content {
        display: flex !important;
        flex-direction: column !important;
        gap: 3px !important;
    }
    .tooltip-greet {
        font-weight: 700 !important;
        font-size: 13px !important;
        color: var(--primary-color) !important;
        letter-spacing: 0.02em !important;
    }
    .tooltip-text {
        font-size: 13px !important;
        font-weight: 500 !important;
        color: var(--text-dark) !important;
        line-height: 1.4 !important;
    }
    .tooltip-close {
        position: absolute !important;
        top: -8px !important;
        right: -8px !important;
        width: 24px !important;
        height: 24px !important;
        background: white !important;
        color: #94a3b8 !important;
        border-radius: 50% !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1) !important;
        cursor: pointer !important;
        transition: 0.2s !important;
        border: 1px solid rgba(0, 0, 0, 0.05) !important;
    }
    .tooltip-close:hover {
        color: #f87171 !important;
        transform: scale(1.1) !important;
    }
</style>
