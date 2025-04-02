# 📂 FolderFit 
Efficient Selection, Perfect Fit for Folders and Files

This application efficiently selects folders from a given list that best fit within a user-defined maximum storage capacity. It calculates the size of each source folder and optimizes the selection to maximize the utilization of the available space, minimizing any unused space. This utility is ideal for users who need to identify the optimal set of folders to fit within a specific storage limit.

```
Author note:
I created this project because I hate to leave free space on CDs/DVDs.
I had originally created it in VB.NET 20 years ago, now it is multiplatform thanks to Golang.
```
[Original project in Forum](https://www.gs-zone.org/temas/fileset-optimal-design-seleccion-inteligente-de-archivos.89846/)

## 📚 How to Use

```sh
> folderfit <sources> -size=<totalsize> [-verbose]
```

- `<sources>`: A list of folders to be selected from.
- `<totalsize>`: The total storage capacity in bytes.
- `-verbose`: Optional flag to enable verbose output.

## Example

```sh
> folderfit folder50kb file100kb -size=150000

FolderFit v 1.0.1
file100kb - 97 KB
folder50kb - 48 KB

Selection size: 146 KB / 146 KB
Free space: 0 bytes
```

## 🤝 Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

## 🔗 Links

- [GitHub Repository](https://github.com/jonathanhecl/folderfit)
- [Report Issues](https://github.com/jonathanhecl/folderfit/issues)
- [Releases](https://github.com/jonathanhecl/folderfit/releases)
