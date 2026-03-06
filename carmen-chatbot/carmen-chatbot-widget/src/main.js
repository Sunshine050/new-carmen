import App from './App.svelte';
import { mount } from 'svelte';

// ==========================================
// 🚀 Auto-init: Script Tag Usage
// ==========================================
// <script src="carmen-widget.js" bu="HR" user="Guest" api-base="http://localhost:8000"></script>
const script = document.currentScript || document.querySelector('script[src*="carmen-widget"]');
if (script && (script.hasAttribute('bu') || script.hasAttribute('data-bu'))) {
  const config = {
    bu: script.getAttribute('bu') || script.getAttribute('data-bu'),
    user: script.getAttribute('user') || script.getAttribute('data-user'),
    apiBase: script.getAttribute('api-base') || script.getAttribute('data-api-base'),
    theme: script.getAttribute('theme') || script.getAttribute('data-theme'),
    title: script.getAttribute('title') || script.getAttribute('data-title'),
    promptExtend: script.getAttribute('prompt-extend') || script.getAttribute('data-prompt-extend'),
  };
  initCarmen(config);
}

// ==========================================
// 🏗️ Programmatic Init
// ==========================================
// window.CarmenBot({ bu: 'HR', user: 'Admin', apiBase: 'http://...' })
function initCarmen(config) {
  // Create a host element with Shadow DOM
  const host = document.createElement('div');
  host.id = 'carmen-widget-host';
  document.body.appendChild(host);
  const shadow = host.attachShadow({ mode: 'open' });

  // Create a mount target inside shadow
  const mountTarget = document.createElement('div');
  shadow.appendChild(mountTarget);

  mount(App, {
    target: mountTarget,
    props: { config, shadowRoot: shadow }
  });
}

/** @type {any} */ (window).CarmenBot = initCarmen;

// ==========================================
// 🧩 Custom Element
// ==========================================
if ('customElements' in window) {
  class CarmenChatbotElement extends HTMLElement {
    connectedCallback() {
      const bu = this.getAttribute('bu');
      const user = this.getAttribute('user');
      const apiBase = this.getAttribute('api-base');
      if (bu && user && apiBase) {
        initCarmen({
          bu, user, apiBase,
          theme: this.getAttribute('theme'),
          title: this.getAttribute('title'),
        });
      }
    }
  }
  customElements.define('carmen-chatbot', CarmenChatbotElement);
}
