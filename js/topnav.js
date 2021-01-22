const template = document.createElement('template');
template.innerHTML = `
<div class="topnav">
<a href="https://www.bitsmasher.net/bsidesindy/index.html">BSidesIndy</a>
<a href="https://github.com/thedevilsvoice/homelab">HomeLab</a>
<a href="https://www.bitsmasher.net/minecraft/index.html">Minecraft</a>
<a href="https://www.bitsmasher.net/netlab/index.html">NetLab</a>
<a href="https://www.bitsmasher.net/projects/index.html">Projects</a>
<a href="https://www.bitsmasher.net/sara2001conf/index.html">Radio Astronomy</a>
</div>
`

document.body.appendChild(template.content);
