<?php
 
// thumb width
$square = 150;
$large = 200;
$small = 100;
 
////////////////////////////////////////////////////////////////////////////////// square
if( isset($_GET[img]) &&  ( $_GET[type] == "square" || $_GET[type] == "" ) ){
 
// thumb size
$thumb_width = $square;
$thumb_height = $square;
 
// align
$align = $_GET[align];
 
// image source
$imgSrc = $_GET[img];
$imgExt = substr($imgSrc,-3);
 
// image extension
if($imgExt == "jpg"){ $myImage = imagecreatefromjpeg($imgSrc); }
if($imgExt == "gif"){ $myImage = imagecreatefromgif($imgSrc); }
if($imgExt == "png"){ $myImage = imagecreatefrompng($imgSrc); }
 
// getting the image dimensions  
list($width_orig, $height_orig) = getimagesize($imgSrc);   
 
// ratio
$ratio_orig = $width_orig/$height_orig;
 
// landscape or portrait?
if ($thumb_width/$thumb_height > $ratio_orig) {
   $new_height = $thumb_width/$ratio_orig;
   $new_width = $thumb_width;
} else {
   $new_width = $thumb_height*$ratio_orig;
   $new_height = $thumb_height;
}
 
// middle
$x_mid = $new_width/2;
$y_mid = $new_height/2;
 
// create new image
$process = imagecreatetruecolor(round($new_width), round($new_height)); 
imagecopyresampled($process, $myImage, 0, 0, 0, 0, $new_width, $new_height, $width_orig, $height_orig);
$thumb = imagecreatetruecolor($thumb_width, $thumb_height); 
 
// alignment
if($align == ""){
imagecopyresampled($thumb, $process, 0, 0, ($x_mid-($thumb_width/2)), ($y_mid-($thumb_height/2)), $thumb_width, $thumb_height, $thumb_width, $thumb_height);
}
if($align == "top"){
imagecopyresampled($thumb, $process, 0, 0, ($x_mid-($thumb_width/2)), 0, $thumb_width, $thumb_height, $thumb_width, $thumb_height);
}
if($align == "bottom"){
imagecopyresampled($thumb, $process, 0, 0, ($x_mid-($thumb_width/2)), ($new_height-$thumb_height), $thumb_width, $thumb_height, $thumb_width, $thumb_height);
}
if($align == "left"){
imagecopyresampled($thumb, $process, 0, 0, 0, ($y_mid-($thumb_height/2)), $thumb_width, $thumb_height, $thumb_width, $thumb_height);
}
if($align == "right"){
imagecopyresampled($thumb, $process, 0, 0, ($new_width-$thumb_width), ($y_mid-($thumb_height/2)), $thumb_width, $thumb_height, $thumb_width, $thumb_height);
}
 
imagedestroy($process);
imagedestroy($myImage); 
 
if($imgExt == "jpg"){ imagejpeg($thumb, null, 100); }
if($imgExt == "gif"){ imagegif($thumb); }
if($imgExt == "png"){ imagepng($thumb, null, 9); }
 
}
 
////////////////////////////////////////////////////////////////////////////////// normal
if( isset($_GET[img]) && ( $_GET[type] == "large" || $_GET[type] == "small" )  ){
 
if( $_GET[type] == "large" ){ $thumb_width = $large; }
if( $_GET[type] == "small" ){ $thumb_width = $small; }
 
// image source
$imgSrc = $_GET[img];
$imgExt = substr($imgSrc,-3);
 
// image extension
if($imgExt == "jpg"){ $myImage = imagecreatefromjpeg($imgSrc); }
if($imgExt == "gif"){ $myImage = imagecreatefromgif($imgSrc); }
if($imgExt == "png"){ $myImage = imagecreatefrompng($imgSrc); }
 
//getting the image dimensions  
list($width_orig, $height_orig) = getimagesize($imgSrc);   
 
// ratio
$ratio_orig = $width_orig/$height_orig;
$thumb_height = $thumb_width/$ratio_orig;
 
// new dimensions
$new_width = $thumb_width;
$new_height = $thumb_height;
 
// middle
$x_mid = $new_width/2;
$y_mid = $new_height/2;
 
// create new image
$process = imagecreatetruecolor(round($new_width), round($new_height));
 
imagecopyresampled($process, $myImage, 0, 0, 0, 0, $new_width, $new_height, $width_orig, $height_orig);
$thumb = imagecreatetruecolor($thumb_width, $thumb_height); 
imagecopyresampled($thumb, $process, 0, 0, ($x_mid-($thumb_width/2)), ($y_mid-($thumb_height/2)), $thumb_width, $thumb_height, $thumb_width, $thumb_height);
 
if($imgExt == "jpg"){ imagejpeg($thumb, null, 100); }
if($imgExt == "gif"){ imagegif($thumb); }
if($imgExt == "png"){ imagepng($thumb, null, 9); }
 
}
?>

