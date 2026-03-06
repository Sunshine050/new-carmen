<script>
    import { STRINGS } from "../lib/constants.js";
    import { onMount } from "svelte";

    let {
        rooms = [],
        currentRoomId = "",
        showDropdown = $bindable(false),
        on_switch_room,
        on_create_new,
        on_confirm_delete,
    } = $props();

    let roomListEl;

    // Scroll Containment on Room List
    onMount(() => {
        if (!roomListEl) return;
        roomListEl.addEventListener(
            "wheel",
            (e) => {
                const isScrollingDown = e.deltaY > 0;
                const isScrollingUp = e.deltaY < 0;
                const isAtBottom =
                    roomListEl.scrollTop + roomListEl.clientHeight >=
                    roomListEl.scrollHeight - 1;
                const isAtTop = roomListEl.scrollTop <= 0;

                if (
                    (isScrollingDown && isAtBottom) ||
                    (isScrollingUp && isAtTop)
                ) {
                    e.preventDefault();
                }
                e.stopPropagation();
            },
            { passive: false },
        );
    });
</script>

<div class="room-dropdown-container" id="carmenRoomDropdownContainer">
    <button
        type="button"
        class="icon-btn room-dropdown-btn"
        onclick={() => (showDropdown = !showDropdown)}
        title="ประวัติการสนทนา"
    >
        <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
        >
            <path
                d="M14 9a2 2 0 0 1-2 2H6l-4 4V4c0-1.1.9-2 2-2h8a2 2 0 0 1 2 2z"
            ></path>
            <path d="M18 9h2a2 2 0 0 1 2 2v11l-4-4h-6a2 2 0 0 1-2-2v-1"></path>
        </svg>
    </button>
    <div class="room-dropdown-menu" class:show={showDropdown}>
        <div class="dropdown-header">
            <span>ประวัติแชท</span>
            <button
                class="new-chat-btn"
                title="เริ่มแชทใหม่"
                onclick={() => {
                    showDropdown = false;
                    on_create_new();
                }}>+</button
            >
        </div>
        <div class="room-list" bind:this={roomListEl}>
            {#each rooms as room}
                <!-- svelte-ignore a11y_click_events_have_key_events -->
                <!-- svelte-ignore a11y_no_static_element_interactions -->
                <div
                    class="room-dropdown-item"
                    class:active={room.room_id === currentRoomId}
                    onclick={() => on_switch_room(room.room_id)}
                >
                    <div class="room-title" title={room.title || "บทสนทนาใหม่"}>
                        {room.title || "บทสนทนาใหม่"}
                    </div>
                    <button
                        class="delete-room-btn"
                        onclick={(e) => {
                            e.stopPropagation();
                            on_confirm_delete(room.room_id);
                        }}>×</button
                    >
                </div>
            {/each}
            {#if rooms.length === 0}
                <div
                    style="text-align:center; padding:16px; color:#94a3b8; font-size:13px;"
                >
                    {STRINGS.history_empty}
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    /* Room Dropdown */
    .room-dropdown-container {
        position: relative !important;
        display: none !important;
    }
    :global(.chat-box.expanded) .room-dropdown-container {
        display: block !important;
    }

    /* Toggle Button - matches ChatHeader .icon-btn */
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
        border: none !important;
        color: white !important;
    }
    .icon-btn:hover {
        background: rgba(255, 255, 255, 0.25) !important;
        transform: scale(1.05) !important;
    }
    :global(.icon-btn svg) {
        stroke: white !important;
    }

    .room-dropdown-menu {
        position: absolute !important;
        top: 100% !important;
        right: 0 !important;
        width: 250px !important;
        margin-top: 12px !important;
        background: rgba(30, 41, 59, 0.98) !important;
        backdrop-filter: blur(16px) saturate(180%) !important;
        -webkit-backdrop-filter: blur(16px) saturate(180%) !important;
        border: 1px solid rgba(255, 255, 255, 0.1) !important;
        border-radius: 16px !important;
        box-shadow:
            0 20px 40px rgba(0, 0, 0, 0.5),
            0 0 0 1px rgba(255, 255, 255, 0.05) inset !important;
        overflow: hidden !important;
        display: none !important;
        flex-direction: column !important;
        z-index: 200 !important;
        animation: dropdownIn 0.3s cubic-bezier(0.16, 1, 0.3, 1) !important;
        transform-origin: top right !important;
    }
    .room-dropdown-menu.show {
        display: flex !important;
    }
    @keyframes dropdownIn {
        from {
            opacity: 0;
            transform: scale(0.9) translateY(-10px);
        }
        to {
            opacity: 1;
            transform: scale(1) translateY(0);
        }
    }

    .dropdown-header {
        padding: 16px !important;
        border-bottom: 1px solid rgba(255, 255, 255, 0.05) !important;
        display: flex !important;
        justify-content: space-between !important;
        align-items: center !important;
        color: white !important;
        font-weight: 600 !important;
        font-size: 14px !important;
    }

    .new-chat-btn {
        width: 28px !important;
        height: 28px !important;
        background: var(--primary-gradient) !important;
        border: none !important;
        color: white !important;
        border-radius: 8px !important;
        cursor: pointer !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        font-size: 18px !important;
        transition: all 0.2s !important;
    }
    .new-chat-btn:hover {
        transform: scale(1.1) !important;
        box-shadow: 0 4px 12px rgba(52, 85, 139, 0.4) !important;
    }

    .room-list {
        flex: 1 !important;
        max-height: 280px !important;
        overflow-y: auto !important;
        padding: 8px !important;
        display: flex !important;
        flex-direction: column !important;
        gap: 4px !important;
    }
    .room-list::-webkit-scrollbar {
        width: 4px !important;
    }
    .room-list::-webkit-scrollbar-track {
        background: transparent !important;
    }
    .room-list::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.1) !important;
        border-radius: 10px !important;
    }
    .room-list:hover::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.2) !important;
    }

    .room-dropdown-item {
        padding: 12px 14px !important;
        border-radius: 10px !important;
        color: #94a3b8 !important;
        cursor: pointer !important;
        display: flex !important;
        align-items: center !important;
        justify-content: space-between !important;
        transition: all 0.2s !important;
        background: transparent !important;
        border: 1px solid transparent !important;
    }
    .room-dropdown-item:hover {
        background: rgba(255, 255, 255, 0.08) !important;
        color: white !important;
        border-color: rgba(255, 255, 255, 0.05) !important;
    }
    .room-dropdown-item.active {
        background: rgba(52, 85, 139, 0.2) !important;
        color: white !important;
        border-color: rgba(52, 85, 139, 0.4) !important;
    }

    .room-title {
        font-size: 13.5px !important;
        font-weight: 500 !important;
        white-space: nowrap !important;
        overflow: hidden !important;
        text-overflow: ellipsis !important;
        flex: 1 !important;
        margin-right: 8px !important;
    }

    .delete-room-btn {
        opacity: 0;
        color: #f87171 !important;
        font-size: 18px !important;
        border: none !important;
        background: transparent !important;
        cursor: pointer !important;
        transition: 0.2s !important;
        padding: 0 4px !important;
        line-height: 1 !important;
    }
    .room-dropdown-item:hover .delete-room-btn {
        opacity: 1 !important;
    }
    .delete-room-btn:hover {
        color: #ef4444 !important;
        transform: scale(1.2) !important;
    }

    @media (max-width: 768px) {
        .room-dropdown-menu {
            width: calc(100vw - 32px) !important;
            max-width: 300px !important;
            right: 16px !important;
        }
        #carmenRoomDropdownContainer {
            display: none !important;
        }
    }
</style>
