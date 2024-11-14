package algo

import (
	"context"
	"errors"
	"fmt"
	"github.com/ramadhia/mnc-test/internal/config"
	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/usecase"
	"math"
	"strconv"
	"strings"
	"time"
)

type AlgoIml struct {
	config config.Config
}

func NewAlgo(p *provider.Provider) *AlgoIml {
	return &AlgoIml{
		config: p.Config(),
	}
}

func (t AlgoIml) AlgoTest1(ctx context.Context, args []string) (string, error) {

	if len(args) < 2 {
		err := errors.New("argument tidak sesuai")
		return "false", err
	}

	totalInput, err := strconv.Atoi(args[0])
	if err != nil {
		err := errors.New("argument tidak sesuai")
		return "false", err
	}

	inputs := args[1:]
	if totalInput != len(inputs) {
		err := errors.New("total argument tidak sesuai")
		return "false", err
	}

	// mencari firstword pertama yang match
	var firstWord string
	wordIndex := make(map[string]int)
	minDistance := len(inputs)
	for i, v := range inputs {
		if prevIndex, ok := wordIndex[v]; ok {
			distance := i - prevIndex
			if distance < minDistance {
				minDistance = distance
				firstWord = v
			}
		} else {
			wordIndex[v] = i
		}
	}

	var matchedWords []int
	for i, v := range inputs {
		if t.isMatchString(v, firstWord) {
			matchedWords = append(matchedWords, i+1)
		}
	}

	if len(matchedWords) > 0 {
		return fmt.Sprintf("%d", matchedWords), nil
	}

	return "false", nil
}

func (t AlgoIml) AlgoTest2(ctx context.Context, totalBelanja float32, uangDibayar float32) (string, error) {
	if totalBelanja == 0 || uangDibayar == 0 {
		return "", errors.New("nilai tidak boleh kosong")
	}

	if uangDibayar < totalBelanja {
		return "", errors.New("uang yang dibayarkan kurang dari total belanja")
	}

	// Hitung kembalian
	kembalian := uangDibayar - totalBelanja
	kembalianBulat := math.Floor(float64(kembalian)/100) * 100
	pecahan := []float32{100000, 50000, 20000, 10000, 5000, 2000, 1000, 500, 200, 100}

	pecahanUang := make(map[float32]float32)
	modKembalian := kembalianBulat
	for _, p := range pecahan {
		count := math.Floor(float64(modKembalian) / float64(p))

		if count > 0 {
			pecahanUang[p] = float32(count)
			modKembalian = math.Mod(modKembalian, float64(p))
		}
	}

	var pecahanOutput string
	for _, p := range pecahan {
		if count, found := pecahanUang[p]; found {
			if p >= 1000 {
				pecahanOutput += fmt.Sprintf("%d lembar %.0f\n", int(count), p)
			} else {
				pecahanOutput += fmt.Sprintf("%d koin %.0f\n", int(count), p)
			}
		}
	}

	return fmt.Sprintf("Kembalian yang harus diberikan kasir: %.0f, dibulatkan menjadi %.0f\nPecahan uang:\n%s", kembalian, kembalianBulat, pecahanOutput), nil
}

func (t AlgoIml) AlgoTest3(ctx context.Context, arg string) (bool, error) {
	runes := map[rune]rune{
		'<': '>',
		'{': '}',
		'[': ']',
	}

	stack := []rune{}
	for _, char := range arg {

		// Jika karakter adalah pembuka, masukkan ke stack
		if char == '<' || char == '{' || char == '[' {
			stack = append(stack, char)
		} else if char == '>' || char == '}' || char == ']' {
			// Jika stack kosong, berarti ada penutup yang tidak cocok
			if len(stack) == 0 {
				return false, errors.New("empty stack")
			}

			// Ambil karakter paling atas dari stack dan periksa apakah cocok
			top := stack[len(stack)-1]
			if runes[top] == char {
				stack = stack[:len(stack)-1]
			} else {
				return false, errors.New("invalid character pasangan")
			}
		}
	}

	return len(stack) == 0, nil
}

func (t AlgoIml) AlgoTest4(ctx context.Context, req usecase.AlgoTest4Request) (bool, string, error) {
	totalCutiKantor := 14
	minimumHariKerja := 180

	// Parsing string input tanggalJoin dan tanggalCuti menjadi tipe time.Time
	joinDate, err := time.Parse("2006-01-02", req.JoinDate)
	if err != nil {
		return false, "", errors.New("format tanggal join tidak valid")
	}

	cutiDate, err := time.Parse("2006-01-02", req.CutiDate)
	if err != nil {
		return false, "", errors.New("format tanggal cuti tidak valid")
	}

	// Hitung tanggal pertama kali karyawan boleh cuti setelah 180 hari dari tanggal join
	startCanDate := joinDate.AddDate(0, 0, minimumHariKerja)

	// Jika tanggal rencana cuti lebih awal dari tanggal boleh cuti, karyawan belum bisa mengambil cuti
	if cutiDate.Before(startCanDate) {
		return false, "Belum memenuhi masa kerja 180 hari untuk cuti pribadi", nil
	}

	// Hitung jumlah hari dari tanggal boleh cuti hingga akhir tahun
	endOfYear := time.Date(cutiDate.Year(), 12, 31, 0, 0, 0, 0, time.UTC)
	availableDays := endOfYear.Sub(startCanDate).Hours() / 24
	if availableDays < 0 {
		return false, "", errors.New("tidak ada hari yang tersisa untuk cuti pribadi")
	}

	// Hitung cuti pribadi yang bisa diambil di tahun tersebut (pembulatan ke bawah)
	totalCutiPribadi := totalCutiKantor - req.CutiBersama
	quotaCuti := int(math.Floor(availableDays / 365 * float64(totalCutiPribadi)))

	// Kuota cuti pribadi untuk anak baru
	if req.CutiDurasi > quotaCuti {
		return false, fmt.Sprintf("karena hanya boleh mengambil %d hari cuti", quotaCuti), nil
	}

	// Jika kuota cuti pribadi sudah habis
	if availableDays <= 0 {
		return false, "", errors.New("kuota cuti pribadi untuk tahun pertama habis")
	}

	// Jika durasi cuti lebih dari 3 hari berturut-turut
	if req.CutiDurasi > 3 {
		return false, "", errors.New("durasi cuti pribadi tidak boleh lebih dari 3 hari berturut-turut")
	}

	// Jika semua syarat terpenuhi
	return true, "Cuti pribadi disetujui", nil
}

func (t AlgoIml) isMatchString(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if strings.ToLower(string(s1[i])) != strings.ToLower(string(s2[i])) {
			return false
		}
	}
	return true
}
