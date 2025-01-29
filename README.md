# InfoGath

InfoGath is a tool designed to gather information from URLs quickly and efficiently. It supports various flags to customize its behavior according to your needs.

## Installation

To install InfoGath, clone the repository and compile the source code using the following command:

```sh
go build -o InfoGath *.go
```

## Usage

```golang
./InfoGath --file <file> [-f <file>] --threads <int> [-t <int>] --output <file> [-o <file>] --crawl <file> [-c <file>]
```

## Flags
- file, -f <file>: Specify the file containing URLs to fetch.
- threads, -t <int>: Indicate the number of threads you want to use.
- output, -o <file>: Indicate the name of the output file.
- crawl, -c <file>: Indicate whether to detect inputs as forms or input labels.

## Example

### Default Example

Uses 1 thread and default output file is: active_subdomains.
```golang
./InfoGath -f subdomains.txt -d -o
```

### Custom Example

```golang
./InfoGath -file subdomains.txt --threads 4 --output results.txt --crawl results.txt -dw subdomains.txt
```

## Notes

The tool compiles by running go build *.go in the src directory.
Replace subdomains.txt, results.txt, and crawl.txt with the actual file paths containing URLs, output, and crawl data, respectively.
