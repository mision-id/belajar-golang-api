package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type Kategori struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Goreng", Harga: 3000, Stok: 10},
	{ID: 2, Nama: "Teh Botol Sosro", Harga: 1500, Stok: 20},
	{ID: 3, Nama: "Kacang Atom Garuda", Harga: 2000, Stok: 15},
}

var kategori = []Kategori{
	{ID: 1, Name: "Makanan", Description: "Makanan berat"},
	{ID: 2, Name: "Minuman", Description: "Minuman"},
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	//Gunakan trim prefix untuk mengambil request id dengan menghilangkan prefix di depan request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/produk/")
	//Ubah request id jadi integer
	id, err := strconv.Atoi(idStr)
	//Tampilkan HTTP Header Bad request jika id tidak ditemukan
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	//Cari id produk dan tampilkan
	for _, p := range produk {
		if p.ID == id {
			//Set jadi konsensus JSON
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
}

func tambahProduk(w http.ResponseWriter, r *http.Request) {
	// Buat temporary variable untuk menyimpan data dari request ke variable ProdukBaru
	var ProdukBaru Produk
	err := json.NewDecoder(r.Body).Decode(&ProdukBaru)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	//Masukan data
	//Define ID ProdukBaru dengan mengambil list ID produk ditambah 1
	ProdukBaru.ID = len(produk) + 1
	//Masukan data Produkbaru ke produk
	produk = append(produk, ProdukBaru)
	//Buat Header status Created 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Produk baru berhasil ditambahkan",
	})
	json.NewEncoder(w).Encode(ProdukBaru)
}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	//Gunakan trim prefix untuk mengambil request id dengan menghilangkan prefix di depan request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/produk/")
	//Ubah request id jadi integer
	id, err := strconv.Atoi(idStr)
	//Tampilkan HTTP Header Bad request jika id tidak ditemukan
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	//Ganti data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Looping produk, cari ID, ganti sesuai request
	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}
	http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	//Gunakan trim prefix untuk mengambil request id dengan menghilangkan prefix di depan request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/produk/")
	//Ubah request id jadi integer
	id, err := strconv.Atoi(idStr)
	//Tampilkan HTTP Header Bad request jika id tidak ditemukan
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}
	//Loop Produk cari ID, dan mendapatkan index yang mau di hapus
	for i, p := range produk {
		if p.ID == id {
			//Buat slice baru dengan data sebelum dan sesudah index
			produk = append(produk[:i], produk[i+1:]...)
			//Tampilkan data produk setelah dihapus
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Data berhasil di hapus",
			})
			return
		}
	}
	http.Error(w, "Data tidak ditemukan", http.StatusNotFound)
}

// Function Model Kategori
// Fungsi get kategori by ID
func getKategoriByID(w http.ResponseWriter, r *http.Request) {
	//Gunakan trim prefix untuk mengambil request id dengan menghilangkan prefix di depan request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/kategori/")
	//Ubah id jadi integer
	id, err := strconv.Atoi(idStr)
	//Tampilkan HTTP Bad request jika tidak ditemukan
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	//Loop kategori ID dan tampilkan
	for _, k := range kategori {
		if k.ID == id {
			//Set jadi konsensus JSON
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(k)
			return
		}
	}
	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
}

// Fungsi menambahkan kategori
func tambahKategori(w http.ResponseWriter, r *http.Request) {
	// Siapkan temporary variabel kategoriBaru untuk menyimpan data request ke kategori
	var KategoriBaru Kategori
	err := json.NewDecoder(r.Body).Decode(&KategoriBaru)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	//Tambahkan data
	// Ambil list ID kemudian setelah ketemu tambah ID terkahir dengan 1
	KategoriBaru.ID = len(kategori) + 1
	// Masukan data KategoriBaru ke kategori
	kategori = append(kategori, KategoriBaru)
	//Buat Header status created
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Kategori berhasil ditambahkan",
	})
	json.NewEncoder(w).Encode(KategoriBaru)
}

// Fungsi Mengedit kategori berdasarkan ID kategori
func updateKategoriByID(w http.ResponseWriter, r *http.Request) {
	//Cari ID kategori berdasarkan request yang berupa string
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/kategori/")
	//ubah request id jadi integer
	id, err := strconv.Atoi(idStr)
	//Tampilkan status Invalid Request jika request tidak sesuai
	if err != nil {
		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
		return
	}

	//Buat temporary variabel UpdateKategori untuk menyimpan data request
	var UpdateKategori Kategori
	err = json.NewDecoder(r.Body).Decode(&UpdateKategori)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
	}
	//Masukan data UpdateKategori ke katgeori
	for i := range kategori {
		if kategori[i].ID == id {
			UpdateKategori.ID = id
			kategori[i] = UpdateKategori
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kategori)
			return
		}
	}
	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
}

// Fungsi Hapus Kategori
func deleteKategori(w http.ResponseWriter, r *http.Request) {
	//Cari ID kategori dengan trim prefix berdasarkan request string
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/kategori/")
	//Ubah string jadi integer
	id, err := strconv.Atoi(idStr)
	//Tampilkan header bad request jika id tidak ditemukan
	if err != nil {
		http.Error(w, "Invalid Katgori ID", http.StatusBadRequest)
	}
	//Loop cari kategori ID dan hapus data berdasarkan ID
	for i, k := range kategori {
		if k.ID == id {
			//Buat slice baru dengan data sebelum dan sesudah index
			kategori = append(kategori[:i], kategori[i+1:]...)
			//Tampilkan data produk setelah dihapus
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Kategori berhasil di hapus",
			})
			return
		}
	}
	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
}

func main() {
	fmt.Println("Server Running on localhost:8080")

	//Endpoint route /api/v1/produk
	http.HandleFunc("/api/v1/produk", func(w http.ResponseWriter, r *http.Request) {
		//Set jadi konsensus JSON
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			// Tampilkan data produk
			json.NewEncoder(w).Encode(produk)
		}

	})

	//Endpoint route /api/v1/kategori
	http.HandleFunc("/api/v1/kategori", func(w http.ResponseWriter, r *http.Request) {
		//Set jadi konsensus JSON
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			// Tampilkan data produk
			json.NewEncoder(w).Encode(kategori)
		}

	})

	//Management Endpoint route /api/v1/produk/
	http.HandleFunc("/api/v1/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			//Get produk by ID
			getProductByID(w, r)
		} else if r.Method == "POST" {
			//Tambah Produk
			tambahProduk(w, r)
		} else if r.Method == "PUT" {
			//Update Produk
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			//Delete Produk by ID
			deleteProduk(w, r)
		}
	})

	//Management Endpoint route /api/v1/kategori/
	http.HandleFunc("/api/v1/kategori/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			//Get kategori by ID
			getKategoriByID(w, r)
		} else if r.Method == "POST" {
			//Tambah Kategori
			tambahKategori(w, r)
		} else if r.Method == "PUT" {
			//edit kategori by ID
			updateKategoriByID(w, r)
		} else if r.Method == "DELETE" {
			//Delete kategori by ID
			deleteKategori(w, r)
		}
	})

	//Buat route default untuk keterangan server running
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//set jadi konsensus JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "Ok",
		})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		//set jadi konsensus JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "Ok",
			"message": "API Running",
		})
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to Start")
	}

}
