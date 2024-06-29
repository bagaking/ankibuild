# AnkiBuild

AnkiBuild is a Go application designed to automate the creation of Anki flashcard packages (.apkg files) from configuration files.

## Overview

This project facilitates the creation of Anki flashcard packages by providing a simple way to define questions and answers (QnAs) in `.apkg.toml` configuration files. The application will search for these configuration files within the current directory and its subdirectories to generate corresponding `.apkg` files.

## Quick Start

1. Clone the repository.
2. Navigate to the root directory of the project.
3. Make sure you have Go installed on your machine.
4. Run `go build` to build the binary or `go run main.go` to directly run the application.

### Example

1. Clone the repository.
2. execute `go run .` or `make example`
3. the .apkg file of `./example/example.apkg.toml` are generated in `./example`

## Usage

To generate `.apkg` files:

1. Define your flashcards in `.apkg.toml` files within the project directory.
2. Make sure each `.apkg.toml` file follows the `KnowledgePage` struct defined in `conf.go`, which includes the title of the package and a slice of QnACards.
3. Run the binary or use `go run main.go` to start the application. This will generate `.apkg` files in the origin directory.

configuration detail @see [FORMAT](./FORMAT.md)

## Project Structure

- `builder.go`: Contains the logic to parse `.apkg.toml` files and build `.apkg` files.
- `conf.go`: Defines the structures for the configuration files (`KnowledgePage` and `QnACard`).
- `main.go`: The entry point of the application which invokes the build process.
- `apkg/`: Root package of the apkg formatter writer.

## Contributing

Contributions to AnkiBuild are welcome. Please ensure that your contributions adhere to the following guidelines:

- Write clean, readable, and maintainable code.
- Provide documentation for any new features or changes.
- Submit Pull Requests with detailed descriptions of changes.
- Make sure that the changes do not break existing functionality.

## License

This project is licensed under the [MIT License](LICENSE).

## Authors

- [Bagaking](https://github.com/bagaking) - Initial work

Feel free to contact the maintainers of this project if you have any questions or feedback.

## Acknowledgments

- Anki - for providing an awesome flashcard application.
- BurntSushi - for the `toml` package used in parsing configuration files.
- Contributors to the `anki_barn` and `goulp` packages for their utility functions in the project.