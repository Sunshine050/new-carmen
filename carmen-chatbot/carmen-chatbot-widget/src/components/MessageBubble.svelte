<script>
    import { ICONS } from "../lib/icons.js";
    import { formatMessageContent } from "../lib/content-formatter.js";

    let {
        sender,
        message,
        msgId = null,
        sources = null,
        apiBase = "",
        isStreaming = false,
        isQueued = false,
        isError = false,
        timestamp = null,
        on_feedback = null,
        on_retry = null,
    } = $props();
    let copied = $state(false);
    let feedbackGiven = $state(null);

    function formatTime(ts) {
        if (!ts) return "";
        try {
            const d = new Date(ts);
            return d.toLocaleTimeString("th-TH", {
                hour: "2-digit",
                minute: "2-digit",
            });
        } catch {
            return "";
        }
    }

    function handleFeedbackClick(score, type) {
        if (feedbackGiven) return;
        feedbackGiven = type;
    }

    function copyText() {
        const tempEl = document.createElement("div");
        tempEl.innerHTML = formatMessageContent(message, apiBase);
        const text = tempEl.textContent || tempEl.innerText;

        const fallbackCopy = () => {
            const ta = document.createElement("textarea");
            ta.value = text;
            ta.style.position = "fixed";
            ta.style.opacity = "0";
            document.body.appendChild(ta);
            ta.select();
            try {
                document.execCommand("copy");
                copied = true;
                setTimeout(() => (copied = false), 2000);
            } catch (err) {
                console.error("Fallback copy failed", err);
            }
            document.body.removeChild(ta);
        };

        if (navigator.clipboard && window.isSecureContext) {
            navigator.clipboard
                .writeText(text)
                .then(() => {
                    copied = true;
                    setTimeout(() => (copied = false), 2000);
                })
                .catch(fallbackCopy);
        } else {
            fallbackCopy();
        }
    }
</script>

<div
    class="msg {sender === 'user' ? 'user' : 'bot-msg'}"
    class:queued={isQueued}
    class:error={isError}
