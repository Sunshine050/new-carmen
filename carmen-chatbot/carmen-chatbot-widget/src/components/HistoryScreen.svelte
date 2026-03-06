<script>
    import { fade, fly, scale } from "svelte/transition";
    import { cubicOut } from "svelte/easing";
    import { ICONS } from "../lib/icons.js";
    import { STRINGS } from "../lib/constants.js";

    let {
        rooms = [],
        currentRoomId = "",
        isExpanded = false,
        on_switch_room,
        on_create_new,
        on_confirm_delete,
        on_close,
    } = $props();

    function formatDate(dateStr) {
        if (!dateStr) return "";
        try {
            const date = new Date(dateStr);
            return date.toLocaleDateString("th-TH", {
                day: "numeric",
                month: "short",
                year: "2-digit",
                hour: "2-digit",
                minute: "2-digit",
            });
        } catch (e) {
            return dateStr;
        }
    }

    // Custom "decay" transition for room deletion (blur + scale + fade)
    function decay(node, { duration = 400 }) {
        return {
            duration,
            css: (t) => {
                const eased = cubicOut(t);
                return `
                    opacity: ${eased};
                    transform: scale(${0.96 + 0.04 * eased});
                    filter: blur(${(1 - eased) * 12}px);
                `;
            },
        };
    }
</script>

<div
    class="history-screen"
    class:expanded={isExpanded}
    in:scale={{ start: 0.95, duration: 400, opacity: 1, easing: cubicOut }}
    out:fade={{ duration: 200 }}
