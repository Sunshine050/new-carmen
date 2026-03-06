// src/lib/styles.js
// Extracted from old carmen-chatbot-widget/src/core/styles.js
// Returns the CSS text with a dynamic theme color

export function getCssText(themeColor = '#34558b') {
    return `
    :host {
        display: block !important;
        position: fixed !important; 
        bottom: 32px !important; 
        right: 32px !important; 
        z-index: 2000000 !important;
        width: 0 !important; height: 0 !important;
        background: transparent !important;
        pointer-events: none !important;
        --primary-color: ${themeColor};
        --primary-gradient: linear-gradient(135deg, ${themeColor} 0%, ${themeColor} 110%);
        --glass-bg: rgba(255, 255, 255, 0.75);
        --glass-border: rgba(255, 255, 255, 0.5);
        --glass-shadow: 0 12px 48px rgba(0, 0, 0, 0.12);
        --sidebar-bg: #0f172a;
        --sidebar-hover: #1e293b;
        --sidebar-text: #f1f5f9;
        --bg-light: #f8fafc;
        --text-dark: #0f172a;
        --text-gray: #64748b;
        --radius-xl: 32px;
        --radius-lg: 24px;
        --radius-md: 16px;
        --font-inter: 'Inter', system-ui, -apple-system, sans-serif;
        --font-sarabun: 'Sarabun', sans-serif;
        --transition-snappy: all 0.7s cubic-bezier(0.34, 1.56, 0.64, 1);
        --transition-sidebar: width 0.5s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.4s ease;
        font-family: var(--font-inter) !important;
        line-height: 1.5 !important; 
        text-align: left !important; color: var(--text-dark) !important;
    }
    .chat-btn, .chat-box, .chat-tooltip { pointer-events: auto !important; }
    * { box-sizing: border-box !important; }
    ::-webkit-scrollbar { width: 5px !important; height: 5px !important; }
    ::-webkit-scrollbar-track { background: transparent !important; }
    ::-webkit-scrollbar-thumb { background: rgba(0,0,0,0.1) !important; border-radius: 10px !important; transition: background 0.3s !important; }
    *:hover::-webkit-scrollbar-thumb { background: rgba(0,0,0,0.2) !important; }
    `;
}

