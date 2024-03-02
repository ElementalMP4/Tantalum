# Tantalum

A Go based mirroring tool, Tantalum copies files from one directory to another simply and efficiently.

## Configuration

Here is an example config file:

```json
{
    "couples": [
        {
            "left": "C:\\Users\\ElementalMP4\\Documents",
            "right": "\\\\elementalmp4-backup-server\\Documents",
            "forceUpdate": false
        }
    ],
    "output": true
}
```

This config will copy all files and directories from `left` and duplicate them on `right`, only copying the latest copy of each file.

To copy every file regardless of modified time, set `forceUpdate` to true

It is possible to have as many couples as you wish, and you can copy from and to network shares, so long as you have already authenticated.

## Running

1) Download this repository
2) Compile with `go build`
3) Create a `config.json` file which MUST be in the same directory as `Tantalum`

Note: Tantalum works on both Windows and Linux.