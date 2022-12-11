# echow
echow is a small tool that generates echo commands to write binary file.

[简体中文](README.zh-CN.md) | [English](README.md)



# Usage

The usage of echow is very simple.

### Basic usage

```
NAME:
   echow - Generate echo commands to write binary file

USAGE:
   echow [global options] [arguments...]

GLOBAL OPTIONS:
   --file file, -f file      file path
   --number value, -n value  split the file into the specified number of parts (default: 1)
   --format hex/oct          choose octal or hex format: hex/oct (default: "hex")
   --name value              specify the file name (default: -f parameter value)
   --help, -h                show help (default: false)
```

### Run the Example

> easy example

```
PS C:\> .\echow.exe -f .\Test.bin -n 3
echo -n -e "\x68\x65\x6c" >> Test.bin
echo -n -e "\x6c\x6f\x20" >> Test.bin
echo -n -e "\x77\x6f\x72\x6c\x64" >> Test.bin
```

> Output the result to a file

Use `>` redirection symbol in CMD:

```
.\echow.exe -f .\Test.bin -n 2 > test.txt
```

If in PowerShell, the following command is recommended:

```powershell
.\echow.exe -f .\Test.bin -n 2 | Out-File -Encoding ASCII test.txt
```



# Tips

> A problem occurred: extra data `0xFFFE` in file header

![](PowerShell_Redirector_Output.png)

This is because when you write to a file using the `>` redirection in PowerShell, it treats the file as Unicode. Therefore, `0xFFFE` is automatically added to the header of the file, which is a [Unicode byte order mark](https://learn.microsoft.com/en-us/windows/win32/intl/using-byte-order-marks?redirectedfrom=MSDN) used to indicate the Unicode encoding scheme used in the file.

> Solution

Using the [Out-File](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/out-file) redirection command in PowerShell instead of the `>` symbol, and using `-Encoding ASCII` to specify the encoding as ASCII will not cause this problem.
