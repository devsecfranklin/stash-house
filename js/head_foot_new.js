class NavBar extends HTMLElement {
    connectedCallback() {
        this.innerHTML = `
        <div class="topnav">
            <a href="/bsidesindy/index.html">BSidesIndy</a>
            <a href="/lab/index.html">Lab</a>
            <a href="https://minecraft.bitsmasher.net/">Minecraft</a>
            <a href="https://franklin-resume.herokuapp.com/">My Resume</a>
        </div>`
    }
}

//Footer

class Footer extends HTMLElement {
    connectedCallback() {
        var modiDate = new Date(document.lastModified);
        var showAs = (modiDate.getMonth() + 1) + "/" + modiDate.getDate() + "/" + modiDate.getFullYear();
        var modiDate = new Date();
        var Seconds

        if (modiDate.getSeconds() < 10) {
            Seconds = "0" + modiDate.getSeconds();
        } else {
            Seconds = modiDate.getSeconds();
        }

        var modiDate = new Date();
        var CurTime = modiDate.getHours() + ":" + modiDate.getMinutes() + ":" + Seconds

        this.innerHTML = `
    <!-- Footer -->
    <div class="footer">
        Copyright © 2010-2022 All Rights Reserved.<br />
        Last updated on ` + showAs + ` at ` + CurTime + `
    </div>
    <!-- Footer -->
        `;
    }
}

customElements.define('nav-bar', NavBar);
customElements.define('main-footer', Footer);

