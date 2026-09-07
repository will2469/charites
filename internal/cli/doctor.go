package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// CheckStatus mendefinisikan level status hasil diagnosis.
type CheckStatus string

// Nilai kanonikal status pemeriksaan dokter.
const (
	StatusPass CheckStatus = "PASS"
	StatusWarn CheckStatus = "WARN"
	StatusFail CheckStatus = "FAIL"
)

// DoctorItem merepresentasikan hasil pengujian tunggal dokter sistem.
type DoctorItem struct {
	Category string
	Title    string
	Status   CheckStatus
	Message  string
	Details  []string
	Remedy   string
}

// RunDoctor mengorkestrasi diagnosis kesehatan instalasi, PATH, izin berkas, dan konektivitas.
func RunDoctor(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	var verbose bool
	fs.BoolVar(&verbose, "v", false, "Tampilkan rincian diagnosis lengkap")
	fs.BoolVar(&verbose, "verbose", false, "Tampilkan rincian diagnosis lengkap")

	if err := fs.Parse(args); err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: %v. Run 'charites --help' for usage.\n", err)
		return ExitOperational
	}

	_, _ = fmt.Fprintln(stdout, "Charites Doctor (v"+Version+")")
	_, _ = fmt.Fprintln(stdout, "Auditing runtime environment, PATH health, file permissions, and connectivity...")
	_, _ = fmt.Fprintln(stdout, "----------------------------------------------------------------------")

	var items []DoctorItem

	// 1. Audit Versi dan Arsitektur Host
	items = append(items, checkHostRuntime())

	// 2. Audit $PATH dan Tabrakan Biner
	items = append(items, checkPathAndBinaries()...)

	// 3. Audit Izin Direktori untuk In-Place Self-Update
	items = append(items, checkBinaryPermissions())

	// 4. Audit Konektivitas API GitHub
	items = append(items, checkUpdateConnectivity())

	// 5. Audit Proyek Lokal / Workspace (opsional jika dijalankan dalam repositori)
	if wsItem, ok := checkWorkspaceConfig(); ok {
		items = append(items, wsItem)
	}

	// Render Hasil Diagnosis
	failCount := 0
	warnCount := 0

	for _, item := range items {
		symbol := "[]"
		switch item.Status {
		case StatusWarn:
			symbol = "[!]"
			warnCount++
		case StatusFail:
			symbol = "[]"
			failCount++
		}

		_, _ = fmt.Fprintf(stdout, "%-3s %-28s %s\n", symbol, item.Title+":", item.Message)
		for _, d := range item.Details {
			_, _ = fmt.Fprintf(stdout, "    • %s\n", d)
		}
		if item.Remedy != "" && (item.Status == StatusWarn || item.Status == StatusFail) {
			_, _ = fmt.Fprintf(stdout, "    ↳ Remediasi: %s\n", item.Remedy)
		}
	}

	_, _ = fmt.Fprintln(stdout, "----------------------------------------------------------------------")

	if failCount > 0 {
		_, _ = fmt.Fprintf(stdout, "Doctor Verdict: %d error(s) and %d warning(s) detected. Please follow remedies above.\n", failCount, warnCount)
		return ExitOperational
	}

	if warnCount > 0 {
		_, _ = fmt.Fprintf(stdout, "Doctor Verdict: Charites is operational with %d warning(s). Review recommendations above.\n", warnCount)
		return ExitClean
	}

	_, _ = fmt.Fprintln(stdout, "Doctor Verdict: Everything looks healthy! No issues found.")
	return ExitClean
}

func checkHostRuntime() DoctorItem {
	return DoctorItem{
		Category: "Runtime",
		Title:    "Host Platform",
		Status:   StatusPass,
		Message:  fmt.Sprintf("%s/%s (Go runtime %s)", runtime.GOOS, runtime.GOARCH, runtime.Version()),
	}
}

