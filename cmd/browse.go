/*
Copyright © 2023 tigerinus

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/spf13/cobra"
)

const lastUpdatedKey = "lastUpdated"

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use: "browse",
	RunE: func(cmd *cobra.Command, args []string) error {
		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			return err
		}

		lastUpdated := time.Now()

		ctx := context.WithValue(cmd.Context(), lastUpdatedKey, &lastUpdated)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		go func() {
			for {
				select {
				case <-ctx.Done():
					break
				case <-time.After(1 * time.Second):
					lastUpdated := ctx.Value(lastUpdatedKey).(*time.Time)
					if time.Since(*lastUpdated) > 10*time.Second {
						cancel()
					}
				}
			}
		}()

		services := make(chan *zeroconf.ServiceEntry)

		go discoverServices(ctx, services, dicoverServiceEntries)

		if err := resolver.Browse(ctx, "_services._dns-sd._udp", "", services); err != nil {
			return err
		}

		<-ctx.Done()

		return nil
	},
}

func discoverServices(ctx context.Context, services chan *zeroconf.ServiceEntry, handle func(ctx context.Context, service, domain string)) {
	for {
		select {
		case <-ctx.Done():
			break
		case service := <-services:
			go handle(ctx, service.Instance, service.Domain)
		}
	}
}

func dicoverServiceEntries(ctx context.Context, service string, domain string) {
	service = strings.TrimSuffix(service, ".local")

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		fmt.Printf("error when creating a new resolver: %s\n", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case entry, ok := <-entries:
				if !ok {
					break
				}

				fmt.Printf("%s:\t%s\n", service, entry.Instance)

				lastUpdated := ctx.Value(lastUpdatedKey).(*time.Time)
				*lastUpdated = time.Now()
			}
		}
	}()

	if err := resolver.Browse(ctx, service, domain, entries); err != nil {
		fmt.Printf("error when browsing for service entries: %s\n", err.Error())
	}
}

func init() {
	rootCmd.AddCommand(browseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// browseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// browseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
