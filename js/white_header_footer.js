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
        Copyright © 2010-2021 All Rights Reserved.<br />
        Last updated on ` + showAs + ` at ` + CurTime + `
    </div>
    <!-- Footer -->
        `;
    }
}

// customElements.define('main-header', Header);
customElements.define('white-footer', Footer);
