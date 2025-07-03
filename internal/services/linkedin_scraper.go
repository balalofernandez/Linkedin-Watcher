package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/spf13/viper"
)

type LinkedInConnection struct {
	Name             string
	ProfileURL       string
	ConnectionDegree string
	Location         string
}

type Cookie struct {
	Domain   string  `json:"domain"`
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expirationDate"`
	HTTPOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
}

func randomDelay(minMs, maxMs int) {
	d := time.Duration(rand.Intn(maxMs-minMs)+minMs) * time.Millisecond
	time.Sleep(d)
}

func loginLinkedIn(ctx context.Context) error {
	username := viper.GetString("LINKEDIN_USERNAME")
	password := viper.GetString("LINKEDIN_PASSWORD")
	if username == "" || password == "" {
		return fmt.Errorf("LinkedIn credentials not set in environment variables")
	}
	fmt.Println("[linkedin_scraper] Starting LinkedIn login...")
	randomDelay(500, 1500)
	res := chromedp.Run(ctx,
		chromedp.Navigate("https://www.linkedin.com/login"),
	)
	randomDelay(500, 1500)
	res2 := chromedp.Run(ctx,
		chromedp.WaitVisible(`#username`, chromedp.ByID),
	)
	randomDelay(500, 1500)
	res3 := chromedp.Run(ctx,
		chromedp.SendKeys(`#username`, username, chromedp.ByID),
	)
	randomDelay(500, 1500)
	res4 := chromedp.Run(ctx,
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
	)
	randomDelay(500, 1500)
	res5 := chromedp.Run(ctx,
		chromedp.Click(`button[type='submit']`, chromedp.ByQuery),
	)
	fmt.Println("[linkedin_scraper] Login actions finished, waiting for 2FA or navigation...")

	waitCh := make(chan error, 1)
	ctx2, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	go func() {
		waitErr := chromedp.Run(ctx2,
			chromedp.WaitVisible(`ul`, chromedp.ByQuery),
		)
		waitCh <- waitErr
	}()

	select {
	case err := <-waitCh:
		if err == nil {
			fmt.Println("[linkedin_scraper] Navigation succeeded before 2FA timeout.")
			return nil
		}
		fmt.Println("[linkedin_scraper] Navigation error while waiting for 2FA:", err)
		return err
	case <-ctx2.Done():
		fmt.Println("[linkedin_scraper] 2FA timeout reached, continuing...")
	}

	res6 := chromedp.Run(ctx,
		chromedp.WaitNotPresent(`#username`, chromedp.ByID),
	)
	fmt.Println("[linkedin_scraper] Login finished, errors:", res, res2, res3, res4, res5, res6)
	return res6
}

func loadCookiesFromFile(path string) ([]*network.CookieParam, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var raw []map[string]interface{}
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, err
	}
	var cookies []*network.CookieParam
	for _, c := range raw {
		cookie := &network.CookieParam{
			Name:   c["name"].(string),
			Value:  c["value"].(string),
			Domain: c["domain"].(string),
			Path:   c["path"].(string),
			Secure: c["secure"].(bool),
		}
		if exp, ok := c["expirationDate"]; ok {
			if f, ok := exp.(float64); ok {
				fmt.Println("expirationDate", f)
				t := time.Unix(int64(f), 0)
				epoch := cdp.TimeSinceEpoch(t)
				cookie.Expires = &epoch
			}
		}
		if httpOnly, ok := c["httpOnly"]; ok {
			cookie.HTTPOnly = httpOnly.(bool)
		}
		cookies = append(cookies, cookie)
	}
	return cookies, nil
}

