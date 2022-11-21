package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "file-share",
	Short: "A file sharing tool.",
	Long:  `A file sharing tool.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		pathname, err := flags.GetString("path")
		if err != nil {
			return err
		}
		_, err = os.Stat(pathname)
		if err == os.ErrNotExist {
			return err
		}

		host, err := flags.GetString("host")
		if err != nil {
			return nil
		}

		serveMux := http.NewServeMux()
		serveMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			uri := request.RequestURI
			absPath, err := filepath.Abs(pathname + uri)
			if err != nil {
				return
			}
			_, err = os.Stat(absPath)
			if err == os.ErrNotExist {
				http.NotFound(writer, request)
				return
			}

			file, err := os.Open(absPath)
			if err != nil {
				return
			}
			defer file.Close()

			_, fileName := filepath.Split(uri)
			if fileName == "" {
				http.NotFound(writer, request)
				return
			}

			writer.Header().Set("Content-Type", "application/octet-stream")
			writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
			writer.Header().Set("Content-Transfer-Encoding", "binary")

			_, err = io.Copy(writer, file)
			if err != nil {
				return
			}

			fmt.Println(absPath + uri)
		})

		fmt.Printf("监听 - %s\n", host)

		return http.ListenAndServe(host, serveMux)
		//router := gin.Default()
		//router.GET("file/*", func(c *gin.Context) {
		//	uri := c.Request.RequestURI
		//	_, file := filepath.Split(uri)
		//	c.Header("Content-Type", "application/octet-stream")
		//	c.Header("Content-Disposition", "attachment; filename="+file)
		//	c.Header("Content-Transfer-Encoding", "binary")
		//	absPath, err := filepath.Abs(pathname + uri)
		//	if err != nil {
		//		return
		//	}
		//	c.File(absPath)
		//})
		//router.Run(host)
		//
		//return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.file-share.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("path", "p", "./", "file path to be shared")
	rootCmd.Flags().StringP("host", "H", ":8888", "access path")
}
