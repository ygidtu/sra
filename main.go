package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/voxelbrain/goptions"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

func loadRBPs(path string) []string {
	res := make([]string, 0)
	rbps := map[string]int{}
	if _, err := os.Stat((path)); !os.IsNotExist(err) {
		file, err := os.Open(path) //
		if err != nil {
			log.Error("找不到CSV檔案路徑:", path, err)
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
				log.Error(err)
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
		Debug   bool          `goptions:"--debug, description='显示debug log'"`
		Help    goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		SRA:     "https://www.ncbi.nlm.nih.gov/sra/",
		Output:  "./output",
		Timeout: 10 * time.Second,
		Param:   `"Homo sapiens"[orgn:__txid9606] AND(rna seq[Strategy])`,
	}

	goptions.ParseAndFail(&options)

	// setup logger
	SetLogger(options.Debug)

	// 检查输出文件夹
	if _, err := os.Stat(options.Output); os.IsNotExist(err) {
		err = os.MkdirAll(options.Output, 0777)
		if err != nil {
			log.Error("无法创建输出文件夹", err)
		}
	}

	progress := &Progress{
		Path: filepath.Join(options.Output, "progress.json"),
		Data: map[string]int{},
	}
	progress.Load()

	KEYWORD := []string{"knockdown", "knockout", "overexpression"}
	RBPs := loadRBPs(options.RBP)

	// 默认下载目录和文件
	user, err := user.Current()
	if err != nil {
		log.Error(err)
	}

	defaultDownload := filepath.Join(user.HomeDir, "Downloads")
	sra_result := filepath.Join(defaultDownload, "sra_result.csv")

	if _, ok := os.Stat(sra_result); !os.IsNotExist(ok) {
		os.Remove(sra_result)
	}

	// create context
	log.Debug("Start chrome under headless ? ", !options.Open)
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
		chromedp.WithLogf(log.Infof),
		chromedp.WithDebugf(log.Debugf),
		chromedp.WithErrorf(log.Infof),
	}
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocContext, contextOpts...)
	defer cancel()

	log.Debug("Open: ", options.SRA)
	if err := chromedp.Run(ctx,
		chromedp.Navigate(options.SRA),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		log.Error(err)
	}

	for idx, rbp := range RBPs {
		log.Infof("[%d/%d] %s", idx+1, len(RBPs), rbp)
		for _, i := range KEYWORD {
			log.Info(i)

			name := fmt.Sprintf("%v_%v.csv", rbp, i)
			output := filepath.Join(options.Output, name)

			if progress.Has(name) {
				log.Info("already finished, skip")
				continue
			}

			// 检查文件是否已存在
			if _, err := os.Stat(output); !os.IsNotExist(err) {
				continue
			}

			// 构建查询语句
			term := fmt.Sprintf("(\"%s\"[Title] OR \"%s\"[Description]) AND (\"%s\"[Title] OR \"%s\"[Description]) AND %s", i, i, rbp, rbp, options.Param)

			// 输入查询语句，查询
			log.Debug("Term: ", term)
			if err := chromedp.Run(ctx,
				chromedp.WaitReady("#term", chromedp.ByID),
				chromedp.SetValue("#term", term, chromedp.ByID),
				chromedp.Sleep(3*time.Second),
				chromedp.WaitReady("#search", chromedp.ByID),
				chromedp.Click("#search", chromedp.ByID),
				chromedp.Sleep(options.Timeout),
			); err != nil {
				log.Error("search failed", err)
			}

			// 根据页面title判断是否有结果
			title := ""
			if err := chromedp.Run(ctx, chromedp.Title(&title)); err != nil {
				log.Error("failed to get title from page", err)
			}

			if strings.Contains(title, "No items found") {
				log.Infof("No items found, next")
				progress.Add(name)
				progress.Dump()
				continue
			}

			log.Debug("Wait for click")
			if err := chromedp.Run(ctx,
				chromedp.WaitReady("#sendto", chromedp.ByID),
				chromedp.Evaluate(`var h4 = document.getElementById("sendto"); h4.click()`, nil),
				chromedp.Sleep(1*time.Second),
				chromedp.WaitReady("#dest_File", chromedp.ByID),
				chromedp.Click("#dest_File", chromedp.ByID),
				chromedp.Sleep(1*time.Second),
				chromedp.WaitReady(`//div[@id="submenu_File"]/button`, chromedp.BySearch),
				chromedp.Click(`//div[@id="submenu_File"]/button`, chromedp.BySearch),
			); err != nil {
				log.Error("failed to click download", err)
			}

			time.Sleep(2 * time.Second)

			_, err = os.Stat((sra_result))
			// 只有下载完成才能推出循环
			for os.IsNotExist(err) {
				log.Info("wait for download")
				time.Sleep(3 * time.Second)
				_, err = os.Stat((sra_result))
			}
			os.Rename(sra_result, output)
			progress.Add(name)
			progress.Dump()
		}
		time.Sleep(3 * time.Second)
	}
}

// IGF2BP2 knockdown
// ("knockdown"[Title] OR "knockdown"[Description]) AND ("IGF2BP2"[Title] OR "IGF2BP2"[Description]) AND "Homo sapiens"[orgn:__txid9606] AND(rna seq[Strategy])