// ScrapeLinkedInConnections scrapes all connections for a given LinkedIn ID.
func ScrapeLinkedInConnections(ctx context.Context, linkedinID string) ([]LinkedInConnection, error) {
	fmt.Println("[linkedin_scraper] ScrapeLinkedInConnections called with id:", linkedinID)
	url := fmt.Sprintf("https://www.linkedin.com/search/results/people/?connectionOf=%%5B\"%s\"%%5D", linkedinID)
	var connections []LinkedInConnection

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Try to load cookies
	cookiePath := "internal/services/cookie.json"
	cookies, err := loadCookiesFromFile(cookiePath)
	if err == nil && len(cookies) > 0 {
		fmt.Println("[linkedin_scraper] Using cookies from cookie.json for authentication...")
		err = chromedp.Run(ctx,
			chromedp.Navigate("https://www.linkedin.com"),
			chromedp.ActionFunc(func(ctx context.Context) error {
				for _, c := range cookies {
					err := network.SetCookie(c.Name, c.Value).
						WithDomain(c.Domain).
						WithPath(c.Path).
						WithExpires(c.Expires).
						WithHTTPOnly(c.HTTPOnly).
						WithSecure(c.Secure).
						Do(ctx)
					if err != nil {
						fmt.Println("[linkedin_scraper] Failed to set cookie:", c.Name, err)
					}
				}
				return nil
			}),
		)
		if err == nil {
			fmt.Println("[linkedin_scraper] Cookies set, navigating to connections page:", url)
			if err := chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.WaitVisible(`ul[role='list']`, chromedp.ByQuery),
			); err == nil {
				fmt.Println("[linkedin_scraper] Navigation with cookies successful, scraping nodes...")
				var nodes []*cdp.Node
				if err := chromedp.Run(ctx,
					chromedp.Nodes(`ul[role='list'] > li`, &nodes, chromedp.ByQueryAll),
				); err == nil {
					fmt.Printf("[linkedin_scraper] Found %d connection nodes\n", len(nodes))
					for i, node := range nodes {
						var name, profileURL, degree, location string
						var anchors []*cdp.Node
						chromedp.Run(ctx, chromedp.Nodes("a", &anchors, chromedp.ByQueryAll, chromedp.FromNode(node)))
						for _, a := range anchors {
							var href, text string
							chromedp.Run(ctx, chromedp.AttributeValue("a", "href", &href, nil, chromedp.ByQuery, chromedp.FromNode(a)))
							chromedp.Run(ctx, chromedp.Text("a", &text, chromedp.ByQuery, chromedp.FromNode(a)))
							fmt.Printf("[linkedin_scraper] li %d anchor: href=%s, text=%s\n", i, href, text)
							if strings.HasPrefix(href, "https://www.linkedin.com/in/") && profileURL == "" {
								profileURL = strings.Split(href, "?")[0]
								name = text
							}
						}
						chromedp.Run(ctx, chromedp.Text("span.entity-result__badge", &degree, chromedp.ByQuery, chromedp.FromNode(node)))
						chromedp.Run(ctx, chromedp.Text("div.entity-result__primary-subtitle", &location, chromedp.ByQuery, chromedp.FromNode(node)))
						fmt.Printf("[linkedin_scraper] Node %d: name=%s, url=%s, degree=%s, location=%s\n", i, name, profileURL, degree, location)
						connections = append(connections, LinkedInConnection{
							Name:             name,
							ProfileURL:       profileURL,
							ConnectionDegree: degree,
							Location:         location,
						})
					}
					fmt.Println("[linkedin_scraper] Scraping complete, returning results.")
					return connections, nil
				}
			}
		}
	}
	fmt.Println("[linkedin_scraper] Cookie-based auth failed or not present, falling back to login flow.")

	fmt.Println("[linkedin_scraper] Logging in...")
	if err := loginLinkedIn(ctx); err != nil {
		fmt.Println("[linkedin_scraper] Login error:", err)
		return nil, err
	}
	fmt.Println("[linkedin_scraper] Login successful, navigating to:", url)

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`ul[role='list']`, chromedp.ByQuery),
	); err != nil {
		fmt.Println("[linkedin_scraper] Navigation or wait error:", err)
		return nil, err
	}
	fmt.Println("[linkedin_scraper] Navigation successful, scraping nodes...")

	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Nodes(`ul[role='list'] > li`, &nodes, chromedp.ByQueryAll),
	); err != nil {
		fmt.Println("[linkedin_scraper] Node scrape error:", err)
		return nil, err
	}
	fmt.Printf("[linkedin_scraper] Found %d connection nodes\n", len(nodes))

	for i, node := range nodes {
		var name, profileURL, degree, location string
		var anchors []*cdp.Node
		chromedp.Run(ctx, chromedp.Nodes("a", &anchors, chromedp.ByQueryAll, chromedp.FromNode(node)))
		for _, a := range anchors {
			var href, text string
			chromedp.Run(ctx, chromedp.AttributeValue("a", "href", &href, nil, chromedp.ByQuery, chromedp.FromNode(a)))
			chromedp.Run(ctx, chromedp.Text("a", &text, chromedp.ByQuery, chromedp.FromNode(a)))
			fmt.Printf("[linkedin_scraper] li %d anchor: href=%s, text=%s\n", i, href, text)
			if strings.HasPrefix(href, "https://www.linkedin.com/in/") && profileURL == "" {
				profileURL = strings.Split(href, "?")[0]
				name = text
			}
		}
		chromedp.Run(ctx, chromedp.Text("span.entity-result__badge", &degree, chromedp.ByQuery, chromedp.FromNode(node)))
		chromedp.Run(ctx, chromedp.Text("div.entity-result__primary-subtitle", &location, chromedp.ByQuery, chromedp.FromNode(node)))
		fmt.Printf("[linkedin_scraper] Node %d: name=%s, url=%s, degree=%s, location=%s\n", i, name, profileURL, degree, location)
		connections = append(connections, LinkedInConnection{
			Name:             name,
			ProfileURL:       profileURL,
			ConnectionDegree: degree,
			Location:         location,
		})
	}
	fmt.Println("[linkedin_scraper] Scraping complete, returning results.")
	return connections, nil
}
