## symbol-to-offset

### Usage

```
Usage: symbol-to-offset EXECUTABLE SYMBOL

        EXECUTABLE      Path to the executable ELF file.
        SYMBOL          The name of the symbol in the executable file.
```

For instance:

```shell
$ symbol-to-offset ./utrace runtime.asmcgocall.abi0
+-------------------------+-----------+---------------+---------+------------+----------------+
|         SYMBOL          | SYMBOL VA | SYMBOL OFFSET | SECTION | SECTION VA | SECTION OFFSET |
+-------------------------+-----------+---------------+---------+------------+----------------+
| runtime.asmcgocall.abi0 |    475380 |         75380 | .text   | 4019A0     | 19A0           |
+-------------------------+-----------+---------------+---------+------------+----------------+
```

### Installation

```shell
go install github.com/maxgio92/symbol-to-offset@latest
```