>
    {#if sender === "user"}
        {message}
        {#if isQueued}
            <div class="queued-badge">🕐 รอคิว</div>
        {/if}
        {#if timestamp && !isQueued}
            <div class="timestamp user-timestamp">{formatTime(timestamp)}</div>
        {/if}
    {:else}
        {#if isError}
            <div class="error-container">
                <span>⚠️ เกิดข้อผิดพลาดในการเชื่อมต่อ</span>
                {#if on_retry}
                    <button class="retry-btn" onclick={on_retry}>
                        🔄 ลองอีกครั้ง
                    </button>
                {/if}
            </div>
        {:else}
            <div class="bot-content">
                <span>{@html formatMessageContent(message, apiBase)}</span>
            </div>
        {/if}

        {#if !isStreaming && !isError && message}
            <div class="tools-row">
                {#if timestamp}
                    <div class="timestamp bot-timestamp">
                        {formatTime(timestamp)}
                    </div>
                {/if}
                <div class="tools-container">
                    <button
                        class="copy-btn"
                        title={copied ? "คัดลอกแล้ว!" : "คัดลอกข้อมูล"}
                        onclick={copyText}
                    >
                        {#if copied}
                            ✓
                        {:else}
                            {@html ICONS.copy}
                        {/if}
                    </button>
                    {#if msgId && on_feedback}
                        <div class="separator"></div>
                        <div class="feedback-group">
                            <button
                                class="feedback-btn {feedbackGiven === 'up'
                                    ? 'active-up'
                                    : ''}"
                                title="มีประโยชน์"
                                onclick={() => handleFeedbackClick(1, "up")}
                                disabled={feedbackGiven !== null}
                            >
                                {@html ICONS.thumbsUp}
                            </button>
                            <button
                                class="feedback-btn {feedbackGiven === 'down'
                                    ? 'active-down'
                                    : ''}"
                                title="ไม่ถูกต้อง"
                                onclick={() => handleFeedbackClick(-1, "down")}
                                disabled={feedbackGiven !== null}
                            >
                                {@html ICONS.thumbsDown}
                            </button>
                        </div>
                    {/if}
                </div>
            </div>
        {/if}
    {/if}
</div>

<style>
    /* Message Bubbles */
    .msg {
        width: fit-content !important;
        max-width: 88% !important;
        padding: 14px 20px !important;
        font-size: 15px !important;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04) !important;
        font-family: var(--font-sarabun) !important;
        position: relative !important;
        animation: msgPop 0.5s cubic-bezier(0.22, 1, 0.36, 1) backwards !important;
        will-change: transform, opacity !important;
    }
    @keyframes msgPop {
        from {
            opacity: 0;
            transform: scale(0.9) translateY(10px);
        }
        to {
            opacity: 1;
            transform: scale(1) translateY(0);
        }
    }

    :global(.msg img) {
        max-width: 100% !important;
        height: auto !important;
        border-radius: 12px !important;
        margin: 8px 0 !important;
        display: block !important;
    }

    .msg.user {
        background: #0f172a !important;
        color: white !important;
        align-self: flex-end !important;
        margin-left: auto !important;
        border-radius: 20px 20px 4px 20px !important;
        box-shadow: 0 4px 12px rgba(15, 23, 42, 0.15) !important;
        padding-bottom: 10px !important;
    }
    .msg.bot-msg {
        background: white !important;
        color: var(--text-dark) !important;
        align-self: flex-start !important;
        margin-right: auto !important;
        border-radius: 20px 20px 20px 4px !important;
        border: 1px solid #f1f5f9 !important;
        padding-bottom: 12px !important;
    }

    /* Timestamp */
    .timestamp {
        font-size: 10px !important;
        font-family: var(--font-sarabun) !important;
        margin-top: 4px !important;
        opacity: 0.5 !important;
    }
    .user-timestamp {
        color: rgba(255, 255, 255, 0.7) !important;
        text-align: right !important;
    }
    .bot-timestamp {
        color: var(--text-gray) !important;
    }

    /* Tools Row (timestamp + copy/feedback on same line) */
    .tools-row {
        display: flex !important;
        align-items: center !important;
        justify-content: space-between !important;
        margin-top: 6px !important;
        gap: 8px !important;
    }

    /* Message Tools (Copy/Feedback) */
    .tools-container {
        display: flex !important;
        align-items: center !important;
        gap: 8px !important;
        opacity: 0.5 !important;
        transition: 0.2s !important;
        margin-left: auto !important;
    }
    .msg:hover .tools-container {
        opacity: 1 !important;
    }

    .copy-btn,
    .feedback-btn {
        background: transparent !important;
        border: none !important;
        cursor: pointer !important;
        padding: 4px !important;
        display: flex !important;
        align-items: center !important;
        color: var(--text-gray) !important;
        transition: 0.2s !important;
    }
    .copy-btn:hover:not(:disabled) {
        color: var(--primary-color) !important;
    }
    .feedback-btn:hover:not(:disabled) {
        color: #f59e0b !important;
    }
    .feedback-btn:disabled {
        cursor: default !important;
    }
    .feedback-btn.active-up {
        color: #10b981 !important;
    }
    .feedback-btn.active-down {
        color: #ef4444 !important;
    }

    .separator {
        width: 1px !important;
        height: 12px !important;
        background: #e2e8f0 !important;
    }
    .feedback-group {
        display: flex !important;
        gap: 4px !important;
    }

    /* Error & Retry */
    .msg.error {
        border: 1px solid #fecaca !important;
        background: #fff5f5 !important;
    }
    .error-container {
        display: flex !important;
        flex-direction: column !important;
        gap: 8px !important;
        align-items: flex-start !important;
        color: #dc2626 !important;
        font-size: 14px !important;
    }
    .retry-btn {
        background: var(--primary-gradient) !important;
        color: white !important;
        border: none !important;
        border-radius: 8px !important;
        padding: 6px 14px !important;
        font-size: 13px !important;
        font-weight: 500 !important;
        cursor: pointer !important;
        font-family: var(--font-sarabun) !important;
        transition: 0.2s !important;
    }
    .retry-btn:hover {
        transform: scale(1.02) !important;
        box-shadow: 0 4px 12px rgba(52, 85, 139, 0.3) !important;
    }

    /* Markdown Styles */
    :global(.msg ul),
    :global(.msg ol) {
        margin: 8px 0 !important;
        padding-left: 24px !important;
    }
    :global(.msg li) {
        margin-bottom: 4px !important;
        line-height: 1.5 !important;
    }
    :global(.msg ul) {
        list-style-type: disc !important;
    }
    :global(.msg ol) {
        list-style-type: decimal !important;
    }
    :global(.carmen-link) {
        color: var(--primary-color) !important;
        text-decoration: none !important;
        font-weight: 500 !important;
    }
    :global(.carmen-link:hover) {
        text-decoration: underline !important;
    }

    /* Image Preview */
    :global(.image-preview-container) {
        padding: 12px 24px !important;
        background: #f8fafc !important;
        border-top: 1px solid #e2e8f0 !important;
        position: relative !important;
    }
    :global(.preview-box) {
        width: 60px !important;
        height: 60px !important;
        border-radius: 10px !important;
        overflow: hidden !important;
        border: 2px solid white !important;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
    }
    :global(.preview-box img) {
        width: 100% !important;
        height: 100% !important;
        object-fit: cover !important;
    }
    :global(#clear-image-btn) {
        position: absolute !important;
        top: 6px !important;
        left: 66px !important;
        background: #fee2e2 !important;
        color: #ef4444 !important;
        border: none !important;
        border-radius: 50% !important;
        width: 24px !important;
        height: 24px !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        cursor: pointer !important;
    }

    @media (max-width: 768px) {
        .msg {
            padding: 12px 16px !important;
            max-width: 90% !important;
        }
    }

    /* Queued user message */
    .msg.queued {
        opacity: 0.7 !important;
    }
    .queued-badge {
        font-size: 11px !important;
        color: rgba(255, 255, 255, 0.85) !important;
        margin-top: 4px !important;
        font-family: var(--font-sarabun) !important;
        font-weight: 400 !important;
    }
</style>
