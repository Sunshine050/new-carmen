import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import cssInjectedByJsPlugin from 'vite-plugin-css-injected-by-js';
import { resolve } from 'path';
import { copyFileSync, existsSync, mkdirSync } from 'fs';

// Post-build: copy to frontend/user/public/
function copyWidgetPlugin() {
  const targets = [
    resolve(__dirname, '../../frontend/user/public'),
  ];

  return {
    name: 'copy-widget',
    closeBundle() {
      const src = resolve(__dirname, 'dist/carmen-widget.js');
      if (!existsSync(src)) return;

      for (const dir of targets) {
        if (!existsSync(dir)) mkdirSync(dir, { recursive: true });
        const dest = resolve(dir, 'carmen-widget.js');
        copyFileSync(src, dest);
        console.log(`\n📦 Copied carmen-widget.js → ${dest}`);
      }
    },
  };
}

export default defineConfig({
  plugins: [
    svelte({
      emitCss: true,
    }),
    cssInjectedByJsPlugin({
      injectCodeFunction: function (cssCode) {
        try {
          // Try to inject into Shadow DOM first (our widget's host)
          const host = document.getElementById('carmen-widget-host');
          if (host && host.shadowRoot) {
            const style = document.createElement('style');
            style.textContent = cssCode;
            host.shadowRoot.appendChild(style);
            return;
          }
          // Fallback: wait for the host element to appear
          const observer = new MutationObserver(function (mutations, obs) {
            const el = document.getElementById('carmen-widget-host');
            if (el && el.shadowRoot) {
              const style = document.createElement('style');
              style.textContent = cssCode;
              el.shadowRoot.appendChild(style);
              obs.disconnect();
            }
          });
          observer.observe(document.body, { childList: true, subtree: true });
        } catch (e) {
          // Ultimate fallback: inject into document head
          const style = document.createElement('style');
          style.textContent = cssCode;
          document.head.appendChild(style);
        }
      }
    }),
    copyWidgetPlugin(),
  ],
  build: {
    lib: {
      entry: 'src/main.js',
      name: 'CarmenWidget',
      formats: ['iife'],
      fileName: () => 'carmen-widget.js',
    },
    cssCodeSplit: false,
    rollupOptions: {
      output: {
        assetFileNames: '[name][extname]',
      },
    },
  },
});
