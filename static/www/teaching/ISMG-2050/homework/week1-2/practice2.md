# Practice 2: Functions, AutoFill, Formatting, and Sheet Management

## Overview

In this project, you will open the workbook created in Practice 1. You will insert statistical functions using the AutoSum tool, propagate formulas using the Fill Handle, format numerical and text data, adjust row and column geometry, and duplicate and prepare a worksheet for subsequent reporting periods.

* **File Needed:** `[your initials]Practice01.xlsx` *(created in Practice 1)*
* **Completed Project File Name:** `[your initials]Practice02.xlsx`

---

## Instructions

### Step 1: Open and Save As

1. Open `[your initials]Practice01.xlsx`. *(Click **Enable Editing** if opened in Protected View).*
2. Resave the workbook as `[your initials]Practice02.xlsx`.
3. Click **Save**.

---

### Step 2: Calculate Totals Using SUM

1. Select cell **H5** and click the **AutoSum** button *(Home tab > Editing group)*.
2. Press **Enter** to accept the suggested range (`D5:G5`). The result is `31`.
3. Type `Total` in cell **H4** and press **Enter**. *(Center alignment is automatically applied based on adjacent column formatting).*
4. Select cell **D19**.
5. Click the **AutoSum** button *(Home tab > Editing group)* to insert the `SUM` function with arguments `D5:D18`. The result is `63`.

---

### Step 3: Calculate Average Rental Days by Month

1. Select cell **D20**.
2. Click the **AutoSum** button dropdown arrow *(Home tab > Editing group)* and select **Average**.
3. In the formula bar or inline editor, adjust the argument range from the suggested range to exclude the total in row 19:
* Change `D5:D19` to `D5:D18`.


4. Press **Enter**. The result is `4.5`.

---

### Step 4: Copy Functions Using the Fill Handle

1. Select cell **H5**.
2. Double-click the **Fill Handle** (the square pointer in the lower-right corner of the active cell). The `SUM` formula will fill down through cell **H19**.
3. Select the range **D19:D20**.
4. Click and drag the Fill Handle across through column H (**D19:H20**).
5. Verify the relative cell references:
* Select cell **E19** and inspect the Formula Bar to confirm the `SUM` formula adjusted for column E.
* Select cell **E20** and inspect the Formula Bar to confirm the `AVERAGE` formula adjusted for column E.



---

### Step 5: Apply Number and Label Formatting

1. Select the range **D20:H20**.
2. Click the **Number Format** dropdown list *(Home tab > Number group)* and choose **Number** (formatted to two decimal places).
3. Select the header range **A4:H4** and click **Bold** *(Home tab > Font group)*.
4. Select cells **A1:A2**, click the **Font Size** dropdown *(Home tab > Font group)*, and select **16 pt**.
5. Select cell **A2** and edit the text to display `First Period` instead of `First Quarter`.

---

### Step 6: Adjust Column Widths

1. Click and drag across the column headings for **B** and **C**.
2. Double-click the boundary between column headers **C** and **D** to **AutoFit** the selected columns.
3. Click and drag across the column headings from **D** through **H**.
4. Click **Format** *(Home tab > Cells group)* > select **Column Width...**
5. Type `7` in the Column Width dialog box and click **OK**.

---

### Step 7: Adjust Row Heights

1. Select the heading for **Row 4**.
2. Hold `Ctrl` and select the headings for rows **19** and **20** (three rows selected simultaneously).
3. Right-click the **Row 4** heading and select **Row Height...**
4. Set the row height to `24` and click **OK**.

---

### Step 8: Center Labels Across Selection

1. Select the cell range **A1:H2**.
2. Click the **Alignment Settings** launcher arrow *(Home tab > Alignment group)*.
3. In the *Horizontal* alignment dropdown, select **Center Across Selection**.
4. Click **OK**.

---

### Step 9: Insert a Record and Enter Data

1. Right-click the row heading for **Row 13** and select **Insert** *(existing rows shift down)*.
2. Enter the following property record across row 13:
* **A13:** `Our Weekend Cottage` *(press Tab)*
* **B13:** `Walker` *(press Tab)*
* **C13:** Right-click > select **Pick From Drop-down List** > select **MN** *(press Tab)*
* **D13:** Type the first value and press Tab. Enter the remaining monthly values for columns E, F, and G.


3. Verify that cell **H13** automatically evaluates the total (`14`) via the copied relative formulas.

---

### Step 10: Apply Borders, Fill Colors, and Alignment

1. Select cells **A4:H21**, click the **Borders** dropdown *(Home tab > Font group)*, and choose **All Borders**.
2. Select the header range **A4:H4**, click the **Fill Color** dropdown, and choose **Blue, Accent 1, Lighter 80%** *(Column 5, Row 2)*.
3. With cell **A4** selected, click **Format Painter** *(Home tab > Clipboard group)*, then drag across **A20:H21** to apply matching bold and center properties.
4. Select cells **A20:A21** and click **Align Right** *(Home tab > Alignment group)*.
5. Select cells **D20:H21** and click **Align Right**.
6. Select cells **D21:H21** and use the **Decrease Decimal** / **Increase Decimal** buttons to enforce consistent two-decimal-place formatting across the averages.

---

### Step 11: Configure Sheet Name and Tab Color

1. Double-click the **Sheet1** tab name, type `Rental Days`, and press **Enter**.
2. Right-click the `Rental Days` tab, select **Tab Color**, and choose **Blue, Accent 1** *(Column 5, Row 1)*.
3. Select cell **A1**.

---

### Step 12: Duplicate and Clear Data for Next Period

1. Right-click the **Rental Days** sheet tab and select **Move or Copy...**
2. In the dialog box:
* Select `(move to end)` in the *Before sheet* list.
* Check the **Create a copy** checkbox.
* Click **OK**.


3. Rename the duplicated tab (`Rental Days (2)`) to `Next Period`.
4. Right-click the **Next Period** tab, select **Tab Color**, and choose **Gold, Accent 4** *(Column 8, Row 1)*.
5. Select the raw data cells **D5:G19** and press **Delete**.
*(Formulas dependent on these cells will display `0` or `#DIV/0!` errors due to empty division parameters).*
6. Update styling accents on the new sheet:
* Select **A4:H4** > set Fill Color to **Gold, Accent 4, Lighter 80%** *(Column 8, Row 2)*.
* Select **A20:H21** > set Fill Color to **Gold, Accent 4, Lighter 80%**.


7. Press `Ctrl + Home` to reset the view to cell **A1**.

---

### Step 13: Finalize and Save

1. Switch back to the **Rental Days** sheet tab and press `Ctrl + Home`.
2. Save your workbook (`[your initials]Practice02.xlsx`).
3. Submit the completed `.xlsx` file via the designated Canvas assignment portal.
