package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"github.com/PLAB-IO/aws-cloudformation-debugger/internal/cloudformation"
	"github.com/PLAB-IO/aws-cloudformation-debugger/internal/ui"
	awsCF "github.com/aws/aws-sdk-go/service/cloudformation"
	"log"
	"regexp"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var stackName string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cfdbg",
	Short: "Debugs cloudformation deploys",
	Long: ` Automatically searches the failed AWS event since the 
	latest stable state of your stack.`,

	// App main action
	Run: func(cmd *cobra.Command, args []string) {
		profilestatus, _ := cmd.Flags().GetString("profile")
		regionstatus, _ := cmd.Flags().GetString("region")
		stacknamestatus, _ := cmd.Flags().GetString("stack-name")

		if profilestatus == "" {
			log.Fatal("Please provide --profile option")
		}
	
		if regionstatus == "" {
			log.Fatal("Please provide --region option")
		}
	
		if stacknamestatus == "" {
			log.Fatal("Please provide --stack-name option")
		}
	
		stackName = stacknamestatus
		cloudformation.Region = regionstatus
	
		if err := cloudformation.SetProfile(profilestatus); err != nil {
			panic(err)
		}
		events := lookupOriginalFailed(stackName)

		for _, event := range events {
			rows := [][]string {
				{"Timestamp", event.Timestamp.Format(time.RFC1123)},
				{"Stack Name", *event.StackName},
				{"Stack Status", *event.ResourceStatus},
				{"Logical Resource Id", *event.LogicalResourceId},
				{"FAILED Reason", *event.ResourceStatusReason},
			}
			ui.PaintTable(rows)
		}
	 },
}

func lookupOriginalFailed(stackName string) []awsCF.StackEvent {
	response :=  make([] awsCF.StackEvent, 0)
	events := cloudformation.GetFailEvents(stackName)

	for _, event := range events {
		if !strings.Contains(*event.ResourceStatus, "FAILED") {
			continue
		}
		if "Resource creation cancelled" == *event.ResourceStatusReason {
			continue
		}

		re := regexp.MustCompile(`Embedded stack (arn:.*) was not successfully`)
		matched := re.FindStringSubmatch(*event.ResourceStatusReason)

		if len(matched) == 2 {
			response = append(response, lookupOriginalFailed(matched[1])...)
			continue
		}

		response = append(response, event)
	}
	return response
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persisten flags and configuration settings here,
	// which will be global for the application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.testcobra.yaml)")

	// Local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().String("profile", "default", "AWS Profile")
	rootCmd.Flags().String("region", "eu-west-1", "AWS Region")
	rootCmd.Flags().String("stack-name", "", "Cloudformation Stack Name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".testcobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".testcobra")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
