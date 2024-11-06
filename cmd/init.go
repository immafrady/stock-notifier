package cmd

import (
	"errors"
	"github.com/immafrady/stock-notifier/templates"
	"github.com/immafrady/stock-notifier/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	path2 "path"
	"path/filepath"
)

var t string
var p string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化配置文件",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		var err error
		if p == "" {
			path, err = os.Getwd()
			utils.PanicOnError(err)
		} else {
			path, err = filepath.Abs(p)
			utils.PanicOnError(err)

			info, err := os.Stat(path)
			utils.PanicOnError(err)

			if !info.IsDir() {
				panic(errors.New("请输入文件夹的路径"))
			}
		}

		var tmpl string
		switch t {
		case "jsonc":
			tmpl = templates.JsoncTmpl
			path = path2.Join(path, "config.jsonc")
		case "yaml":
			tmpl = templates.YamlTmpl
			path = path2.Join(path, "config.yaml")
		default:
			log.Fatalln("[error]仅支持jsonc和yaml格式")
		}

		err = os.WriteFile(path, []byte(tmpl), 0644)
		utils.PanicOnError(err)
		log.Println("成功写入文件，请编辑配置文件后运行程序")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&t, "type", "t", "jsonc", "导出类型，支持jsonc和yaml")
	initCmd.Flags().StringVarP(&p, "path", "p", "", "输出路径，默认工作目录")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
