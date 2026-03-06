<script>
    import { ICONS } from "../lib/icons.js";
    import RoomDropdown from "./RoomDropdown.svelte";

    let {
        title = "",
        rooms = [],
        currentRoomId = "",
        isExpanded = false,
        showDropdown = $bindable(false),
        on_create_new,
        on_switch_room,
        on_confirm_delete,
        on_expand_toggle,
        on_close,
        on_drag_start,
    } = $props();
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
    class="chat-header"
    onmousedown={on_drag_start}
    ontouchstart={(e) => on_drag_start(e, true)}
    style="cursor: move;"
>
    <div class="header-info">
        <div class="avatar-wrapper">
            {@html ICONS.botAvatar}
        </div>
        <div class="title-wrapper">
            <h3>{title}</h3>
            <div class="status-indicator">
                <span class="dot"></span> คลังความรู้ AI พร้อมบริการ
            </div>
        </div>
    </div>
    <div
        class="header-tools"
        onmousedown={(e) => e.stopPropagation()}
        ontouchstart={(e) => e.stopPropagation()}
    >
        <RoomDropdown
            {rooms}
            {currentRoomId}
            bind:showDropdown
            {on_create_new}
            {on_switch_room}
            {on_confirm_delete}
        />

        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <div
            class="icon-btn"
            id="carmen-expand-btn"
            title="ขยายหน้าจอ"
            onclick={on_expand_toggle}
        >
            ⛶
        </div>
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <div
            class="icon-btn"
            id="carmen-close-btn"
            title="ปิด"
            onclick={on_close}
        >
            {@html ICONS.close}
        </div>
    </div>
</div>

<style>
    /* Header */
    .chat-header {
        background: var(--primary-gradient) !important;
        padding: 24px 28px !important;
        display: flex !important;
        justify-content: space-between !important;
        align-items: center !important;
        color: white !important;
        flex-shrink: 0 !important;
        position: relative !important;
        z-index: 100 !important;
        user-select: none !important;
        -webkit-user-select: none !important;
    }
    .header-info {
        display: flex !important;
        align-items: center !important;
        gap: 14px !important;
    }
    .avatar-wrapper {
        width: 48px !important;
        height: 48px !important;
        background: linear-gradient(
            135deg,
            rgba(255, 255, 255, 0.25),
            rgba(255, 255, 255, 0.1)
        ) !important;
        border-radius: 16px !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        backdrop-filter: blur(12px) !important;
        border: 1px solid rgba(255, 255, 255, 0.3) !important;
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1) !important;
        animation: avatarFloat 4s infinite ease-in-out !important;
    }
    @keyframes avatarFloat {
        0%,
        100% {
            transform: translateY(0);
        }
        50% {
            transform: translateY(-4px);
        }
    }
    :global(.avatar-wrapper svg) {
        width: 28px !important;
        height: 28px !important;
        filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1)) !important;
    }

    .chat-header h3 {
        font-size: 19px !important;
        font-weight: 700 !important;
        letter-spacing: -0.02em !important;
        margin: 0 !important;
        text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1) !important;
    }
    .status-indicator {
        font-size: 11px !important;
        font-weight: 500 !important;
        opacity: 0.8 !important;
        display: flex !important;
        align-items: center !important;
        gap: 6px !important;
        margin-top: 4px !important;
        text-transform: uppercase !important;
        letter-spacing: 0.05em !important;
    }
    .status-indicator .dot {
        width: 7px !important;
        height: 7px !important;
        background: #22c55e !important;
        border-radius: 50% !important;
        display: inline-block !important;
        border: 1.5px solid rgba(255, 255, 255, 0.4) !important;
        box-shadow: 0 0 8px #22c55e !important;
    }
    .header-tools {
        display: flex !important;
        gap: 10px !important;
        align-items: center !important;
    }

    .icon-btn {
        width: 36px !important;
        height: 36px !important;
        border-radius: 14px !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        background: rgba(255, 255, 255, 0.15) !important;
        cursor: pointer !important;
        transition: 0.2s !important;
    }
    .icon-btn:hover {
        background: rgba(255, 255, 255, 0.25) !important;
        transform: scale(1.05) !important;
    }
    :global(.icon-btn svg) {
        width: 22px !important;
        height: 22px !important;
        fill: white !important;
    }

    @media (max-width: 768px) {
        .icon-btn {
            width: 44px !important;
            height: 44px !important;
        }
    }
</style>
