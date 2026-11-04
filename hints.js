(function() {
    if (window.__wimgoHints) {
        window.__wimgoHints.clear();
    }

    const HINT_CHARS = "asdfghjkl";

    function setBadgeText(badge, donePart, restPart) {
        while (badge.firstChild) {
            badge.removeChild(badge.firstChild);
        }
        if (donePart) {
            const doneSpan = document.createElement("span");
            doneSpan.style.color = "#64748b";
            doneSpan.style.fontWeight = "normal";
            doneSpan.textContent = donePart;
            badge.appendChild(doneSpan);
        }
        const restSpan = document.createElement("span");
        restSpan.style.color = "#0f172a";
        restSpan.style.fontWeight = "700";
        restSpan.textContent = restPart;
        badge.appendChild(restSpan);
    }

    window.__wimgoHints = {
        hints: [],
        container: null,

        generateKeys(count) {
            let keys = [""];
            while (keys.length < count) {
                const base = keys.shift();
                for (let i = 0; i < HINT_CHARS.length; i++) {
                    keys.push(base + HINT_CHARS[i]);
                    if (keys.length >= count) break;
                }
            }
            return keys.slice(0, count);
        },

        show() {
            this.clear();

            const selector = [
                'a[href]',
                'button',
                'input:not([type="hidden"])',
                'textarea',
                'select',
                'summary',
                '[role="button"]',
                '[role="link"]',
                '[role="tab"]',
                '[role="menuitem"]',
                '[role="checkbox"]',
                '[role="option"]',
                '[role="listitem"]',
                '[onclick]',
                '[jsaction]',
                '[tabindex]:not([tabindex="-1"])',
                '[contenteditable="true"]'
            ].join(', ');

            const elements = Array.from(document.querySelectorAll(selector));
            const visible = [];
            const vw = window.innerWidth;
            const vh = window.innerHeight;

            for (let el of elements) {
                const rect = el.getBoundingClientRect();
                if (rect.width > 2 && rect.height > 2 &&
                    rect.bottom >= 0 && rect.top <= vh &&
                    rect.right >= 0 && rect.left <= vw) {

                    const style = window.getComputedStyle(el);
                    if (style.visibility !== 'hidden' && style.display !== 'none' && style.opacity !== '0') {
                        visible.push({ el, rect });
                    }
                }
            }

            if (visible.length === 0) return;

            this.container = document.createElement("div");
            this.container.id = "wimgo-hints-container";
            this.container.style.cssText = `
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                z-index: 2147483647;
                pointer-events: none;
            `;

            const keys = this.generateKeys(visible.length);

            this.hints = visible.map((item, i) => {
                const badge = document.createElement("span");
                const hintKey = keys[i];

                badge.style.cssText = `
                    position: fixed;
                    left: ${Math.max(0, item.rect.left)}px;
                    top: ${Math.max(0, item.rect.top)}px;
                    background: rgba(254, 240, 138, 0.82);
                    color: #0f172a;
                    font-family: monospace;
                    font-size: 11px;
                    font-weight: 700;
                    line-height: 12px;
                    padding: 2px 4px;
                    border: 1px solid rgba(15, 23, 42, 0.2);
                    border-radius: 3px;
                    box-shadow: 0 2px 4px rgba(0,0,0,0.25);
                    pointer-events: none;
                    user-select: none;
                    text-transform: lowercase;
                `;

                badge.textContent = hintKey;
                this.container.appendChild(badge);

                return { key: hintKey, el: item.el, badge: badge };
            });

            (document.body || document.documentElement).appendChild(this.container);
        },

        filter(prefix) {
            if (!this.hints.length) return "nomatch";

            const matched = this.hints.filter(h => h.key.startsWith(prefix));

            if (matched.length === 1 && matched[0].key === prefix) {
                const target = matched[0].el;
                this.clear();

                const tag = target.tagName.toLowerCase();
                const isInput = tag === 'input' || tag === 'textarea' || target.isContentEditable;

                target.focus();
                target.click();

                return isInput ? "click:insert" : "click:normal";
            }

            if (matched.length > 0) {
                for (let h of this.hints) {
                    if (h.key.startsWith(prefix)) {
                        h.badge.style.display = "inline-block";
                        const donePart = prefix;
                        const restPart = h.key.slice(prefix.length);
                        setBadgeText(h.badge, donePart, restPart);
                    } else {
                        h.badge.style.display = "none";
                    }
                }
                return "pending";
            }

            this.clear();
            return "nomatch";
        },

        clear() {
            if (this.container && this.container.parentNode) {
                this.container.parentNode.removeChild(this.container);
            }
            this.container = null;
            this.hints = [];
        }
    };

    window.__wimgoHints.show();
})();
