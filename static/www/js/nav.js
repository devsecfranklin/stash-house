/**
 * bitsmasher.net — Dynamic Navigation Bar
 * Injects a modular nav into #site-nav with active-route suppression.
 * Renders: <nav class="w3-container w3-margin-top">
 *            <div class="w3-bar w3-border custom-nav-bar">...</div>
 *          </nav>
 */
(function () {
  "use strict";

  // ── Nav Configuration ────────────────────────────────────────────────
  var NAV_ITEMS = [
    { label: "home",      href: "/",                          external: false },
    { label: "nostr",      href: "/nostr",                         external: false },
    { label: "work",      href: "https://g.dev/franklin",     external: true },
    { label: "minecraft", href: "/minecraft",                 external: false },
    { label: "research",  href: "/research",                  external: false },
    { label: "teaching",  href: "/teaching",                  external: false },
    { label: "training",  href: "/training",                  external: false },
    { label: "discord",   href: "https://discord.gg/mdJGV73Ub", external: true }
  ];

  // ── Route Detection ──────────────────────────────────────────────────
  function currentRoute() {
    var p = window.location.pathname;
    if (p === "/" || p === "") return "/";

    // Normalize: strip trailing slash
    p = p.replace(/\/+$/, "");

    // Direct match against internal nav hrefs
    for (var i = 0; i < NAV_ITEMS.length; i++) {
      var item = NAV_ITEMS[i];
      if (!item.external) {
        if (p === item.href || p.startsWith(item.href + "/")) return item.href;
      }
    }

    // Subdirectory match: "/research/about" -> matches /research
    var parts = p.split("/").filter(Boolean);
    for (var j = 0; j < NAV_ITEMS.length; j++) {
      var candidate = NAV_ITEMS[j];
      if (!candidate.external && parts[0] === candidate.href.replace(/^\//, "")) {
        return candidate.href;
      }
    }

    return null;
  }

  // ── Render ───────────────────────────────────────────────────────────
  function renderNav() {
    var container = document.getElementById("site-nav");
    if (!container) return;

    var current = currentRoute();
    var html = '<nav class="w3-container w3-margin-top">' +
               '<div class="w3-bar w3-border custom-nav-bar">';

    for (var i = 0; i < NAV_ITEMS.length; i++) {
      var item = NAV_ITEMS[i];

      // Suppress active internal route link
      if (!item.external && current === item.href) continue;

      var rel = "";
      if (item.external) rel = ' target="_blank" rel="noopener noreferrer"';

      html += '<a href="' + item.href + '"' + rel +
              ' class="w3-bar-item w3-button">' + item.label + '</a>';
    }

    html += "</div></nav>";
    container.outerHTML = html;
  }

  // Execute on DOM ready (defer script handles this)
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", renderNav);
  } else {
    renderNav();
  }
})();
