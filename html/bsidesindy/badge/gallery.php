<html>
<head>
<style type="text/css">
.thumb{
    border:1px solid #CCC;
    padding:5px;
    background:#FFF;
    box-shadow:0px 0px 5px #DDD;
    -webkit-box-shadow:0px 0px 5px #DDD;
    -moz-box-shadow:0px 0px 5px #DDD;   
}
</style>
</head>
<body>
<table border=1>
<?php
// define directory path
$dir = ".";

// iterate through files
// look for JPEGs
if (is_dir($dir)) {
    if ($dh = opendir($dir)) {
        while (($file = readdir($dh)) !== false) {
            if (preg_match("/.jpg/", $file)) {
                // read EXIF headers
                $exif = exif_read_data($file, 0, true);
                echo "<tr>";
                // get thumbnail
                // link to full image
                echo "<td valign=top><a href=$file><img src=thumbnail.php?img=$file&type=square&align=top class=thumb></a><td>";
                //echo "<td valign=top><font size=-1>";
                // get file name
                echo "File: <b>" . $exif['FILE']['FileName'] . "</b><br/>";
                // get timestamp
                echo "Timestamp: "  . $exif['IFD0']['DateTime'] . "<br/>";
                // get image dimensions
                echo "Dimensions: " . $exif['COMPUTED']['Height'] . " x " . $exif['COMPUTED']['Height'] . " <br/>";
                // get camera make and model
                echo "Camera: " . $exif['IFD0']['Model'];
                echo "</font></td>";
                echo "</tr>";
           }
        }
       closedir($dh);
   }
}
?>
</table>
</body>
</html>
