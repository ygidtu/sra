package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/voxelbrain/goptions"
)

func loadRBPs(path string) []string {
	res := make([]string, 0)
	rbps := map[string]int{}
	if _, err := os.Stat((path)); !os.IsNotExist(err) {
		file, err := os.Open(path) //
		if err != nil {
			log.Fatal("找不到CSV檔案路徑:", path, err)
		}

		// read
		r := csv.NewReader(file)

		header := []string{}
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}

			if len(header) < 1 {
				header = record
				continue
			}

			for _, i := range record {
				if i != "Yes" && i != "No" && i != "" {
					rbps[i] = 0
				}
			}
		}
	}

	for i := range rbps {
		res = append(res, i)
	}
	sort.Strings(res)
	return res
}

func main() {

	options := struct {
		RBP     string        `goptions:"-i, --input, obligatory, description='RBP的list，csv格式，带列名'"`
		SRA     string        `goptions:"-u, --url, description='SRA的官方链接'"`
		Proxy   string        `goptions:"-x, --proxy, description='代理链接地址，比如：http://127.0.0.1:7890'"`
		Output  string        `goptions:"-o, --output, description='输出文件夹'"`
		Param   string        `goptions:"-p, --param, description='额外的查询参数'"`
		Timeout time.Duration `goptions:"-t, --timeout, description='Connection timeout in seconds'"`
		Open    bool          `goptions:"--open, description='是否打开chrome的图形化界面'"`
		Help    goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		SRA:     "https://www.ncbi.nlm.nih.gov/sra/",
		Output:  "./output",
		Timeout: 10 * time.Second,
		Param:   `"Homo sapiens"[orgn:__txid9606] AND(rna seq[Strategy])`,
	}

	goptions.ParseAndFail(&options)

	if _, err := os.Stat(options.Output); os.IsNotExist(err) {
		err = os.MkdirAll(options.Output, 0777)
		if err != nil {
			log.Fatal("无法创建输出文件夹", err)
		}
	}

	var err error
	KEYWORD := []string{"knockdown", "knockout", "overexpression"}
	RBPs := loadRBPs(options.RBP)

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	defaultDownload := filepath.Join(user.HomeDir, "Downloads")

	// create context
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", !options.Open),
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	}

	if options.Proxy != "" {
		opts = append(
			opts,
			chromedp.ProxyServer(options.Proxy),
		)
	}

	contextOpts := []chromedp.ContextOption{
		chromedp.WithLogf(log.Printf),
		// chromedp.WithDebugf(log.Printf),
		chromedp.WithErrorf(log.Printf),
	}
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocContext, contextOpts...)
	defer cancel()

	err = chromedp.Run(ctx,
		chromedp.Navigate(options.SRA),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, rbp := range RBPs {
		log.Println(rbp)
		for _, i := range KEYWORD {
			log.Println(i)
			output := filepath.Join(options.Output, fmt.Sprintf("%v_%v.csv", rbp, i))

			if _, err := os.Stat(output); !os.IsNotExist(err) {
				continue
			}
			term := fmt.Sprintf("(\"%s\"[Title] OR \"%s\"[Description]) AND (\"%s\"[Title] OR \"%s\"[Description]) AND %s", i, i, rbp, rbp, options.Param)
			// log.Println(term)
			err = chromedp.Run(ctx,
				chromedp.SetValue(`#term`, term, chromedp.ByID),
				chromedp.Sleep(3*time.Second),
				chromedp.Click("#search", chromedp.ByID),
				chromedp.Sleep(options.Timeout),
			)
			if err != nil {
				log.Fatal(err)
			}

			// err = chromedp.Run(ctx,
			// 	chromedp.ActionFunc(func(ctx context.Context) error {
			// 		var nodes []*cdp.Node
			// 		if err := chromedp.Nodes("#sendto", &nodes, chromedp.AtLeast(0)).Do(ctx); err != nil {
			// 			return err
			// 		}
			// 		if len(nodes) == 0 {
			// 			return err
			// 		} // nothing to do
			// 		return chromedp.Evaluate(`var h4 = document.getElementById("sendto"); h4.click()`, nil).Do(ctx)
			// 	}),
			// )

			if err != nil {
				log.Println(err)
				continue
			}

			err = chromedp.Run(ctx,
				chromedp.Evaluate(`var h4 = document.getElementById("sendto"); h4.click()`, nil),
				chromedp.Sleep(1*time.Second),
				chromedp.Click("#dest_File", chromedp.ByID),
				chromedp.Sleep(1*time.Second),
				chromedp.Click(`//div[@id="submenu_File"]/button`, chromedp.BySearch),
			)

			if err != nil {
				log.Println(err)
				continue
			}

			sra_result := filepath.Join(defaultDownload, "sra_result.csv")

			if _, err := os.Stat((sra_result)); !os.IsNotExist(err) {
				os.Rename(sra_result, output)
			}
		}
		time.Sleep(3 * time.Second)
	}
}
