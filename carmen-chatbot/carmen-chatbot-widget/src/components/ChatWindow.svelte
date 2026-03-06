<script>
    import { onMount, tick } from "svelte";
    import { ICONS } from "../lib/icons.js";
    import { STRINGS, SUGGESTED_QUESTIONS } from "../lib/constants.js";
    import { formatMessageContent } from "../lib/content-formatter.js";
    import WelcomeScreen from "./WelcomeScreen.svelte";
    import MessageBubble from "./MessageBubble.svelte";
    import AlertOverlay from "./AlertOverlay.svelte";
    import ChatHeader from "./ChatHeader.svelte";
    import ChatInput from "./ChatInput.svelte";

    let {
        api,
        bu,
        username,
        apiBase,
        theme,
        title,
        promptExtend,
        isExpanded = false,
        requestClose = false,
        currentRoomId = $bindable(null),
        on_close,
        on_toggle_expand,
    } = $props();

    // Watch for requestClose from Launcher
    $effect(() => {
        if (requestClose && !isClosing) {
            handleClose();
        }
    });

    // State
    let messages = $state([]);
    let messageQueue = $state([]);
    let inputText = $state("");
    let isStreaming = $state(false);
    let streamingText = $state("");
    let streamingStatusText = $state("");
    let showDropdown = $state(false);
    let rooms = $state([]);
    let showAlert = $state(false);
    let alertData = $state({
        icon: "",
        title: "",
        text: "",
        confirmText: "",
        cancelText: "",
        onConfirm: null,
    });
    let showWelcome = $state(true);
    let currentImageBase64 = $state(null);
    let isClosing = $state(false);

    let chatBodyEl;
    let chatBoxEl;
    let chatInputComp;
    const roomKey = "carmen_current_room";
    let closeTimeout = null;
    let userHasScrolledUp = $state(false);

    onMount(async () => {
        await loadRoomList();

        // Restore saved position
        const savedPos = localStorage.getItem(`carmen_chat_pos_${bu}`);
        if (savedPos && chatBoxEl) {
            try {
                const pos = JSON.parse(savedPos);
                chatBoxEl.style.bottom = pos.bottom;
                chatBoxEl.style.right = pos.right;
            } catch (e) {}
        }

        // Track user scroll to detect manual scroll-up
        if (chatBodyEl) {
            // Wheel up = user wants to read older messages
            chatBodyEl.addEventListener(
                "wheel",
                (e) => {
                    if (e.deltaY < 0) {
                        userHasScrolledUp = true;
                    }
                },
                { passive: true },
            );

            // Touch scroll detection
            let touchStartY = 0;
            chatBodyEl.addEventListener(
                "touchstart",
                (e) => {
                    touchStartY = e.touches[0].clientY;
                },
                { passive: true },
            );
            chatBodyEl.addEventListener(
                "touchmove",
                (e) => {
                    if (e.touches[0].clientY > touchStartY) {
                        userHasScrolledUp = true;
                    }
                },
                { passive: true },
            );

            // When user scrolls back to bottom, re-enable auto-scroll
            chatBodyEl.addEventListener(
                "scroll",
                () => {
                    if (isProgrammaticScroll) return;
                    const distanceFromBottom =
                        chatBodyEl.scrollHeight -
                        chatBodyEl.scrollTop -
                        chatBodyEl.clientHeight;
                    if (distanceFromBottom < 80) {
                        userHasScrolledUp = false;
                    }
                },
                { passive: true },
            );

            // Watch for newly added images and re-scroll when they load
            const observer = new MutationObserver((mutations) => {
                for (const mutation of mutations) {
                    for (const node of mutation.addedNodes) {
                        if (node.nodeType === 1) {
                            const el = /** @type {Element} */ (node);
                            const imgs =
                                el.tagName === "IMG"
                                    ? [el]
                                    : el.querySelectorAll("img") || [];
                            for (const img of imgs) {
                                const imgEl = /** @type {HTMLImageElement} */ (
                                    img
                                );
                                if (!imgEl.complete) {
                                    img.addEventListener(
                                        "load",
                                        () => scrollToBottom(),
                                        { once: true },
                                    );
                                }
                            }
                        }
                    }
                }
            });
            observer.observe(chatBodyEl, { childList: true, subtree: true });
        }
    });

    let isProgrammaticScroll = false;
    let autoScrollInterval = null;

    function scrollToBottom(force = false) {
        if (chatBodyEl) {
            tick().then(() => {
                if (force || !userHasScrolledUp) {
                    isProgrammaticScroll = true;
                    chatBodyEl.scrollTop = chatBodyEl.scrollHeight;
                    requestAnimationFrame(() => {
                        isProgrammaticScroll = false;
                    });
                }
            });
        }
    }

    // Start a throttled auto-scroll during streaming (every 300ms)
    function startAutoScroll() {
        if (autoScrollInterval) return;
        autoScrollInterval = setInterval(() => {
            if (!userHasScrolledUp && chatBodyEl) {
                isProgrammaticScroll = true;
                chatBodyEl.scrollTop = chatBodyEl.scrollHeight;
                requestAnimationFrame(() => {
                    isProgrammaticScroll = false;
                });
            }
        }, 50);
    }

    function stopAutoScroll() {
        if (autoScrollInterval) {
            clearInterval(autoScrollInterval);
            autoScrollInterval = null;
        }
    }

    // ==========================================
    // 🚪 Close with Animation
    // ==========================================
    function handleClose() {
        if (isClosing) return;
        isClosing = true;
        showDropdown = false;

        closeTimeout = setTimeout(() => {
            isClosing = false;
            on_close();
        }, 800);
    }

    // ==========================================
    // 🖱️ Draggable Chat Window
    // ==========================================
    function handleDragStart(e, isTouch = false) {
        if (isExpanded || !chatBoxEl) return;
        const event = isTouch ? e.touches[0] : e;
        let startX = event.clientX;
        let startY = event.clientY;

        const style = window.getComputedStyle(chatBoxEl);
        let startBottom = parseInt(style.bottom) || 84;
        let startRight = parseInt(style.right) || 0;
        let isMoving = false;
        let rafId = null;

        const onMove = (eMove) => {
            const ev = eMove.type.startsWith("touch")
                ? eMove.touches[0]
                : eMove;
            const currentX = ev.clientX;
            const currentY = ev.clientY;

            if (
                !isMoving &&
                Math.abs(currentX - startX) < 5 &&
                Math.abs(currentY - startY) < 5
            )
                return;

            if (!isMoving) {
                isMoving = true;
                chatBoxEl.style.top = "auto";
                chatBoxEl.style.left = "auto";
            }

            if (eMove.cancelable) eMove.preventDefault();

            if (rafId) cancelAnimationFrame(rafId);
            rafId = requestAnimationFrame(() => {
                const deltaX = startX - currentX;
                const deltaY = startY - currentY;

                const nextBottom = startBottom + deltaY;
                const nextRight = startRight + deltaX;

                const rect = chatBoxEl.getBoundingClientRect();
                const minB = -22;
                const maxB = window.innerHeight - rect.height - 42;
                const minR = -22;
                const maxR = window.innerWidth - rect.width - 42;

                chatBoxEl.style.bottom =
                    Math.min(Math.max(minB, nextBottom), maxB) + "px";
                chatBoxEl.style.right =
                    Math.min(Math.max(minR, nextRight), maxR) + "px";
            });
        };

        const onEnd = () => {
            window.removeEventListener("mousemove", onMove);
            window.removeEventListener("mouseup", onEnd);
            window.removeEventListener("touchmove", onMove);
            window.removeEventListener("touchend", onEnd);

            if (rafId) cancelAnimationFrame(rafId);

            chatBoxEl.style.transition = "";
            document.body.style.userSelect = "";
            document.body.style.webkitUserSelect = "";

            if (isMoving) {
                localStorage.setItem(
                    `carmen_chat_pos_${bu}`,
                    JSON.stringify({
                        bottom: chatBoxEl.style.bottom,
                        right: chatBoxEl.style.right,
                    }),
                );
            }
        };

        window.addEventListener("mousemove", onMove, { passive: false });
        window.addEventListener("mouseup", onEnd);
        window.addEventListener("touchmove", onMove, { passive: false });
        window.addEventListener("touchend", onEnd);

        chatBoxEl.style.transition = "none";
        document.body.style.userSelect = "none";
        document.body.style.webkitUserSelect = "none";
    }

    // ==========================================
    // 🏠 Room Management
    // ==========================================
    async function loadRoomList() {
        rooms = await api.getRooms();
    }

    async function createNewChat() {
        messageQueue = [];
        if (currentRoomId) {
            await api.clearHistory(currentRoomId);
        }
        currentRoomId = null;
        localStorage.removeItem(roomKey);
        messages = [];
        userHasScrolledUp = false;
        showWelcome = true;
        await loadRoomList();
    }

    async function switchRoom(roomId) {
        messageQueue = [];
        if (currentRoomId === roomId) return;
        currentRoomId = roomId;
        localStorage.setItem(roomKey, roomId);
        showDropdown = false;
        await loadHistory(roomId);
        await loadRoomList();
    }

    async function loadHistory(roomId) {
        const data = await api.getRoomHistory(roomId);
        userHasScrolledUp = false;
        if (data.messages && data.messages.length > 0) {
            messages = data.messages;
            showWelcome = false;
            scrollToBottom(true);
        } else {
            messages = [];
            showWelcome = true;
        }
    }

    function confirmDeleteRoom(roomId) {
        alertData = {
            icon: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>`,
            title: STRINGS.delete_room_confirm_title,
            text: STRINGS.delete_room_confirm_desc,
            confirmText: "ลบทิ้ง",
            cancelText: STRINGS.alert_cancel,
            onConfirm: async () => {
                await api.deleteRoom(roomId);
                showAlert = false;
                if (currentRoomId === roomId) {
                    await createNewChat();
                } else {
                    await loadRoomList();
                }
            },
        };
        showAlert = true;
    }

    // ==========================================
    // 💬 Sending Messages
    // ==========================================
    async function sendMessage(text = null) {
        const msgText = text || inputText.trim();
        if (!msgText && !currentImageBase64) return;

        // Determine if this message will be queued (bot is still busy)
        const willBeQueued = isStreaming || messageQueue.length > 0;

        // Add user message visually immediately to keep UI responsive
        const userMsg = {
            sender: "user",
            message: msgText,
            timestamp: new Date().toISOString(),
            isQueued: willBeQueued,
        };
        messages = [...messages, userMsg];

        // If queued, also add a bot placeholder showing "รอคิว"
        if (willBeQueued) {
            messages = [
                ...messages,
                {
                    id: "queued-" + Date.now(),
                    sender: "bot",
                    message: "",
                    isQueued: true,
                    statusText: "รอคิว...",
                },
            ];
        }

        showWelcome = false;
        inputText = "";
        currentImageBase64 = null;

        // Reset textarea height via Component API if available
        if (chatInputComp) chatInputComp.resetHeight();
        scrollToBottom(true);

        // Ensure room exists
        if (!currentRoomId) {
            const roomTitle =
                msgText.substring(0, 30) + (msgText.length > 30 ? "..." : "");
            const newRoom = await api.createRoom(bu, username, roomTitle);
            currentRoomId = newRoom.room_id;
            localStorage.setItem(roomKey, currentRoomId);
            await loadRoomList();
        }

        // Save user message to database, then queue it for the bot response
        await api.saveMessage(currentRoomId, userMsg);

        // Enqueue the message for processing
        messageQueue.push({ text: msgText });

        // Trigger queue processor
        processQueue();
    }

    async function processQueue() {
        if (isStreaming || messageQueue.length === 0) return;

        isStreaming = true;
        const currentTask = messageQueue.shift();
        const msgText = currentTask.text;

        // Clear "isQueued" flag from the matching user message
        for (let i = 0; i < messages.length; i++) {
            if (
                messages[i].sender === "user" &&
                messages[i].isQueued &&
                messages[i].message === msgText
            ) {
                messages[i].isQueued = false;
                break;
            }
        }

        // Find and replace the queued bot placeholder for this message
        let currentBotIdx = -1;
        for (let i = 0; i < messages.length; i++) {
            if (messages[i].sender === "bot" && messages[i].isQueued) {
                // Convert queued placeholder → streaming placeholder
                messages[i].isQueued = false;
                messages[i].isStreaming = true;
                messages[i].statusText = "กำลังค้นหาเอกสารที่เกี่ยวข้อง";
                messages[i].id = "temp-" + Date.now();
                currentBotIdx = i;
                break;
            }
        }

        // If no queued placeholder found (first message, not queued), add a new one
        if (currentBotIdx === -1) {
            currentBotIdx = messages.length;
            messages = [
                ...messages,
                {
                    id: "temp-" + Date.now(),
                    sender: "bot",
                    message: "",
                    isStreaming: true,
                    statusText: "กำลังค้นหาเอกสารที่เกี่ยวข้อง",
                },
            ];
        }

        streamingText = "";
        streamingStatusText = "กำลังค้นหาเอกสารที่เกี่ยวข้อง";
        scrollToBottom();

        // Status messages timer
        const statusMessages = [
            { delay: 8000, text: "กำลังวิเคราะห์ข้อมูล" },
            { delay: 20000, text: "กำลังเรียบเรียงคำตอบ" },
            { delay: 45000, text: "ใกล้เสร็จแล้ว กรุณารอสักครู่" },
        ];
        const statusTimers = statusMessages.map((msg) =>
            setTimeout(() => {
                if (
                    messages[currentBotIdx] &&
                    messages[currentBotIdx].isStreaming
                ) {
                    messages[currentBotIdx].statusText = msg.text;
                }
            }, msg.delay),
        );

        let fullBotText = "";
        let sourcesData = null;
        let messageId = null;

        // Typing animation state
        let typingBuffer = "";
        let displayedText = "";
        let isTypingAnim = true;

        const processTyping = () => {
            if (typingBuffer.length > 0) {
                const charsToTake =
                    typingBuffer.length > 50
                        ? 5
                        : typingBuffer.length > 10
                          ? 2
                          : 1;
                displayedText += typingBuffer.substring(0, charsToTake);
                typingBuffer = typingBuffer.substring(charsToTake);
                if (messages[currentBotIdx]) {
                    messages[currentBotIdx].message = displayedText;
                }
            }

            if (isTypingAnim || typingBuffer.length > 0) {
                requestAnimationFrame(() => setTimeout(processTyping, 15));
            }
        };

        startAutoScroll();
        processTyping();

        try {
            const h_data = await api.getRoomHistory(currentRoomId);
            const history = h_data ? h_data.messages : [];

            const response = await fetch(`${apiBase}/api/chat/stream`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    text: msgText,
                    bu,
                    username,
                    room_id: currentRoomId,
                    prompt_extend: promptExtend,
                    history,
                }),
            });

            const reader = response.body.getReader();
            const decoder = new TextDecoder("utf-8");

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;

                const chunk = decoder.decode(value, { stream: true });
                const lines = chunk.split("\n");

                for (const line of lines) {
                    if (!line.trim()) continue;
                    try {
                        const json = JSON.parse(line);
                        if (json.type === "chunk") {
                            typingBuffer += json.data;
                            fullBotText += json.data;
                        } else if (json.type === "status") {
                            if (messages[currentBotIdx]) {
                                messages[currentBotIdx].statusText = json.data;
                            }
                        } else if (json.type === "sources") {
                            sourcesData = json.data;
                        } else if (json.type === "done") {
                            messageId = json.id;
                        }
                    } catch (e) {
                        /* ignore parse errors */
                    }
                }
            }
        } catch (error) {
            console.error("Stream Error:", error);
            fullBotText = "⚠️ เกิดข้อผิดพลาดในการเชื่อมต่อ";
            typingBuffer += fullBotText;
            // Mark as error for retry button
            if (messages[currentBotIdx]) {
                messages[currentBotIdx].isError = true;
                messages[currentBotIdx].errorText = msgText;
            }
        }

        // Stop typing animation and auto-scroll
        isTypingAnim = false;
        statusTimers.forEach((t) => clearTimeout(t));
        stopAutoScroll();

        // Wait for typing animation to finish
        await new Promise((resolve) => {
            const checkDone = () => {
                if (typingBuffer.length === 0) {
                    resolve();
                } else {
                    setTimeout(checkDone, 50);
                }
            };
            checkDone();
        });

        // Finalize bot message in-place
        if (messages[currentBotIdx]) {
            messages[currentBotIdx].isStreaming = false;
            messages[currentBotIdx].message = fullBotText;
            messages[currentBotIdx].id = messageId;
            messages[currentBotIdx].sources = sourcesData;
            messages[currentBotIdx].timestamp = new Date().toISOString();
        }

        // Save bot message
        if (messageId && messages[currentBotIdx]) {
            await api.saveMessage(currentRoomId, messages[currentBotIdx]);
        }
        await loadRoomList();
        scrollToBottom();

        // Release streaming lock and process next in queue
        isStreaming = false;
        processQueue();
    }

    function handleSuggestionClick(q) {
        sendMessage(q);
    }

    function handleRetry(errorText) {
        // Remove the error bot message
        messages = messages.filter(
            (m) => !(m.isError && m.errorText === errorText),
        );
        // Re-queue the message
        messageQueue.push({ text: errorText });
        processQueue();
    }

    async function handleFeedback(msgId, score) {
        await api.sendFeedback(msgId, score);
    }

    function handleExpandToggle() {
        if (!isExpanded && chatBoxEl) {
            // Expanding — clear custom position
            chatBoxEl.style.bottom = "";
            chatBoxEl.style.right = "";
        } else if (chatBoxEl) {
            // Collapsing — restore saved position
            const savedPos = localStorage.getItem(`carmen_chat_pos_${bu}`);
            if (savedPos) {
                try {
                    const pos = JSON.parse(savedPos);
                    chatBoxEl.style.bottom = pos.bottom;
                    chatBoxEl.style.right = pos.right;
                } catch (e) {}
            }
        }
        on_toggle_expand();
    }
</script>

<div
    class="chat-box open"
    class:expanded={isExpanded}
    class:closing={isClosing}
    bind:this={chatBoxEl}
>
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        class="chat-main"
        onmousedown={(e) => e.stopPropagation()}
        ontouchstart={(e) => e.stopPropagation()}
    >
        <AlertOverlay bind:showAlert bind:alertData />

        <ChatHeader
            {title}
            {rooms}
            {currentRoomId}
            {isExpanded}
            bind:showDropdown
            on_create_new={createNewChat}
            on_switch_room={switchRoom}
            on_confirm_delete={confirmDeleteRoom}
            on_expand_toggle={handleExpandToggle}
            on_close={handleClose}
            on_drag_start={handleDragStart}
        />

        <div class="chat-body" id="carmenChatBody" bind:this={chatBodyEl}>
            {#if showWelcome && messages.length === 0}
                <WelcomeScreen on_suggestion_click={handleSuggestionClick} />
            {/if}

            {#each messages as msg}
                {#if msg.isQueued && msg.sender === "bot"}
                    <!-- Queued bot placeholder -->
                    <div class="msg bot-msg queued-indicator-container">
                        <div class="queued-status-text">
                            🕐 {msg.statusText || "รอคิว..."}
                        </div>
                    </div>
                {:else if msg.isStreaming && !msg.message}
                    <div class="msg bot-msg typing-indicator-container">
                        <div class="typing-status-text">
                            {msg.statusText || "..."}
                        </div>
                        <div class="typing-dots">
                            <span></span><span></span><span></span>
                        </div>
                    </div>
                {:else}
                    <MessageBubble
                        sender={msg.sender}
                        message={msg.message}
                        msgId={msg.id}
                        sources={msg.sources}
                        {apiBase}
                        isStreaming={msg.isStreaming}
                        isQueued={msg.isQueued}
                        isError={msg.isError || false}
                        timestamp={msg.timestamp}
                        on_feedback={handleFeedback}
                        on_retry={msg.isError
                            ? () => handleRetry(msg.errorText)
                            : null}
                    />
                {/if}
            {/each}

            <!-- Scroll to bottom button -->
            {#if userHasScrolledUp}
                <button
                    class="scroll-to-bottom-btn"
                    onclick={() => scrollToBottom(true)}
                    title="เลื่อนลงล่างสุด"
                >
                    ⬇
                </button>
            {/if}
        </div>

        <ChatInput
            bind:inputText
            bind:this={chatInputComp}
            on_send={() => sendMessage()}
        />
    </div>
</div>

<style>
    /* Chat Window */
    .chat-box {
        position: absolute !important;
        bottom: 84px;
        right: 0;
        width: 360px !important;
        height: 600px !important;
        max-height: 85vh !important;
        background: linear-gradient(
            135deg,
            rgba(255, 255, 255, 0.98) 0%,
            rgba(250, 250, 252, 0.98) 50%,
            rgba(255, 255, 255, 0.98) 100%
        ) !important;
        border-radius: var(--radius-xl) !important;
        box-shadow:
            0 20px 60px rgba(0, 0, 0, 0.15),
            0 8px 24px rgba(0, 0, 0, 0.1),
            0 0 0 1px rgba(0, 0, 0, 0.05) !important;
        border: 1px solid rgba(0, 0, 0, 0.08) !important;
        display: flex !important;
        flex-direction: row !important;
        overflow: hidden !important;
        transform-origin: bottom right !important;
        animation: fadeIn 0.4s ease-out forwards !important;
        transition:
            width 0.6s cubic-bezier(0.34, 1.56, 0.64, 1),
            height 0.6s cubic-bezier(0.34, 1.56, 0.64, 1),
            transform 0.6s cubic-bezier(0.34, 1.56, 0.64, 1),
            opacity 0.6s ease,
            top 0.6s cubic-bezier(0.34, 1.56, 0.64, 1),
            left 0.6s cubic-bezier(0.34, 1.56, 0.64, 1),
            bottom 0.6s cubic-bezier(0.34, 1.56, 0.64, 1),
            right 0.6s cubic-bezier(0.34, 1.56, 0.64, 1),
            border-radius 0.6s ease !important;
        z-index: 1000000 !important;
        backface-visibility: hidden !important;
    }
    .chat-box.closing {
        animation: fadeOut 0.3s ease-out forwards !important;
        pointer-events: none !important;
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: scale(0.95);
            filter: blur(4px);
        }
        to {
            opacity: 1;
            transform: scale(1);
            filter: blur(0);
        }
    }
    @keyframes fadeOut {
        from {
            opacity: 1;
            transform: scale(1);
            filter: blur(0);
        }
        to {
            opacity: 0;
            transform: scale(0.95);
            filter: blur(4px);
        }
    }

    /* Main Content Area */
    .chat-main {
        flex: 1 !important;
        display: flex !important;
        flex-direction: column !important;
        min-width: 0 !important;
        height: 100% !important;
        background: transparent !important;
        position: relative !important;
    }

    /* Chat Body */
    .chat-body {
        flex: 1 !important;
        padding: 24px !important;
        overflow-y: auto !important;
        background: linear-gradient(
            to bottom,
            rgba(252, 252, 254, 0.5) 0%,
            rgba(248, 249, 252, 0.5) 100%
        ) !important;
        background-image: radial-gradient(
            circle at 2px 2px,
            rgba(0, 0, 0, 0.02) 1px,
            transparent 0
        ) !important;
        background-size: 24px 24px !important;
        display: flex !important;
        flex-direction: column !important;
        gap: 20px !important;
        scroll-behavior: smooth !important;
    }

    /* Queued Indicator inside Chat Body */
    .queued-indicator-container {
        padding: 12px 18px !important;
        width: auto !important;
        max-width: 100% !important;
        display: flex !important;
        align-items: center !important;
        box-sizing: border-box !important;
        opacity: 0.6 !important;
    }
    .queued-status-text {
        font-size: 13px !important;
        color: #94a3b8 !important;
        font-family: var(--font-sarabun) !important;
        font-weight: 400 !important;
        font-style: italic !important;
    }

    /* Typing Indicator inside Chat Body */
    .typing-indicator-container {
        padding: 12px 18px !important;
        width: auto !important;
        max-width: 100% !important;
        display: flex !important;
        flex-direction: row !important;
        align-items: center !important;
        gap: 8px !important;
        flex-wrap: wrap !important;
        box-sizing: border-box !important;
    }
    .typing-dots {
        display: flex !important;
        gap: 5px !important;
        align-items: center !important;
        flex-shrink: 0 !important;
    }
    .typing-dots span {
        width: 8px !important;
        height: 8px !important;
        background: var(--primary-color, #34558b) !important;
        border-radius: 50% !important;
        animation: typingBounce 1.4s infinite ease-in-out !important;
        display: block !important;
        opacity: 0.5 !important;
    }
    .typing-dots span:nth-child(1) {
        animation-delay: -0.32s !important;
    }
    .typing-dots span:nth-child(2) {
        animation-delay: -0.16s !important;
    }
    .typing-status-text {
        font-size: 13px !important;
        color: #64748b !important;
        font-family: var(--font-sarabun) !important;
        font-weight: 400 !important;
        white-space: nowrap !important;
        animation:
            statusPulse 2s ease-in-out infinite,
            statusFadeIn 0.3s ease-out !important;
    }
    @keyframes statusPulse {
        0%,
        100% {
            opacity: 1;
        }
        50% {
            opacity: 0.5;
        }
    }
    @keyframes statusFadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }
    @keyframes typingBounce {
        0%,
        80%,
        100% {
            transform: scale(0);
        }
        40% {
            transform: scale(1);
        }
    }

    /* Expanded View */
    .chat-box.expanded {
        position: fixed !important;
        top: 20px !important;
        bottom: 84px !important;
        left: 20px !important;
        right: 20px !important;
        width: calc(100vw - 40px) !important;
        height: auto !important;
        max-height: calc(100vh - 120px) !important;
        border-radius: 24px !important;
        box-shadow: 0 40px 100px rgba(0, 0, 0, 0.5) !important;
    }

    /* Mobile */
    @media (max-width: 768px) {
        :global(:host(:has(.chat-box.open))) {
            top: 0 !important;
            left: 0 !important;
            right: 0 !important;
            bottom: 0 !important;
            width: 100vw !important;
            height: 100vh !important;
        }
        :global(:host(:not(:has(.chat-box.open)))) {
            bottom: 16px !important;
            right: 16px !important;
        }
        :global(:host(:has(.chat-box.open)) .chat-btn) {
            display: none !important;
        }

        .chat-box {
            position: fixed !important;
            top: 0 !important;
            left: 0 !important;
            right: 0 !important;
            bottom: 0 !important;
            width: 100vw !important;
            height: 100vh !important;
            height: 100dvh !important;
            max-height: 100dvh !important;
            max-width: 100vw !important;
            border-radius: 0 !important;
            border: none !important;
            margin: 0 !important;
        }
        #carmenChatBody {
            flex: 1 !important;
            height: auto !important;
            max-height: none !important;
        }
        .chat-main {
            width: 100% !important;
            height: 100% !important;
            padding-bottom: 0 !important;
        }
        .chat-box,
        #carmenChatBody {
            overscroll-behavior-y: none !important;
        }
    }

    /* Scroll to bottom button */
    .scroll-to-bottom-btn {
        position: sticky !important;
        bottom: 12px !important;
        left: 50% !important;
        transform: translateX(-50%) !important;
        width: 36px !important;
        height: 36px !important;
        border-radius: 50% !important;
        background: var(--primary-gradient) !important;
        color: white !important;
        border: none !important;
        padding: 0 !important;
        flex-shrink: 0 !important;
        line-height: 1 !important;
        cursor: pointer !important;
        box-shadow: 0 4px 12px rgba(52, 85, 139, 0.3) !important;
        display: flex !important;
        align-items: center !important;
        justify-content: center !important;
        font-size: 16px !important;
        z-index: 50 !important;
        transition: all 0.2s !important;
        animation: fadeInUp 0.3s ease !important;
        margin: 0 auto !important;
    }
    .scroll-to-bottom-btn:hover {
        transform: translateX(-50%) scale(1.1) !important;
        box-shadow: 0 6px 16px rgba(52, 85, 139, 0.4) !important;
    }
    @keyframes fadeInUp {
        from {
            opacity: 0;
            transform: translateX(-50%) translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateX(-50%) translateY(0);
        }
    }
</style>
