<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">

<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>Networking Lab @bitsmasher.net</title>
<style type="text/css">
body,td,th {
	font-family: Verdana, Geneva, sans-serif;
	font-size: 12px;
	color: #FFF;
}
body {
	background-color: #333;
}
</style>
</head>

<body>
<table width="100%" border="0">
  <tr>
    <td><h1>Networking Lab @bitsmasher.net</h1></td>
    <td> <div align="right">[ <a href="http://www.bitsmasher.net/">Home</a> ] [ <a href="http://www.bitsmasher.net/projects/">projects</a> ] [ <a href="http://www.bitsmasher.net/about.html">about</a> ] [ <a href="https://server14.websitehostserver.net:2083/xferwebmail/">webmail</a> ]</div></td>
  </tr>
</table>


<p>&nbsp;</p>

<table border="0">
  <tr>
    <td valign=top><p>I maintain a small collection of routers and switches that I can use to study various routing and switching configurations. Also I experiment here with botnets, packet shapers, security tools like Snort, Nessus, and so on. </p>
      <p>This lab is isolated from the rest of my network, and also from the Internet so I can pretty much do anything I want without concern for running a vulnerable service or doing something that might be considered &quot;illegal&quot;. Pretty much protecting myself from the bad guys <i>and</i> the good guys by only allowing traffic to come in, but no traffic to go out. </p>
    <p>I started a  collection of practice  <a href="labs">labs</a> but have since found a lot of good sites on the web that already have better collections than I have time to make myself. I will keep dumping random labs in here anyway. </p>
    <h2><u>Tools</u> <a href="tools">(My tools are stored here)</a></h2>
    <p><u>ettcp</u> - ettcp is based on the venerable ttcp application.ttcp allows measurement of   network throughputover TCP or UDP, between two nodes. ettcp adds several useful   features to ttcp, whileretaining bckwards compability. </p>
    <p><u>ping+</u> - I got this tool from Dr. James Yu at DePaul University. It is free to use but please leave the names of the authors in the headers if you copy it.</p>
	<p><u>perl/expect scripts</u> - PERL is a classic language for UNIX admins to automate their tasks and I find it extremely useful when the expect module is added to it. 
I have written (and lost) quite a few of these scripts over the years. For example, <a href="tools/bb1.pl">here is a script that allows a log-in to a single device</a> without having to do all the tedious intermediate steps.
Cisco includes TCL ("tickle") as of 12.3T, <a href="http://www.cisco.com/en/US/docs/ios/12_3t/12_3t2/feature/guide/gt_tcl.html">you can read all about that here.</a>
<p><u>MC Hammer</u> - Love the name, tool from Nortel that allows user to generate Multicast traffic. The link for download is 
<a href="http://www.nortel.com/corporate/nortel_on_nortel/multicast_hammer.html">Multicast Hammer</a>
<p>
Other tools: tacacs+, ettercap



<h2><u>Equipment</u></h2>
<p> <br />
<u>netlab1</u>
<p> <br />

I keep a 1u Sun Sparc v100 running OpenBSD to serve various lab functions like Tacacs server, NTP server, TFTP server for backing up IOS images and router configurations, syslog server, etc. 
<p> <br />
<u>remote power</u>
<p> <br />
You can remotely power the routers and switches up and down using the 
<a href="http://10.10.6.9/">web interface here.</a> The login is "device" and
the password is "apc". 
<p> <br />
<u>Cisco routers</u>
<p> <br />
Some of my router <a href="configs">conifgurations are stored here</a>. Most of them are stored in the /tftpboot directory on the netlab1 host. I've tried a
few times to document the connections between devices in a Visio drawing, but
things tend to change so often that they're quickly obsolete. My plan is to 
eventually get the routers and switches set up in a configuration that matches
one of those popular lab "workbooks" and just leave it that way.  
<p> <br />
<table border="1">
<tr>
<th>Model</th><th>hostname</th><th>DRAM</th><th>Flash</th> <th>IOS Version</th>
<th> local sw port</th><th>remote sw port</th><th>serial #</th>
</tr>
<tr><td>2610</td><td>R1</td><td>64MB</td><td>16MB</td><td>(C2600-IK9O3S3-M), Version 12.3(26)</td><td>e0/0</td><td>11</td><td>9</td></tr>
<tr><td>2621</td><td>R2</td><td>64MB</td><td>8MB</td><td>(C2600-I-M), Version 12.1(4)</td><td>fa0/0,fa0/1</td><td>5,18</td><td>10</td></tr>
<tr><td>2621</td><td>R3</td><td>64MB</td><td>8MB</td><td>(C2600-I-M), Version 12.1(4)</td><td>fa0/0,fa0/1</td><td>8,10</td><td>11</td></tr>
<tr><td>2621</td><td>R4</td><td>64MB</td><td>16MB</td><td>(C2600-DS-M), Version 12.1(2)</td><td>fa0/0,fa0/1</td><td>14,16</td><td>12</td></tr>
<tr><td>2620</td><td>R5</td><td>48MB</td><td>16MB</td><td>(C2600-DS-M), Version 12.1(5)T10</td><td>fa0/0</td><td>12</td><td>13</td></tr>
<tr><td>2621</td><td>R6</td><td>64MB</td><td>16MB</td><td></td><td></td><td></td><td></td></tr>
<tr><td>2501</td><td>BB1</td><td>16MB</td><td>16MB</td><td></td><td></td><td>2</td><td>15</td></tr>
<tr><td>2501</td><td>BB2</td><td>16MB</td><td>16MB</td><td></td><td></td><td>9</td><td>16</td></tr>
<tr><td>2522</td><td>R9 (frame_switch)</td><td>16MB</td><td>16MB</td><td></td><td>e0</td><td>4</td><td>14</td></tr>
</table>
<p>
To connect to the console of one of the devices, look at the number in the "serial #" column. Telnet to labgate and then select portX, where X is the number from the column. For example, to connect to the console port of R1, telnet to labgate and then type "port9" and you will be connected to the console. 
<p>
Also remember that you can Suspend the Telnet session by entering Ctrl-Shift-6 x So easy to use and remember, thanks Cisco! =]
<p>
<h2><u>Books</u></h2>
<p> <br />
  <a href="http://www.informit.com/library/library.aspx?b=CCIE_Practical_Studies_I">CCIE Practical Studies Vol. 1</a></p>