>
    <div class="history-header">
        <button class="back-btn" onclick={on_close} title="กลับไปที่แชท">
            <svg
                viewBox="0 0 24 24"
                width="20"
                height="20"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <path d="M19 12H5M12 19l-7-7 7-7" />
            </svg>
        </button>
        <h2>ประวัติการสนทนา</h2>
        <button
            class="header-new-chat"
            onclick={on_create_new}
            title="เริ่มแชทใหม่"
        >
            +
        </button>
    </div>

    <div class="history-list">
        {#each rooms as room (room.room_id)}
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div
                class="history-item"
                class:active={room.room_id === currentRoomId}
                onclick={() => on_switch_room(room.room_id)}
                out:decay={{ duration: 450 }}
            >
                <div class="item-main">
                    <div class="item-title">{room.title || "บทสนทนาใหม่"}</div>
                    <div class="item-snippet">
                        {room.lastMessage || "ไม่มีข้อความ"}
                    </div>
                    <div class="item-meta">{formatDate(room.updated_at)}</div>
                </div>
                <button
                    class="item-delete"
                    onclick={(e) => {
                        e.stopPropagation();
                        on_confirm_delete(room.room_id);
                    }}
                    title="ลบ"
                >
                    {@html ICONS.trash}
                </button>
            </div>
        {/each}

        {#if rooms.length === 0}
            <div class="empty-state">
                <div class="empty-icon">{@html ICONS.history}</div>
                <p>{STRINGS.history_empty}</p>
                <button class="start-btn" onclick={on_create_new}
                    >เริ่มบทสนทนาแรก</button
                >
            </div>
        {/if}
    </div>
</div>

<style>
    .history-screen {
        position: absolute !important;
        top: 0 !important;
        left: 0 !important;
        right: 0 !important;
        bottom: 0 !important;
        background: #0f172a !important; /* matches --sidebar-bg */
        z-index: 150 !important;
        display: flex !important;
        flex-direction: column !important;
        color: white !important;
    }

    .history-header {
        padding: 20px 24px !important;
        display: flex !important;
        align-items: center !important;
        gap: 16px !important;
        border-bottom: 1px solid rgba(255, 255, 255, 0.05) !important;
        background: rgba(15, 23, 42, 0.8) !important;
        backdrop-filter: blur(10px) !important;
    }

    .history-header h2 {
        flex: 1 !important;
        margin: 0 !important;
        font-size: 18px !important;
        font-weight: 600 !important;
    }

    .back-btn,
    .header-new-chat {
        width: 36px !important;
        height: 36px !important;
        border-radius: 10px !important;
        border: none !important;
        background: rgba(255, 255, 255, 0.05) !important;
        color: white !important;
        cursor: pointer !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        transition: 0.2s !important;
    }

    .back-btn:hover,
    .header-new-chat:hover {
        background: rgba(255, 255, 255, 0.1) !important;
        transform: scale(1.05) !important;
    }

    .header-new-chat {
        background: var(--primary-gradient) !important;
        font-size: 20px !important;
    }

    .history-list {
        flex: 1 !important;
        overflow-y: auto !important;
        padding: 16px !important;
        display: flex !important;
        flex-direction: column !important;
        gap: 12px !important;
        overscroll-behavior: contain !important;
    }

    .history-item {
        background: rgba(255, 255, 255, 0.03) !important;
        border: 1px solid rgba(255, 255, 255, 0.08) !important;
        border-radius: 16px !important;
        padding: 20px 24px !important;
        cursor: pointer !important;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
        display: flex !important;
        align-items: flex-start !important;
        gap: 20px !important;
        position: relative !important;
        margin-bottom: 16px !important;
        height: auto !important;
        width: 100% !important;
        box-sizing: border-box !important;
    }

    .history-screen.expanded .history-item {
        padding: 28px 32px !important;
        margin-bottom: 20px !important;
        border-radius: 20px !important;
    }

    .history-item:hover {
        background: rgba(255, 255, 255, 0.06) !important;
        border-color: rgba(255, 255, 255, 0.15) !important;
        transform: translateY(-2px) !important;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2) !important;
    }

    .history-item.active {
        background: rgba(var(--primary-color-rgb, 52, 85, 139), 0.1) !important;
        border: 1.5px solid var(--primary-color, #34558b) !important;
    }

    .item-main {
        flex: 1 !important;
        min-width: 0 !important;
        display: flex !important;
        flex-direction: column !important;
        gap: 10px !important;
    }

    .item-title {
        font-weight: 700 !important;
        font-size: 16px !important;
        color: #fff !important;
        line-height: 1.4 !important;
        word-break: break-word !important;
        display: block !important;
        overflow: visible !important;
        text-overflow: unset !important;
        white-space: normal !important;
    }

    .history-screen.expanded .item-title {
        font-size: 19px !important;
    }

    .item-snippet {
        font-size: 13.5px !important;
        color: rgba(255, 255, 255, 0.45) !important;
        line-height: 1.5 !important;
        font-family: var(--font-sarabun) !important;
        word-break: break-word !important;
        display: -webkit-box !important;
        -webkit-line-clamp: 3 !important;
        line-clamp: 3 !important;
        -webkit-box-orient: vertical !important;
        overflow: hidden !important;
        text-overflow: ellipsis !important;
        opacity: 0.8 !important;
    }

    .history-screen.expanded .item-snippet {
        font-size: 15px !important;
        -webkit-line-clamp: 4 !important;
        line-clamp: 4 !important;
    }

    .item-meta {
        font-size: 10px !important;
        color: rgba(255, 255, 255, 0.25) !important;
        text-transform: uppercase !important;
        letter-spacing: 0.1em !important;
        margin-top: 6px !important;
    }

    .item-delete {
        width: 36px !important;
        height: 36px !important;
        border-radius: 10px !important;
        border: none !important;
        background: rgba(255, 255, 255, 0.03) !important;
        color: rgba(255, 255, 255, 0.2) !important;
        cursor: pointer !important;
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1) !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        flex-shrink: 0 !important;
        margin-top: 2px !important;
    }

    .history-screen.expanded .item-delete {
        width: 44px !important;
        height: 44px !important;
        border-radius: 14px !important;
    }

    .item-delete:hover {
        background: rgba(239, 68, 68, 0.1) !important;
        color: #ef4444 !important;
        transform: scale(1.1) !important;
    }

    .item-delete :global(svg) {
        width: 18px !important;
        height: 18px !important;
        stroke-width: 2px !important;
    }

    .history-screen.expanded .item-delete :global(svg) {
        width: 20px !important;
        height: 22px !important;
    }

    .empty-state {
        display: flex !important;
        flex-direction: column !important;
        align-items: center !important;
        justify-content: center !important;
        padding: 60px 20px !important;
        text-align: center !important;
        color: #64748b !important;
    }

    .empty-icon {
        width: 64px !important;
        height: 64px !important;
        background: rgba(255, 255, 255, 0.03) !important;
        border-radius: 20px !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        margin-bottom: 16px !important;
    }

    .empty-icon :global(svg) {
        width: 32px !important;
        height: 32px !important;
        opacity: 0.5 !important;
    }

    .start-btn {
        margin-top: 20px !important;
        padding: 10px 24px !important;
        background: var(--primary-gradient) !important;
        border: none !important;
        border-radius: 12px !important;
        color: white !important;
        font-weight: 600 !important;
        cursor: pointer !important;
        transition: 0.2s !important;
    }

    .start-btn:hover {
        transform: scale(1.05) !important;
        box-shadow: 0 8px 16px rgba(52, 85, 139, 0.3) !important;
    }
</style>