func checkPathAndBinaries() []DoctorItem {
	rawPath := os.Getenv("PATH")
	if rawPath == "" {
		return []DoctorItem{{
			Category: "PATH",
			Title:    "$PATH Configuration",
			Status:   StatusFail,
			Message:  "$PATH is empty or not set",
			Remedy:   "Set your $PATH environment variable in your shell profile (~/.bashrc or ~/.zshrc)",
		}}
	}

	dirs := filepath.SplitList(rawPath)
	items := make([]DoctorItem, 0, 3)

	// A. Duplikasi direktori di $PATH
	items = append(items, checkPathDuplication(dirs))

	// B. Tabrakan biner di $PATH
	binItem, foundBins := checkBinaryCollisions(dirs)
	items = append(items, binItem)

	// C. Mismatch biner aktif
	if mismatchItem, ok := checkActiveBinaryMismatch(foundBins); ok {
		items = append(items, mismatchItem)
	}

	return items
}

func checkPathDuplication(dirs []string) DoctorItem {
	seenDirs := make(map[string]int)
	var dupList []string
	for _, d := range dirs {
		cleaned := filepath.Clean(d)
		seenDirs[cleaned]++
		if seenDirs[cleaned] == 2 {
			dupList = append(dupList, cleaned)
		}
	}

	if len(dupList) == 0 {
		return DoctorItem{
			Category: "PATH",
			Title:    "$PATH Cleanliness",
			Status:   StatusPass,
			Message:  "Tidak ada entri duplikat pada $PATH",
		}
	}

	details := make([]string, 0, len(dupList))
	for _, dup := range dupList {
		details = append(details, fmt.Sprintf("%s tercantum %d kali", dup, seenDirs[dup]))
	}

	return DoctorItem{
		Category: "PATH",
		Title:    "$PATH Duplication",
		Status:   StatusWarn,
		Message:  fmt.Sprintf("Found %d duplicate directory entries in $PATH", len(dupList)),
		Details:  details,
		Remedy:   "Hapus pemanggilan berulang 'export PATH=...' di ~/.bashrc, ~/.zshrc, atau ~/.profile",
	}
}

func checkBinaryCollisions(dirs []string) (DoctorItem, []string) {
	var foundBinaries []string
	seenBinPaths := make(map[string]bool)

	binaryName := "charites"
	if runtime.GOOS == "windows" {
		binaryName = "charites.exe"
	}

	for _, dir := range dirs {
		candidate := filepath.Join(dir, binaryName)
		if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
			realPath, rErr := filepath.EvalSymlinks(candidate)
			if rErr != nil {
				realPath = candidate
			}
			if !seenBinPaths[realPath] {
				seenBinPaths[realPath] = true
				foundBinaries = append(foundBinaries, candidate)
			}
		}
	}

	switch len(foundBinaries) {
	case 0:
		return DoctorItem{
			Category: "Binary",
			Title:    "Binary Discovery",
			Status:   StatusWarn,
			Message:  "Biner 'charites' tidak ditemukan di direktori $PATH aktif saat ini",
			Remedy:   "Tambahkan direktori instalasi (~/.local/bin atau /usr/local/bin) ke $PATH",
		}, foundBinaries
	case 1:
		return DoctorItem{
			Category: "Binary",
			Title:    "Binary Discovery",
			Status:   StatusPass,
			Message:  fmt.Sprintf("Tunggal dan unik di %s", foundBinaries[0]),
		}, foundBinaries
	default:
		details := make([]string, 0, len(foundBinaries))
		for i, b := range foundBinaries {
			prefix := "Aktif (utama)"
			if i > 0 {
				prefix = "Tertimpa (shadowed)"
			}
			details = append(details, fmt.Sprintf("[%s] %s", prefix, b))
		}
		return DoctorItem{
			Category: "Binary",
			Title:    "Multiple Binaries",
			Status:   StatusWarn,
			Message:  fmt.Sprintf("Ditemukan %d biner berbeda di direktori $PATH", len(foundBinaries)),
			Details:  details,
			Remedy:   fmt.Sprintf("Hapus biner sekunder selain %s untuk mencegah inkonsistensi versi", foundBinaries[0]),
		}, foundBinaries
	}
}

