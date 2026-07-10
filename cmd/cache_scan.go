package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type cacheDir struct {
	path  string
	size  int64
	items int
}

var (
	cacheMinSizeMB int64
	cacheDeleteAll bool
	cacheDryRun    bool
)

var cacheScanCmd = &cobra.Command{
	Use:   "cache-scan",
	Short: "Scan dan hapus cache filesystem yang tidak penting",
	Long: `Memindai direktori cache umum (seperti ~/Library/Caches, ~/.cache, dll)
dan menampilkan folder cache beserta ukurannya. Pengguna dapat memilih
folder mana yang akan dihapus untuk membebaskan storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		runCacheScan()
	},
}

func init() {
	generateCmd.AddCommand(cacheScanCmd)
	cacheScanCmd.Flags().Int64VarP(&cacheMinSizeMB, "min-size", "m", 10, "Ukuran minimum (MB) untuk ditampilkan")
	cacheScanCmd.Flags().BoolVarP(&cacheDeleteAll, "all", "a", false, "Hapus semua tanpa konfirmasi")
	cacheScanCmd.Flags().BoolVarP(&cacheDryRun, "dry-run", "n", false, "Hanya tampilkan tanpa menghapus")
}

func runCacheScan() {
	minBytes := cacheMinSizeMB * 1024 * 1024

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	roots := scanRoots(home)

	var entries []cacheDir
	for _, root := range roots {
		dirs, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, d := range dirs {
			if !d.IsDir() {
				continue
			}
			dirPath := filepath.Join(root, d.Name())
			size, items := dirSize(dirPath)
			if size >= minBytes {
				entries = append(entries, cacheDir{path: dirPath, size: size, items: items})
			}
		}
	}

	if len(entries) == 0 {
		fmt.Println("Tidak ada cache yang cukup besar untuk ditampilkan.")
		fmt.Printf("(Gunakan --min-size untuk mengubah batas, saat ini %d MB)\n", cacheMinSizeMB)
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].size > entries[j].size
	})

	var totalSize int64
	var totalItems int
	div := strings.Repeat("─", 80)
	fmt.Println("\nDaftar cache yang ditemukan:")
	fmt.Println(div)
	fmt.Printf("%-3s %-55s %-12s %s\n", "#", "Path", "Size", "Files")
	fmt.Println(div)
	for i, e := range entries {
		fmt.Printf("%-3d %-55s %-12s %d\n", i+1, shortenPath(e.path, home), humanSize(e.size), e.items)
		totalSize += e.size
		totalItems += e.items
	}
	fmt.Println(div)
	fmt.Printf("    Total: %s, %d file\n\n", humanSize(totalSize), totalItems)

	if cacheDryRun {
		fmt.Println("Dry-run: tidak ada yang dihapus.")
		return
	}

	if cacheDeleteAll {
		deleteEntries(entries)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Pilih nomor (pisahkan koma, contoh: 1,3,5), 'a' untuk semua, 'q' untuk keluar: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "q", "Q", "":
		fmt.Println("Dibatalkan.")
		return
	case "a", "A", "all":
		deleteEntries(entries)
		return
	}

	var selected []cacheDir
	parts := strings.Split(input, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		idx, err := strconv.Atoi(p)
		if err != nil || idx < 1 || idx > len(entries) {
			fmt.Printf("Input tidak valid: %s\n", p)
			continue
		}
		selected = append(selected, entries[idx-1])
	}

	if len(selected) == 0 {
		fmt.Println("Tidak ada yang dipilih.")
		return
	}
	deleteEntries(selected)
}

func scanRoots(home string) []string {
	var roots []string
	candidates := []string{
		filepath.Join(home, "Library", "Caches"),
		filepath.Join(home, ".cache"),
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			roots = append(roots, c)
		}
	}
	return roots
}

func dirSize(path string) (int64, int) {
	var size int64
	var count int
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
			count++
		}
		return nil
	})
	return size, count
}

func humanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	exp := 0
	for n := bytes; n >= unit; n /= unit {
		exp++
	}
	div := int64(1)
	for i := 0; i < exp; i++ {
		div *= unit
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp-1])
}

func shortenPath(path, home string) string {
	if strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}

func deleteEntries(selected []cacheDir) {
	var freed int64
	var count int
	for _, e := range selected {
		fmt.Printf("Menghapus %s ... ", shortenPath(e.path, mustHome()))
		if err := os.RemoveAll(e.path); err != nil {
			fmt.Printf("Gagal: %v\n", err)
			continue
		}
		freed += e.size
		count++
		fmt.Println("OK")
	}
	fmt.Printf("\nSelesai! %d folder dihapus, %s storage dibebaskan.\n", count, humanSize(freed))
}

func mustHome() string {
	h, _ := os.UserHomeDir()
	return h
}
