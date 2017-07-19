# encoding

Basic tool for fixing mixed NFC/NFD encoding. Will walk recursively throughout a folder and detect NFD encoded filenames. If not in dry-run mode, it will rename them to using NFC normalization form.  
Not working on MacOSX, as MacOSX always store filenames in NFD.

## Usage

```
$ go run main.go --help

Options:

  -h, --help      display help information
  -f, --folder   *Root folder to start scanning for NFD encoded filenames
  -d, --dry-run   Use dry run to see all changes without applying them

```
