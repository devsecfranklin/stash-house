class NavBar extends HTMLElement {
  connectedCallback() {
    this.innerHTML = `
        <div class="topnav">
            <a href="https://franklin-resume.herokuapp.com/" target="_blank">My Resume</a>
            <a href="https://www.bitsmasher.net/minecraft">Minecraft</a>
            <a href="/index.html">Home</a>
        </div>`;
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
    <div class="footer">
    <table><tr><td>
        Copyright © 2010-2023 All Rights Reserved.<br />
        Last updated on ` +
      showAs +
      ` at ` +
      CurTime +
      `</td><td>
      <a href='http://ipv6-test.com'><img src='http://v4v6.ipv6-test.com/imgtest.php?bl=1' alt='ipv6 test'
                title='ipv6 test' border='0' /></a></td>
    </tr></table></div>
        `;
  }
}

customElements.define("nav-bar", NavBar);
customElements.define("main-footer", Footer);
