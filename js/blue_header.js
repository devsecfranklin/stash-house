class Header extends HTMLElement {
    connectedCallback() {
        this.innerHTML = `
    <div class="header">
        <a href="/index.html"> <img src="/images/bit-logo2.jpg" alt="a blue rectangle with the URI of the website"></a>
    </div>
    <div class="topnav">
        <a href="/bsidesindy/index.html">BSidesIndy</a>
        <a href="https://github.com/thedevilsvoice/homelab">HomeLab</a>
        <a href="/minecraft/index.html">Minecraft</a>
        <a href="/netlab/index.html">NetLab</a>
        <a href="/projects/index.html">Projects</a>
        <a href="/sara2001conf/index.html">Radio Astronomy</a>
    </div>
                        `
    }
}

//Footer

class Footer extends HTMLElement {
    connectedCallback() {
        this.innerHTML = `
    <!-- Footer -->
    <div class="footer">
        Copyright © 2010-2021 All Rights Reserved.<br />
        <script type="text/javascript" src="/js/last.js"></script>
    </div>
    <!-- Footer -->
        `;
    }
}

customElements.define('main-header', Header);
customElements.define('main-footer', Footer);
