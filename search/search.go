package search

import (
	"encoding/csv"
	"fmt"
	"github.com/chromedp/cdproto/browser"
	"github.com/ygidtu/sra/client"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger
)

func loadRBPs(path string) []string {

	sugar.Infof("Load RBPs from %s", path)
	res := make([]string, 0)
	RBPs := map[string]int{}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		file, err := os.Open(path) //
		if err != nil {
			sugar.Fatalf("找不到CSV檔案路徑: %v, %v", path, err)
		}

		// read
		r := csv.NewReader(file)

		var header []string
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				sugar.Error(err)
			}

			sugar.Debugf("%v", record)
			if len(header) < 1 {
				header = record
				continue
			}

			for _, i := range record {
				if i != "Yes" && i != "No" && i != "" {
					RBPs[i] = 0
				}
			}
		}
	}

	for i := range RBPs {
		res = append(res, i)
	}
	sort.Strings(res)
	return res
}

func Search(options *Params, sugar_ *zap.SugaredLogger) {
	sugar = sugar_
	sugar.Info(options.String())

	// 检查输出文件夹
	if _, err := os.Stat(options.Output); os.IsNotExist(err) {
		err = os.MkdirAll(options.Output, 0777)
		if err != nil {
			sugar.Error("无法创建输出文件夹", err)
		}
	}

	progress := &Progress{
		Path: filepath.Join(options.Output, "progress.json"),
		Data: map[string]int{},
	}
	progress.Load()

	KEYWORD := []string{
		"knockdown", "knock-down", "knock down",
		"knockout", "knock out", "knock-out",
		"overexpression", "over-expression", "over expression",
		"siRNA", "shRNA", "dCas9", "sgRNA",
		"crispr cas9", "crispr-cas9",
	}
	RBPs := loadRBPs(options.RBP)

	sugar.Infof("there are %d RBPs", len(RBPs))

	sraResult := filepath.Join(options.Output, "sra_result.csv")

	if _, ok := os.Stat(sraResult); !os.IsNotExist(ok) {
		_ = os.Remove(sraResult)
	}

	// create context
	ctx, cancel := client.SetChromeClient(options.Open, options.Proxy, options.Exec, sugar_)
	defer cancel()

	sugar.Debug("Open: ", options.SRA)
	if err := chromedp.Run(ctx,
		chromedp.Navigate(options.SRA),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		sugar.Error(err)
	}

	for idx, rbp := range RBPs {
		for _, i := range KEYWORD {
			sugar.Infof("[%d/%d] %s - %s", idx+1, len(RBPs), rbp, i)

			name := fmt.Sprintf("%v_%v.csv", rbp, i)
			output := filepath.Join(options.Output, name)

			if progress.Has(name) {
				sugar.Info("already finished, skip")
				continue
			}

			// 检查文件是否已存在
			if _, err := os.Stat(output); !os.IsNotExist(err) {
				continue
			}

			// 构建查询语句
			term := fmt.Sprintf("(\"%s\"[Title] OR \"%s\"[Description]) AND (\"%s\"[Title] OR \"%s\"[Description]) AND %s", i, i, rbp, rbp, options.Param)

			// 输入查询语句，查询
			sugar.Debug("Term: ", term)
			if err := chromedp.Run(ctx,
				chromedp.WaitReady("#term", chromedp.ByID),
				chromedp.SetValue("#term", term, chromedp.ByID),
				//chromedp.Sleep(3*time.Second),
				chromedp.WaitReady("#search", chromedp.ByID),
				chromedp.Click("#search", chromedp.ByID),
			); err != nil {
				sugar.Fatal("search failed", err)
			}

			time.Sleep(options.Timeout)

			// 根据页面title判断是否有结果
			title := ""
			if err := chromedp.Run(ctx, chromedp.Title(&title)); err != nil {
				sugar.Fatal("failed to get title from page", err)
			}

			for !strings.Contains(title, fmt.Sprintf("(\"%s\"[Title] OR \"%s\"[Description])", i, i)) && !strings.Contains(title, "No items found") {

				if strings.HasPrefix(title, "GSM") {
					break
				}

				sugar.Warn("page title contains neither term nor 'no items found', is page still loading ?")
				//time.Sleep(3 * time.Second)

				if err := chromedp.Run(ctx, chromedp.Title(&title)); err != nil {
					sugar.Fatal("failed to get title from page", err)
				}
			}

			if strings.Contains(title, "No items found") {
				sugar.Infof("No items found, next")
				progress.Add(name)
				_ = progress.Dump()
				continue
			}

			sugar.Debug("Wait for click")
			if err := chromedp.Run(ctx,
				chromedp.WaitReady("#sendto", chromedp.ByID),
				chromedp.Evaluate(`var h4 = document.getElementById("sendto"); h4.click()`, nil),
				chromedp.Sleep(3*time.Second),
				chromedp.WaitReady("#dest_File", chromedp.ByID),
				chromedp.Click("#dest_File", chromedp.ByID),
				chromedp.Sleep(1*time.Second),
				browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
					WithDownloadPath(options.Output).
					WithEventsEnabled(true),
				chromedp.WaitReady(`//div[@id="submenu_File"]/button`, chromedp.BySearch),
				chromedp.Click(`//div[@id="submenu_File"]/button`, chromedp.BySearch),
			); err != nil {
				sugar.Fatal("failed to click download", err)
			}

			//time.Sleep(2 * time.Second)

			_, err := os.Stat(sraResult)
			// 只有下载完成才能推出循环
			for os.IsNotExist(err) {
				sugar.Info("wait for download")
				time.Sleep(3 * time.Second)
				_, err = os.Stat(sraResult)
			}

			sugar.Debug("rename sra_result to ", output)
			_ = os.Rename(sraResult, output)
			progress.Add(name)
			_ = progress.Dump()
			//time.Sleep(3 * time.Second)
		}
	}
}
