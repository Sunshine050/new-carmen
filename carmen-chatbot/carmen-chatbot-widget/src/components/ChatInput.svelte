<script>
    import { ICONS } from "../lib/icons.js";
    import { STRINGS } from "../lib/constants.js";

    let { inputText = $bindable(""), on_send } = $props();

    let inputEl;

    // Auto-Expand Textarea
    export function setupAutoExpand() {
        if (!inputEl) return;
        inputEl.style.height = "auto";
        inputEl.style.height = inputEl.scrollHeight + "px";
    }

    export function resetHeight() {
        if (inputEl) inputEl.style.height = "auto";
    }

    function handleInput() {
        setupAutoExpand();
    }

    function handleKeydown(e) {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            on_send();
            resetHeight();
        }
    }
</script>

<div class="chat-footer">
    <textarea
        class="chat-input"
        rows="1"
        placeholder={STRINGS.input_placeholder}
        bind:value={inputText}
        bind:this={inputEl}
        oninput={handleInput}
        onkeydown={handleKeydown}
    ></textarea>
    <button
        class="send-btn"
        onclick={() => {
            on_send();
            resetHeight();
        }}
        title="ส่งข้อความ"
    >
        {@html ICONS.send}
    </button>
</div>

<style>
    /* Footer & Input */
    .chat-footer {
        padding: 20px 24px 24px 24px !important;
        background: var(--glass-bg) !important;
        backdrop-filter: blur(12px) !important;
        border-top: 1px solid var(--glass-border) !important;
        display: flex !important;
        gap: 12px !important;
        align-items: center !important;
    }
    .chat-input {
        flex: 1 !important;
        padding: 14px 20px !important;
        border-radius: 16px !important;
        border: 1px solid #e2e8f0 !important;
        outline: none !important;
        background: white !important;
        font-family: var(--font-sarabun) !important;
        font-size: 15px !important;
        transition:
            border-color 0.2s,
            box-shadow 0.2s !important;
        resize: none !important;
        min-height: 50px !important;
        max-height: 120px !important;
        overflow-y: auto !important;
        line-height: 1.5 !important;
    }
    .chat-input:focus {
        border-color: var(--primary-color) !important;
        box-shadow:
            0 0 0 4px rgba(52, 85, 139, 0.1),
            0 4px 12px rgba(52, 85, 139, 0.05) !important;
        background: white !important;
    }
    .send-btn {
        width: 48px !important;
        height: 48px !important;
        background: #0f172a !important;
        color: white !important;
        border-radius: 14px !important;
        cursor: pointer !important;
        border: none !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        transition: var(--transition-snappy) !important;
    }
    .send-btn:hover {
        background: var(--primary-color) !important;
        transform: scale(1.05) !important;
    }
    :global(.send-btn svg) {
        width: 24px !important;
        height: 24px !important;
        fill: white !important;
    }

    .icon-btn-footer {
        background: transparent !important;
        border: none !important;
        color: var(--text-gray) !important;
        cursor: pointer !important;
        transition: 0.2s !important;
        display: flex !important;
        align-items: center !important;
    }
    .icon-btn-footer:hover {
        color: var(--primary-color) !important;
        transform: scale(1.1) !important;
    }

    @media (max-width: 768px) {
        .chat-footer {
            padding: 12px 16px !important;
            padding-bottom: max(
                16px,
                env(safe-area-inset-bottom, 16px)
            ) !important;
            background: rgba(255, 255, 255, 0.98) !important;
            border-top: 1px solid rgba(0, 0, 0, 0.05) !important;
            position: relative !important;
            z-index: 10 !important;
        }
        .chat-input {
            font-size: 16px !important;
            min-height: 48px !important;
        }
        .icon-btn-footer {
            width: 44px !important;
            height: 44px !important;
        }
    }
</style>
