package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var keepAlive bool

var webhookCmdGroup = &cobra.Command{
	Use:   "webhook",
	Short: "commands here are used for interacting with the strava webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		// Nothing to see here
	},
}

var createSubscription = &cobra.Command{
	Use:   "create",
	Short: "create a subscription for your app. this is a one-time run. note that this will spawn a server at the callback url that recieves events",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		stravaApp, _, err := createApp()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		id, server, wg, err := stravaApp.CreateSubscription()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("server is running on %s\n", server.Addr)
		fmt.Printf("subscription id: %d\n", id)
		if keepAlive {
			wg.Wait()
		} else {
			wg.Done()
		}
	},
}

var launchWebhookServer = &cobra.Command{
	Use:   "launch-server",
	Short: "launch the server. only do this if you have already created a webhook subscription.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		stravaApp, _, err := createApp()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		server, wg, err := stravaApp.LaunchWebhookServer()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("server is running on %s\n", server.Addr)
		if keepAlive {
			wg.Wait()
		} else {
			wg.Done()
		}
	},
}

var viewSubscription = &cobra.Command{
	Use:   "view",
	Short: "view subscription for your app",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		stravaApp, _, err := createApp()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = stravaApp.ViewSubscription()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

var deleteSubscription = &cobra.Command{
	Use:   "delete [subscription id]",
	Short: "delete subscription for your app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subscriptionID := args[0]
		stravaApp, _, err := createApp()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = stravaApp.DeleteSubscription(subscriptionID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func init() {
	createSubscription.Flags().BoolVar(&keepAlive, "keep-alive", true, "keep the server alive after creation. the server is needed to get events from the webhook")
	launchWebhookServer.Flags().BoolVar(&keepAlive, "keep-alive", true, "keep the server alive after creation. the server is needed to get events from the webhook")
	webhookCmdGroup.AddCommand(createSubscription)
	webhookCmdGroup.AddCommand(launchWebhookServer)
	webhookCmdGroup.AddCommand(viewSubscription)
	webhookCmdGroup.AddCommand(deleteSubscription)
	RootCmd.AddCommand(webhookCmdGroup)
}
