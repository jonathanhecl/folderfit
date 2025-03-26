# 📂 FolderFit 
Efficient Selection, Perfect Fit

This application efficiently selects folders from a given list that best fit within a user-defined maximum storage capacity. It calculates the size of each source folder and optimizes the selection to maximize the utilization of the available space, minimizing any unused space. This utility is ideal for users who need to identify the optimal set of folders to fit within a specific storage limit.

## 📚 How to Use

```sh
> folderfit <sources> -size=<totalsize> [-verbose]
```

- `<sources>`: A list of folders to be selected from.
- `<totalsize>`: The total storage capacity in bytes.
- `-verbose`: Optional flag to enable verbose output.

## Example

```sh
> folderfit file50kb.dummy file100kb.dummy -size=150000
```

## 🤝 Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

## 🔗 Links

- [GitHub Repository](https://github.com/jonathanhecl/folderfit)
- [Report Issues](https://github.com/jonathanhecl/folderfit/issues)
- [Releases](https://github.com/jonathanhecl/folderfit/releases)