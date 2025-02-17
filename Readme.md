# ScriptGen - AI-Powered CLI Script Generator

ScriptGen is a command-line tool written in Go that leverages OpenAI's language models to generate scripts based on user-defined requirements. The tool supports multiple scripting languages, including Python, NodeJS, PHP, and Bash, and ensures compatibility with both MacOS and Linux environments.

## Features
- Generates scripts dynamically based on user-defined requirements
- Supports structured JSON output for predictable results
- Supports Python, NodeJS, PHP, and Bash
- Automatically writes the script to a file and makes it executable

## Installation
### Prerequisites
- Go 1.21 or later installed
- An OpenAI API key

### Clone the Repository
```sh
git clone https://github.com/ohnotnow/scriptwriter_go
cd scriptwriter_go
```

### Build the Executable
```sh
# I like to use a shorter name for the executable
go build -o ws .
```

## Usage
Run the tool and enter the script requirements when prompted:

```sh
./ws
```

You will be prompted to enter a description of the script you need. The tool will then generate the script based on your input and save it to a file. The filename and script content will be displayed in the terminal.

### Example
```sh
./ws
Enter the requirements for the script: A Bash script that prints system information.
```
#### Output:
```sh
#!/bin/bash
echo "System Information:"
uname -a
```
### Script written to output.sh
```

## License
This project is licensed under the MIT License.
