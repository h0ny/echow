package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	format   string
	filePath string
	number   int64
	fileSize int64
	name     string
)

func main() {
	app := &cli.App{
		Name:            "echow",
		Usage:           "Generate echo commands to write binary file",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "`file` path",
				Destination: &filePath,
				Required:    true,
				Action: func(ctx *cli.Context, v string) error {
					fileInfo, err := os.Stat(filePath)
					if err != nil {
						return fmt.Errorf("%s", err)
					}
					// 判断是否是一个文件
					if !fileInfo.Mode().IsRegular() {
						return fmt.Errorf("file not is a regular file")
					}
					fileSize = fileInfo.Size()
					return nil
				},
			},
			&cli.Int64Flag{
				Name:        "number",
				Aliases:     []string{"n"},
				Value:       1,
				Usage:       "split the file into the specified number of parts",
				Destination: &number,
				Action: func(ctx *cli.Context, v int64) error {
					if v > fileSize || v < 1 {
						return fmt.Errorf("Flag number value %v out of range [1-%v]", v, fileSize)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "format",
				Value:       "hex",
				Usage:       "choose octal or hex format: `hex/oct`",
				Destination: &format,
				Action: func(ctx *cli.Context, s string) error {
					if !strings.EqualFold(s, "hex") && !strings.EqualFold(s, "oct") {
						return fmt.Errorf("invalid value: %s (only hex or oct)\n", s)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "name",
				DefaultText: "-f parameter value",
				Usage:       "specify the file name",
				Destination: &name,
			},
		},

		Action: func(cCtx *cli.Context) error {
			// 判断name参数是否有值传入。如果用户没有传入，则设置默认值。
			if !cCtx.IsSet("name") {
				name = filepath.Base(filePath)
			}

			splitFile, err := SplitFile(filePath, int(number))
			if err != nil {
				return cli.Exit(err, 86)
			}
			// 输出echo写入命令
			if strings.EqualFold(format, "hex") {
				for _, bytes := range splitFile {
					fmt.Print("echo -n -e \"")
					for _, b := range bytes {
						fmt.Printf("\\x%x", b)
					}
					fmt.Printf("\" >> %s\n", name)
				}
			} else if strings.EqualFold(format, "oct") {
				for _, bytes := range splitFile {
					fmt.Print("echo -n -e \"")
					for _, b := range bytes {
						fmt.Printf("\\0%o", b)
					}
					fmt.Printf("\" >> %s\n", name)
				}
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// SplitFile
//
//	@Description: 将文件文拆分成指定数量的 []byte
//	@param filename 要拆分的文件名
//	@param numSlices 拆分后的数量
func SplitFile(filename string, numSlices int) ([][]byte, error) {
	// 读取文件内容，返回 []byte
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	sliceSize := len(data) / numSlices
	// 创建一个二维字节切片，用来存储拆分后的数据块
	chunks := make([][]byte, numSlices)

	for i := 0; i < numSlices; i++ {
		start := i * sliceSize
		end := start + sliceSize
		if i == numSlices-1 {
			// 最后一份可能会比其他份大一些，需要特殊处理
			end = len(data)
		}
		chunks[i] = data[start:end]
	}

	return chunks, nil
}
