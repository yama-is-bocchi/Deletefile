package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Path struct {
	Dirpath string
}

// フィルター
type Filter func(string) bool

// DirPath,FileNameをフィルターする関数
func FilterName(name string, filter Filter) bool {

	if filter(name) == true {
		fmt.Printf("無効な文字列です\n")
		return false
	}

	return true
}

// ディレクトリ名
// true==含んでいる,false==含んでいない
func DirFilter() Filter {
	return func(n string) bool {
		if strings.Contains(n, "*") ||
			strings.Contains(n, "?") || strings.Contains(n, "<") ||
			strings.Contains(n, ">") || strings.Contains(n, "\"") {
			return true
		}
		return false
	}
}

// 日付
// true==数値以外を含んでいる,false==含んでいない
func DateFilter() Filter {
	return func(n string) bool {
		match, _ := regexp.MatchString("^[0-9]+$", n)
		if !match {
			return true
		}
		return false
	}
}

// 文字列を入力して返す
// 戻り値input or ""(空文字)
func enter_str() string {
	var input string  //入力文字列
	var ok_str string //入力確認文字列

	for {
		// Enterが押されるまで入力を待つ
		fmt.Scanln(&input)
		if len(input) > 0 {
			break
		}

		fmt.Printf("文字列の長さが0です\n")
	}
	for {
		fmt.Printf("入力は間違いないですか?1=OK,0=やり直し\n")
		fmt.Scanln(&ok_str)
		if ok_str == "1" {
			return input
		} else if ok_str == "0" {
			break
		}
		fmt.Printf("入力文字が無効です\n")
	}
	return ""
}

// 入力を続けるか確認
// True続ける,False終了
func continue_check() bool {
	fmt.Printf("入力を終了しますか?1=終了,0=続ける\n")
	var ok_str string //入力確認文字列
	for {

		fmt.Scanln(&ok_str)
		if ok_str == "1" {
			return true
		} else if ok_str == "0" {
			return false
		}
		fmt.Printf("入力文字が無効です\n")
	}
}

// 引数の文字列をYYYY-MM-DDの形に変更して返す
func date_converter(date_str string) string {
	y := date_str[:4] + "-"
	temp := date_str[4:]
	m := temp[:2] + "-"
	d := temp[2:]
	var ret string
	ret = y + m + d
	return ret
}

func main() {

	for {
		// フォルダのパスを指定します
		var FolderList []Path
		var folderPath string //フォルダのパス
		var setdate string
		var count int
		count=0

		for {
			pathFilter := DirFilter()
			fmt.Printf("フォルダのパスを入力してください\n")
			folderPath = enter_str()
			if len(folderPath) > 0 && FilterName(folderPath, pathFilter) == true {
				temp := Path{Dirpath: folderPath}
				FolderList = append(FolderList, temp)
				//複数フォルダーを選択するか
				if continue_check() == true {
					break
				}
			}
		}

		for {
			dateFilter := DateFilter()
			fmt.Printf("閾値の日付を入力してください(※YYYYMMDD形式でゼロ埋めしてください)\n")
			setdate = enter_str()
			if len(setdate) == 8 && FilterName(setdate, dateFilter) == true {
				//YYYY-MM-DD型に変換する
				setdate = date_converter(setdate)
				// 文字列をtime.Time型に変換
				dateFormat := "2006-01-02" // 年-月-日のフォーマット
				_, err := time.Parse(dateFormat, setdate)

				if err != nil {
					fmt.Printf("Date型に変換できません\n")
				}else{
					break
				}
			}
		}

		for _, dir := range FolderList {
			// フォルダ内のファイル一覧を取得します
			files, err := filepath.Glob(filepath.Join(dir.Dirpath, "*"))
			if err != nil {
				fmt.Println("ファイルが存在しません", err)
				break
			}
			for _, file := range files {

				info, err := os.Stat(file)
				if err != nil {
					fmt.Println("エラー:", err)
					continue
				}
				//フォルダではない
				if !info.IsDir() {
					//更新日時が指定日付より前か
					// 更新日時を取得
					modifiedTime := info.ModTime()

					// 日付だけをフォーマットして出力
					dateFormat := "2006-01-02" // 年-月-日のフォーマット
					dateOnly := modifiedTime.Format(dateFormat)

					//閾値より下の更新日付のファイルを削除
					ntime, err := time.Parse(dateFormat, setdate)
					if err != nil {
						fmt.Println("文字列をtime.Time型に変換できません:", err)
						return
					}
					if dateOnly < ntime.Format(dateFormat) {
						err := os.Remove( file)
						if err != nil {
							fmt.Println("ファイルを削除できません", err)
							return
						}
						count++
					}

				}
			}
		}
		fmt.Println(count,"件のファイルを削除しました")
		if continue_check() == true {
			break
		}
	}
}
