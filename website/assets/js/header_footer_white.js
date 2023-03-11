//Header
class Header extends HTMLElement {
  connectedCallback() {
    this.innerHTML = `
    <div class="header-white">
        <a href="/index.html"> <img src="/images/bitsmasherflag.gif" alt="The TLD over an American flag."></a>
    </div>
    <div class="topnav">
        <a href="https://www.bitsmasher.net/">Minecraft</a>
    </div>
                        `;
  }
}
//Footer
class Footer extends HTMLElement {
  connectedCallback() {
    var modiDate = new Date(document.lastModified);
    var showAs =
      modiDate.getMonth() +
      1 +
      "/" +
      modiDate.getDate() +
      "/" +
      modiDate.getFullYear();
    var modiDate = new Date();
    var Seconds;

    if (modiDate.getSeconds() < 10) {
      Seconds = "0" + modiDate.getSeconds();
    } else {
      Seconds = modiDate.getSeconds();
    }

    var modiDate = new Date();
    var CurTime =
      modiDate.getHours() + ":" + modiDate.getMinutes() + ":" + Seconds;

    this.innerHTML =
      `
    <!-- Footer -->
    <div class="footer-white">
        Copyright © 2010-2023 All Rights Reserved.<br />
        Last updated on ` +
      showAs +
      ` at ` +
      CurTime +
      `
    </div>
    <!-- Footer -->
        `;
  }
}

customElements.define("main-header", Header);
customElements.define("footer-white", Footer);
