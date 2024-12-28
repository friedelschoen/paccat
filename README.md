# Paccat Package Manager

Paccat is a minimalist and functional package manager designed for simplicity, flexibility, and traceability. Inspired by Nix, it provides a declarative way to manage packages with reproducible builds. Recipes in Paccat are written in a domain-specific language (DSL) that balances expressiveness and simplicity.

## Features

- **Declarative Syntax**: Recipes define dependencies, attributes, and outputs in a clean, functional format.
- **Reproducible Builds**: Build artifacts are hashed, ensuring builds can be reliably reproduced.
- **Modular Design**: Recipes can import other recipes with parameters, enabling reuse and customization.
- **Traceability**: Outputs can trace back to their source expressions for debugging and transparency.

## Recipe Syntax

A Paccat recipe consists of attributes, dependencies, and outputs, described using a functional DSL. Below is an explanation of the syntax.

### Grammar Overview

- **Recipe**: Defines the main structure of the recipe, including optional dependencies (`Require`) and attribute definitions (`Line`).
- **Dependencies**: Declares required attributes and their values.
- **Attributes**: Defines key-value pairs used in the recipe.
- **Outputs**: Specifies how to build the package, with options to control build behavior.
- **Imports**: Allows importing other recipes with parameters.

### Key Constructs

#### 1. **Recipe**
```peg
Recipe <- require:Require? _ lines:Line* _ EOF
```
A recipe consists of optional dependencies (`Require`) and multiple lines of attribute definitions (`Line`).

#### 2. **Dependencies**
```peg
Require <- "[" _ head:RequireContent tail:(_ "," _ val:RequireContent)* _ "]" _ ";"
RequireContent <- key:Key _ "=" _ value:Value / key:Key
```
Dependencies are declared in square brackets and specify attributes that must be provided for the recipe to function.

Example:
```plaintext
[ gcc = "gcc-10", make = "make-4.3" ];
```

#### 3. **Attributes**
```peg
Line <- _ val:Pair? _ ";" _
Pair <- key:Key _ "=" _ value:Value
Key <- [a-zA-Z0-9_]+
```
Attributes define key-value pairs. Keys are alphanumeric strings, and values can be any valid `Value`.

Example:
```plaintext
name = "example-package";
version = "1.0";
```

#### 4. **Outputs**
```peg
Output <- "output" options:OutputOptions* _ script:Value
OutputOptions <- "always" / "interpreter" "=" Value / "try"
```
Outputs define the build script and options for the recipe.

Example:
```plaintext
output always interpreter = "bash" {
    echo "Building example-package";
    make;
    make install;
};
```

#### 5. **Imports**
```peg
Import <- "import" _ path:(Path / String) _ "{" _ params:ImportParams? _ "}"
ImportParams <- Pair ("," Pair)*
```
Imports allow a recipe to include another recipe with parameters.

Example:
```plaintext
import "./base.pcr" { compiler = "gcc", flags = "-O2" };
```

#### 6. **Values**
Values are versatile and include lists, strings, references, and interpolations.

```peg
Value <- List / String / Multiline / Import / Output / Dependencies / Surrounded / Reference
List <- "[" Value ("," Value)* "]"
String <- '"' StringContent* '"'
StringInterpolation <- "${" Value "}"
```

Example:
```plaintext
description = "This is an example ${name} package";
paths = ["/usr/bin/example", "/usr/lib/example"];
```



## Example Recipe

```plaintext
[ gcc = "gcc-10", make = "make-4.3" ];

name = "example-package";
version = "1.0";
description = "A sample package for demonstration purposes";

output always interpreter = "bash" {
    echo "Building ${name} version ${version}";
    ./configure --prefix=/usr;
    make;
    make install;
};
```

## Commands

Paccat supports the following commands:

1. **Hash**: Calculate the hash of a recipe or attribute.
   ```sh
   paccat hash <attribute>
   ```

2. **Install**: Build and install the recipe's attribute.
   ```sh
   paccat install [attribute]
   ```

3. **Evaluate**: Evaluate an attribute and print the result. Use `--ast` for source tracing.
   ```sh
   paccat evaluate <attribute> --ast
   ```

4. **Remove**: Remove a package by its hash.
   ```sh
   paccat remove <hash>
   ```

## Summary

Paccat is a simple yet powerful package manager tailored for developers who value minimalism and reproducibility. Its DSL ensures that recipes remain clean and expressive, while its modular and traceable design keeps package management efficient and transparent.