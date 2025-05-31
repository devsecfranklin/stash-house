# C Projects

## Size Optimization

- Try to reduce the number of string literals in your code.
- Optimize variables.
  - Remove unused variables and be economical with the data types you use.
  - Don’t use long (4 bytes) if only int (2 bytes) are needed.
  - Don’t use float if only whole numbers will be assigned.

- Declare variables as locally as possible to save memory
  - Anything declared at global scope will exist for the entire lifetime of the sketch.
  - Variables that are only used inside one function, can be declared within that function.
  - If a variable needs to be shared between multiple functions, consider passing it as a parameter.

Avoid recursive code to reduce memory usage. Use iterative functions instead.

Minimize array sizes. Don’t declare larger arrays than needed.

Check your libraries to save storage. Avoid including large libraries unless you need to. Save storage by implementing your own code that only does what you need.

Store constant data in storage instead of memory with PROGMEM.
