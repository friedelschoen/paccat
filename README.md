# upkg: A Minimalistic Recipe-based Build System

`upkg` is a lightweight, flexible build system inspired by the Nix language. It provides a simple Domain-Specific Language (DSL) for defining and executing build recipes. Recipes are defined in a human-readable format and allow interpolation, function calls, and dynamic attribute resolution.

## Features

- **Recipe-based builds**: Write modular, reusable build recipes.
- **Dynamic attribute resolution**: Use variables, interpolations, and function calls to dynamically configure builds.
- **Outputs**: Recipes produce build outputs that are automatically managed.
- **Simple dependency management**: Specify required attributes for recipes.
- **Human-readable DSL**: A syntax designed for simplicity and readability.


## Getting Started

### Installation

1. Clone this repository:
   ```sh
   git clone https://github.com/yourusername/upkg.git
   cd upkg
   ```

2. Build the `upkg` binary:
   ```sh
   go build -o upkg upkg.go
   ```

3. Install the binary to your system:
   ```sh
   sudo install -m 755 upkg /usr/local/bin
   ```


## Writing a Recipe

A `upkg` recipe is a plain text file with the following structure:

### Basic Recipe Example

```plaintext
name = "my-app"

# Build command with dynamic attributes
build = ""
    echo 'Building ${name}'
    mkdir -p ${out}
    cp source-code/* ${out}
""
```

### Features of the Recipe DSL

#### Key-Value Attributes

- Define static or interpolated attributes:
  ```plaintext
  key = "value"
  key_with_interpolation = "value-${variable}"
  ```

#### Multiline Strings

- Use `""` for defining multiline build commands:
  ```plaintext
  build = ""
      echo 'Compiling...'
      gcc -o ${out}/app main.c
  ""
  ```

#### Function Calls

- Call reusable recipes using the function syntax:
  ```plaintext
  dependencies = [fetch.url("https://example.com/source.tar.gz")]
  ```


## Using upkg

1. **Define your recipe**: Create a file like `example.recipe` with your build logic.

2. **Run upkg**: Execute the `upkg` binary with your recipe:
   ```sh
   upkg example.recipe
   ```

3. **Result**: The build result will be symlinked to `./result` for easy access.


## Code Overview

### Core Components

1. **Recipes**: Defined in the `recipe` package. 
    - `Recipe` struct holds attributes and build logic.
    - Recipes can include required attributes and dynamic attribute evaluation.

2. **Build Context**: 
    - The `Context` struct manages the state during builds, including variable resolution and build outputs.

3. **DSL Parsing**: 
    - The grammar is defined in a PEG file (`recipe.peg`) and parsed using the `pigeon` parser generator.

4. **Build Execution**: 
    - Attributes and build commands are resolved and executed in a sandboxed environment using `os/exec`.


## Contributing

1. Fork this repository.
2. Create a branch for your feature/bugfix.
3. Submit a pull request with detailed explanations.


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


## Acknowledgments

Inspired by the Nix language and its derivatives, this project aims to simplify the build process with a clear and concise syntax.