func checkActiveBinaryMismatch(foundBinaries []string) (DoctorItem, bool) {
	if len(foundBinaries) == 0 {
		return DoctorItem{}, false
	}

	execPath, err := OsExecutable()
	if err != nil {
		return DoctorItem{}, false
	}

	realExec, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		realExec = execPath
	}

	realFirst, err := filepath.EvalSymlinks(foundBinaries[0])
	if err != nil {
		realFirst = foundBinaries[0]
	}

	if realExec != realFirst {
		return DoctorItem{
			Category: "Binary",
			Title:    "Active Binary Mismatch",
			Status:   StatusWarn,
			Message:  "Biner yang sedang berjalan berbeda dari prioritas pertama di $PATH",
			Details: []string{
				"Sedang berjalan: " + realExec,
				"Prioritas $PATH: " + realFirst,
			},
			Remedy: "Sesuaikan urutan direktori $PATH Anda atau hapus biner sekunder",
		}, true
	}

	return DoctorItem{}, false
}

func checkBinaryPermissions() DoctorItem {
	execPath, err := OsExecutable()
	if err != nil {
		return DoctorItem{
			Category: "Permissions",
			Title:    "Self-Update Permission",
			Status:   StatusWarn,
			Message:  "Tidak dapat membaca path eksekusi biner",
		}
	}

	realExec, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		realExec = execPath
	}

	execDir := filepath.Dir(realExec)

	// Uji write permission dengan membuat probe temporary file
	probeFile, err := os.CreateTemp(execDir, ".charites-probe-*")
	if err != nil {
		return DoctorItem{
			Category: "Permissions",
			Title:    "Self-Update Permission",
			Status:   StatusWarn,
			Message:  fmt.Sprintf("Direktori biner (%s) tidak dapat ditulis langsung", execDir),
			Details: []string{
				"Izin tulis diperlukan untuk perintah 'charites update' mandiri tanpa sudo.",
			},
			Remedy: "Jalankan 'sudo charites update' atau instal biner ke direktori user (~/.local/bin)",
		}
	}

	_ = probeFile.Close()
	_ = os.Remove(probeFile.Name())

	return DoctorItem{
		Category: "Permissions",
		Title:    "Self-Update Permission",
		Status:   StatusPass,
		Message:  fmt.Sprintf("Izin tulis tersedia pada direktori biner (%s)", execDir),
	}
}

func checkUpdateConnectivity() DoctorItem {
	client := &http.Client{Timeout: 3 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getUpdateURL(), nil)
	if err != nil {
		return DoctorItem{
			Category: "Network",
			Title:    "GitHub API Connectivity",
			Status:   StatusWarn,
			Message:  "Gagal membuat request ke GitHub API",
		}
	}
	req.Header.Set("User-Agent", "Charites-Doctor")

	resp, err := client.Do(req)
	if err != nil {
		return DoctorItem{
			Category: "Network",
			Title:    "GitHub API Connectivity",
			Status:   StatusWarn,
			Message:  "Gagal menghubungi GitHub API (koneksi offline atau timeout)",
			Remedy:   "Periksa koneksi internet Anda jika ingin menggunakan 'charites update'",
		}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return DoctorItem{
			Category: "Network",
			Title:    "GitHub API Connectivity",
			Status:   StatusWarn,
			Message:  fmt.Sprintf("GitHub API merespons status HTTP %d", resp.StatusCode),
		}
	}

	return DoctorItem{
		Category: "Network",
		Title:    "GitHub API Connectivity",
		Status:   StatusPass,
		Message:  "GitHub Release API dapat diakses dengan lancar",
	}
}

func checkWorkspaceConfig() (DoctorItem, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return DoctorItem{}, false
	}

	configPath := filepath.Join(wd, "charites.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return DoctorItem{
			Category: "Workspace",
			Title:    "Workspace Config",
			Status:   StatusPass,
			Message:  fmt.Sprintf("Ditemukan berkas konfigurasi lokal: %s", configPath),
		}, true
	}

	return DoctorItem{}, false
}
