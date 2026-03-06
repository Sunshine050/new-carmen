# Carmen Chatbot Widget (Svelte 5)

A modern, embeddable Chatbot Widget for the Carmen project, rebuilt from the ground up using **Svelte 5** and **Vite**.

This widget compiles down to a **single JavaScript file** (`carmen-widget.js`) with fully encapsulated CSS (via Shadow DOM), ensuring it never conflicts with the styling of the host website.

## ✨ Features

- **Svelte 5 + Vite**: Fast, modern frontend with Svelte's Runes reactivity system.
- **Single File Output**: Everything (JS + CSS) is bundled into one lightweight script.
- **Shadow DOM**: Fully isolated styling.
- **Draggable Window**: Users can move the chat window around the screen.
- **Smooth Animations**: High-quality opening, closing, and streaming animations.
- **Streaming Response**: Real-time Markdown rendering for bot responses.
- **Feedback System**: Built-in UI for thumbs up/down feedback on messages.

---

## 🚀 Setup & Development

### 1. Install Dependencies

Make sure you have Node.js installed, then run:

```bash
npm install
```

### 2. Development Server (UI Testing)

To test the widget layout locally without compiling to a single file:

```bash
npm run dev
```

### 3. Build for Production

To build the final widget script:

```bash
npm run build
```

*(Note: A custom Vite plugin in `vite.config.js` will automatically copy the built `carmen-widget.js` to `../../frontend/user/public/carmen-widget.js` for seamless integration).*

---

## 💻 How to Embed

To embed the chatbot on any HTML page, simply include the script tag with the required `data-*` attributes:

```html
<script 
    src="/carmen-widget.js" 
    bu="YOUR_BU"
    user="YOUR_USERNAME"
    title="ผู้ช่วยส่วนตัว (Carmen)"
    api-base="http://localhost:8000" 
    theme="#34558b" 
    defer>
</script>
```

### Configuration Attributes

| Attribute | Description | Default |
| --- | --- | --- |
| `bu` | Business Unit ID (Required) | - |
| `user` | Username or Employee ID (Required) | - |
| `api-base` | Backend FastAPI URL | `http://localhost:8000` |
| `theme` | Primary Theme Color (Hex) | `#34558b` |
| `title` | Widget Header Title | `Carmen Chatbot` |
| `prompt-extend` | Additional hidden context | `""` |

---

## 📁 Project Structure

- **`src/main.js`**: The main entry point. Reads `data-*` attributes from the script tag and mounts the Svelte app into a Shadow DOM.
- **`src/App.svelte`**: The root component. Manages the overall state (`requestClose`, `isOpen`) and toggles between the launcher and the main chat window.
- **`src/components/`**:
  - `ChatWindow.svelte`: The main chat interface, handles drag-and-drop, API streaming, and conversation history.
  - `MessageBubble.svelte`: Renders individual messages (user/bot), Markdown parsing, copy functionality, and feedback buttons.
  - `Launcher.svelte`: The floating action button that opens the chat.
  - `WelcomeScreen.svelte`: The initial greeting and suggested questions view.
- **`src/lib/`**:
  - `api-client.js`: Handles all HTTP communication with the FastAPI backend.
  - `styles.js`: Contains all the raw CSS injected into the Shadow DOM.
  - `icons.js`: SVG icons used throughout the widget.
  - `constants.js`: Strings, configuration, and suggested questions.
  - `content-formatter.js`: Utility for parsing Markdown and formatting responses.
