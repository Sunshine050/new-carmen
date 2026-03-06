<script>
    import { STRINGS } from "../lib/constants.js";

    let { showAlert = $bindable(false), alertData = $bindable({}) } = $props();
</script>

<div class="alert-overlay" class:show={showAlert}>
    <div class="alert-box">
        <div class="alert-icon">{@html alertData.icon || "⚠️"}</div>
        <div class="alert-title">{alertData.title || ""}</div>
        <div class="alert-desc">{alertData.text || ""}</div>
        <div class="alert-actions">
            {#if alertData.cancelText}
                <button
                    class="btn-alert btn-cancel"
                    onclick={() => (showAlert = false)}
                    >{alertData.cancelText}</button
                >
            {/if}
            <button
                class="btn-alert btn-confirm"
                onclick={() => alertData.onConfirm?.()}
                >{alertData.confirmText || STRINGS.alert_confirm}</button
            >
        </div>
    </div>
</div>

<style>
    /* Modals */
    .alert-overlay {
        position: absolute !important;
        inset: 0 !important;
        background: rgba(15, 23, 42, 0.4) !important;
        backdrop-filter: blur(8px) !important;
        z-index: 10000 !important;
        display: none;
        justify-content: center !important;
        align-items: center !important;
        padding: 24px !important;
        animation: fadeIn 0.3s ease !important;
    }
    .alert-overlay.show {
        display: flex !important;
    }

    .alert-box {
        background: rgba(255, 255, 255, 0.95) !important;
        width: 100% !important;
        max-width: 340px !important;
        padding: 28px !important;
        border-radius: 20px !important;
        box-shadow:
            0 20px 40px rgba(0, 0, 0, 0.3),
            0 0 0 1px rgba(255, 255, 255, 0.5) inset !important;
        backdrop-filter: blur(20px) !important;
        text-align: center !important;
        border: 1px solid rgba(0, 0, 0, 0.05) !important;
        animation: scaleUp 0.3s cubic-bezier(0.16, 1, 0.3, 1) !important;
        font-family: var(--font-inter) !important;
    }
    @keyframes scaleUp {
        from {
            transform: scale(0.9);
            opacity: 0;
        }
        to {
            transform: scale(1);
            opacity: 1;
        }
    }

    .alert-icon {
        width: 56px !important;
        height: 56px !important;
        background: #fee2e2 !important;
        color: #ef4444 !important;
        border-radius: 50% !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        margin: 0 auto 16px auto !important;
    }
    :global(.alert-icon svg) {
        width: 28px !important;
        height: 28px !important;
    }

    .alert-title {
        font-weight: 700 !important;
        font-size: 18px !important;
        margin-bottom: 8px !important;
        color: #0f172a !important;
        font-family: var(--font-sarabun) !important;
    }
    .alert-desc {
        color: var(--text-gray) !important;
        font-size: 14px !important;
        line-height: 1.5 !important;
        margin-bottom: 0 !important;
        font-family: var(--font-sarabun) !important;
    }
    .alert-actions {
        display: flex !important;
        gap: 12px !important;
        justify-content: center !important;
        margin-top: 24px !important;
    }

    .btn-alert {
        padding: 12px 20px !important;
        border-radius: 12px !important;
        font-weight: 600 !important;
        cursor: pointer !important;
        transition: 0.2s !important;
        width: 100% !important;
        flex: 1 !important;
        display: flex !important;
        justify-content: center !important;
        align-items: center !important;
        border: transparent;
        font-family: var(--font-inter) !important;
    }
    .btn-confirm {
        background: #ef4444 !important;
        color: white !important;
    }
    .btn-confirm:hover {
        background: #dc2626 !important;
        transform: scale(1.02) !important;
        box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3) !important;
    }
    .btn-cancel {
        background: transparent !important;
        color: var(--text-dark) !important;
        border: 1px solid #cbd5e1 !important;
    }
    .btn-cancel:hover {
        background: #f1f5f9 !important;
    }
</style>