<p>&nbsp;</p>
<h2><u>Network Diagrams</u></h2>

<table width="100%" border="0">
  <tr>
    <td><a href="diagrams/planning-aug-11.png" target="_top"><img src="diagrams/planning-aug-11-thumb.png" alt="network" /></a></td>
<td><p>An updated view of the lab network that includes the new PIX and remote power controller. I think this one is a bit easier to read, though it does not show any serial connections.</p></td> 
  </tr>
</table>
<p></p>

<table width="100%" border="0">
  <tr>
    <td><a href="diagrams/netlab1-serial.png" target="_top"><img src="diagrams/netlab1-serial-thumb.png" alt="network" /></a></td>
<td><p>This is the current configuration for serial connections. It's here so it can be printed out and used to do labs.</p></td>
  </tr>
</table>
<p></p>

<table width="100%" border="0">
  <tr>
    <td><a href="diagrams/netlab1-ethernet.png" target="_top"><img src="diagrams/netlab1-ethernet-thumb.png" alt="network" /></a></td>
<td><p>This is the current configuration for ethernet connections. It's here so it can be printed out and used to do labs.</p></td>
  </tr>
</table>
<p></p>

<table width="100%" border="0">
  <tr>
    <td><a href="diagrams/planning-mar-10.jpg" target="_top"><img src="diagrams/planning-mar-10-thumbnail.jpg" width="100" height="77" alt="network" /></a></td>
    <td><p>Here is a Visio diagram I did earlier this year that shows my whole network. This diagram does not show &quot;everything&quot; believe it or not since things are always changing around and you can't expect to track everything down to the last laptop and game machine these days when everything has an IP address, but it does hit the highlights. If I find the time I may do a smaller diagram that shows just the Network lab components.</p></td>
  </tr>
</table>
<p></p>

<table width="100%" border="0">
<tr>
 <td><a href="images/dl_pod_DePaul.jpg" target="_top"><img src="images/dl_pod_DePaul-thumb.jpg" width="115" height="156" alt="DL Pod at DePaul" /></a></td>
 <td>This is a typical "Distance Learning Pod" at DePaul University in Chicago. Pod is the term they use for a single lab rack, I think they have 4 or 5 of these. Lucky me I get to spend all kinds of time doing labs on these racks =]</td>
</tr>
</table>
<p></p>




<h2><u>Links</u></h2>
<a href="http://iosadventures.blogspot.com/">IOS Adventures</a>
<br />
<a href="http://blog.ine.com/">INE Blog</a>
<br />
<a href="http://www.bitsmasher.net/netlab/books/76.pdf">How to Download a Software Image to a Cisco 2600 via TFTP Using the tftpdnld ROMMON Command</a>
<br />
<a href="http://www.grid.unina.it/software/ITG/sdescr.php">D-ITG, Distributed Internet Traffic Generator</a>




<p>&nbsp;</p></td>
    <td align=right>
<!-- IPv6-test.com widget BEGIN -->
<script type="text/javascript">var _ipv6test_widget_style = {
border: "solid 1px #000",
font_size: "12px",
show_country_flags: true,
show_loading_anim: true,
ipv4_label_color: "#393",
ipv4_background_color: "#eee",
ipv6_label_color: "#339",
ipv6_background_color: "#ddd",
stats_position: "bottom",
stats_font_size: "10px",
stats_color: "#eee",
stats_color_v4: "#beb",
stats_color_v6: "#bbe",
stats_background_color: "#666"
}</script>

<div id="_ipv6test_widget" style="width:250px;display:none">loading <a href="http://ipv6-test.com/">IPv6 connection test</a> ...</div><script type="text/javascript" src="http://ipv6-test.com/api/widget.php?domain=referer" async="async"></script>
<!-- IPv6-test.com widget END -->
<P>
<img src="images/lab_rack.jpg" alt="my lab rack" width="249" height="533" align="right" />

<script type="text/javascript" language="javascript" src="http://ipv6.he.net/certification/badge.js"></script>
<script type="text/javascript">/*
<![CDATA[*/
var user = "bitsmasher";
display_swf(user);
//]]></script>
</td>
  </tr>
</table>
<p>&nbsp;</p>
<h2>&nbsp;</h2>
<p>&nbsp;</p>

<table width="100%" border="0">
  <tr>
    <td><h1>Networking Lab @bitsmasher.net</h1></td>
    <td> <div align="right">[ <a href="http://www.bitsmasher.net/">Home</a> ] [ <a href="http://www.bitsmasher.net/fubaria/">#fubaria</a> ] [ <a href="http://www.bitsmasher.net/school/">teaching</a> ] [ <a href="http://www.bitsmasher.net/projects">projects</a> ] [ <a href="http://www.bitsmasher.net/about.html">about</a> ] [ <a href="http://mail.bitsmasher.net/webmail/">webmail</a> ]</div></td>
  </tr>
</table>

	<!-- end page -->
	<div id="footer">
		<p id="legal">Copyright © 2010,2011 All Rights Reserved. </p>
		<p id="links"><a href="#">Privacy Policy</a> | <a href="#">Terms of Use</a></p>
	</div>
	<!-- end footer -->
</body>

</html>
