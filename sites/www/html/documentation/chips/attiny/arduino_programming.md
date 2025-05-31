# Arduino Programming

* Connect a 10 uF capacitor between “Reset” and “Ground” pins on the Arduino Uno board as shown in the diagram. The stripe of the capacitor with the negative sign (“-”) goes to the Arduino “Ground” pin. This prevents the Arduino Uno from resetting and ensures that the Arduino IDE talks to the ArduinoISP (and not the bootloader) during the upload of the sketches to the ATTiny.
* Connect the ATtiny84/44 Pin 1 (with the little dot) to the 5 volt breadboard rail.
* Connect the ATtiny84/44 Pin 14 to ground.
* RESET: Connect the ATtiny84/44 Pin 4 (Reset) to Arduino Pin 10.
* MOSI: Connect the ATtiny84/44 Pin 7  to Arduino Pin 11.
* MISO: Connect the ATtiny84/44 Pin 8 to Arduino Pin 12.
* CLOCK: Connect the ATTiny84/44 Pin 9 to Arduino Pin 13
* Connect the ground and 5v from the Arduino to the breadboard.

Make sure you select the option “ATTiny84  (internal 1 MHz clock)”, or “ATTiny44  (internal 1 MHz clock)” in the Tools->Board menu (depending on which version of the chip you have):

Upload the sketch using the “File” -> “Upload Using Programmer” command (Ctrl + Shift + U).

If you get an error like this: “avrdude: please define PAGEL and BS2 signals in the configuration file…”, you can safely ignore it. This relate to parallel programmers and is not an issue here.

If all works well, you should have a blinking LED and a proven method for uploading Arduino sketches to your ATTiny84/44. Congratulations!

The next step  (Step 6) is optional. By default, the ATTiny chips run at 1 MHz. You need an extra step to configure the microcontroller to run at 8 MHz. This is a requirement for using some popular Arduino libraries, as well as to take full advantage of your chip.

To do so, you need to use the same wiring that you used to upload the blink sketch in the previous step. Under “Tools”->”Board” menu, this time select the option “ATtiny84 (internal 8 Mhz clock)”, or  “ATtiny44 (internal 8 Mhz clock)”, depending on which chip you have. It is important to pick the correct clock speed from the menu. Choosing the 20 Mhz external clock option, will require you to add wire in a 20 MHz crystal, or resonator, in order to program and use your ATTiny84/44.

Next, again under the “Tools” menu, choose the option “Burn Bootloader”. There is actually, no bootloader on the ATTiny chips, but this step will set the fuses to the values that will allow it to run at 8 Mhz, similarly to what we did for the ATtiny85 in the previous write-up.

Again, ignore any errors related to PAGEL and BS2 signals you may get. To make sure the change worked, you may want to re-upload the modified Blink sketch from step 5, this time using the internal 8 MHz clock option for your chip from the Boards menu.

Next, time to use the ATTiny in a real project! A warning: not all Arduino libraries will work on the ATTiny chips. The Arduino forum is a great place to search for additional info, if you run into issues. Chances are, someone already had this problem and found a solution.
