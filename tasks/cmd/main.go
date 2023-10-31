package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/4everland/ipfs-top/tasks"
	"github.com/4everland/ipfs-top/tasks/conf"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var (
	cmd      *cobra.Command
	flagconf string
	pass     string
	bc       conf.Data
)

func init() {
	cmd = &cobra.Command{
		Short: "sync s3 file to arweave",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if pass == "" {
				fmt.Println("sign dataItem without password")
			}
			c := config.New(
				config.WithSource(
					file.NewSource(flagconf),
				),
			)
			defer c.Close()

			if err := c.Load(); err != nil {
				panic(err)
			}

			if err := c.Scan(&bc); err != nil {
				panic(err)
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
		},
		TraverseChildren: true,
	}
	cmd.PersistentFlags().StringVarP(&flagconf, "config", "c", "config.yaml", "config file")
	cmd.PersistentFlags().StringVarP(&pass, "pass", "p", "", "mnemonic password")

}

func main() {

	var (
		batchSourceFile string
		singleKey       string
	)
	singleKeyCmd := &cobra.Command{
		Use:   "once",
		Short: "send one s3 object to arseeding",

		Run: func(cmd *cobra.Command, args []string) {
			bridge, err := tasks.NewS3Bridge(bc.Arseeding.GatewayAddr, bc.Arseeding.Mnemonic, bc.S3)
			if err != nil {
				panic(err)
			}
			if singleKey == "" {
				fmt.Println("object key must not be empty")
				return
			}
			tx, err2 := bridge.SendtoAR(context.Background(), pass, singleKey, nil)
			if err2 != nil {
				fmt.Println("send to ar failed: ", err2)
				return
			}
			fmt.Println("success send by hash: ", tx.ItemId)
		},
	}

	batchFileCmd := &cobra.Command{
		Use:   "batch",
		Short: "batch send s3 objects to arseeding",
		Run: func(cmd *cobra.Command, args []string) {
			bridge, err := tasks.NewS3Bridge(bc.Arseeding.GatewayAddr, bc.Arseeding.Mnemonic, bc.S3)
			if err != nil {
				panic(err)
			}
			fi, err := os.Open(batchSourceFile)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
			defer fi.Close()
			br := bufio.NewReader(fi)
			for {
				a, _, err := br.ReadLine()
				if err == io.EOF {
					break
				}
				if len(a) == 0 {
					continue
				}
				tx, err2 := bridge.SendtoAR(context.Background(), pass, string(a), nil)
				if err2 != nil {
					fmt.Printf("send object key %s failed: %s\n", string(a), err2)
					return
				}
				fmt.Println("success send by hash: ", tx.ItemId)
			}
		},
	}
	singleKeyCmd.Flags().StringVarP(&singleKey, "key", "k", "", "s3 object key")
	batchFileCmd.Flags().StringVarP(&batchSourceFile, "batch-file", "s", "", "batch file path")

	cmd.AddCommand(singleKeyCmd, batchFileCmd)
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		panic(err)
	}
}
