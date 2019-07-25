package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	comics "github.com/fancxxy/gocomic/download/comic"
	"github.com/fancxxy/gocomic/download/parser"
	"github.com/fancxxy/gocomic/download/rpc"
)

func main() {
	help := flag.Bool("help", false, "command line tool to download comics")
	website := flag.String("website", "", "webiste name")
	comic := flag.String("comic", "", "comic title")
	chapter := flag.String("chapter", "", "chapter index")
	path := flag.String("path", "", "set download directory")
	latest := flag.Bool("latest", false, "check the latest updated chapter")
	serve := flag.Bool("serve", false, "start rpc server")

	flag.Parse()

	if *help {
		h()
		os.Exit(0)
	}

	if *serve {
		if err := rpc.Start(); err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}

	if *website == "" || *comic == "" {
		fmt.Printf("website and comic are mandatory\n")
		os.Exit(-1)
	}

	var chapters []int
	if *chapter != "" {
		nums := strings.Split(*chapter, ",")
		for _, num := range nums {
			chapter, err := strconv.Atoi(num)
			if err != nil {
				fmt.Printf("%d is not a valid chapter index\n", chapter)
				continue
			}
			chapters = append(chapters, chapter)
		}
	}

	url, err := comics.SearchComic(*website, *comic)
	if url == "" {
		fmt.Printf("search %s in %s error: %v\n", *comic, *website, err)
		os.Exit(-1)
	}

	instance, err := comics.NewComic(url)
	if err != nil {
		fmt.Printf("get comic %s information error: %v\n", *comic, err)
		os.Exit(-1)
	}

	if *latest {
		fmt.Printf("%s %s\n", instance.Latest, instance.Date)
	}

	if *chapter == "" && *latest {
		return
	}

	if *path != "" {
		_, err = instance.Home(*path)
		if err != nil {
			fmt.Printf("set download directory error: %v\n", err)
			os.Exit(-1)
		}
	}

	instance.DownloadInCmd(chapters...)
	fmt.Printf("finished\n")
}

func h() {
	fmt.Printf("Usage: gocomic [-website -comic [-chapter...] [-latest] [-path]] [-serve]\n\n")
	fmt.Printf("gocomic is a command line tool to download specific comic or chapter resources\n\n")
	fmt.Printf("Support website:\n")
	for _, label := range parser.Support() {
		fmt.Printf("  %s\n", label)
	}
	fmt.Printf("\nExamples:\n")
	fmt.Printf("  gocomic -website 腾讯动漫 -comic 航海王              		# download all chapters\n")
	fmt.Printf("  gocomic -website 腾讯动漫 -comic 航海王 -chapter 936,937,938	# download chaper 936, 937 and 938\n")
	fmt.Printf("  gocomic -website 腾讯动漫 -comic 航海王 -latest		# check the latest chapter\n")
	fmt.Printf("  gocomic -website 腾讯动漫 -comic 航海王 -path ~/Comics/ 	# download all chapters to ~/Comics/ directory\n")
	fmt.Printf("  gocomic -serve						# start rpc server\n")
}
