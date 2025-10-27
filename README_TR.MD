# Golang için Anıt Sayaç Veri Çekici

**🌐 Dil Seçenekleri:** [English](README.MD) | **Türkçe**

Veri işleme için kullanacak kişilerin JSON dosyasını kullanması önerilir. JSON dosyası her gün gece yarısında (UTC 00:00) güncellenir. Veriler [Anıt Sayaç](https://anitsayac.com/) sitesinden alınmaktadır.

## ❗ Önemli Bilgiler

- ID'leri 70, 69, 114, 78, 71, 511, 605, 607, 3580, 1051, 934, 675, 617, 616, 611, 599, 138, 4185 olan verilerde yanlış veri girişi nedeniyle hatalar olduğu bilinmektedir.
- Boş veriler sitede mevcut olmadığı için o şekilde görünmektedir.

## Kurulum

```bash
go mod tidy
go build -o crawler crawler.go
```

## Kullanım

```bash
go run .
```

## Veri Araştırması için Git Kurulumu ve Kullanımı

### Git Kurulumu

#### Windows

1. [https://git-scm.com/download/win](https://git-scm.com/download/win) adresinden Git'i indirin
2. Kurulum dosyasını çalıştırın ve kurulum sihirbazını takip edin
3. Kurulumu doğrulayın: `git --version`

#### macOS

```bash
# Homebrew kullanarak
brew install git

# Veya https://git-scm.com/download/mac adresinden indirin
```

#### Linux (Ubuntu/Debian)

```bash
sudo apt update
sudo apt install git
```

### Veri Araştırması için Git Komutları

#### Commit geçmişinde belirli isimleri arama

```bash
# Tüm branch'lerde "Fidan Çakır" için arama
git log --oneline -S "Fidan Çakır" --source --all

# Commit geçmişinde "İkbal Uzuner" arama
git log --oneline -S "İkbal Uzuner" --source --all

# Commit mesajlarında grep ile arama
git log --all --grep="Fidan Çakır" --oneline
```

#### Belirli commit'ler için detaylı değişiklikleri görüntüleme

```bash
# Belirli bir commit'teki değişen dosyaları göster
git show --name-only COMMIT_HASH

# Bir commit için tam farkı göster
git show COMMIT_HASH

# Patch detayları ile değişiklikleri göster
git log -S "ARAMA_TERIMİ" --patch --oneline
```

**Not:** `COMMIT_HASH` yerine gerçek commit hash'ini yazın (örn: `38b5d7f` veya `c3e1dfb`). Commit hash'lerini `git log --oneline` komutu ile görebilirsiniz.

#### Dosya içeriklerinde arama

```bash
# Mevcut dosyalarda bir terim ara
git grep "Fidan Çakır"

# Belirli dosya türlerinde ara
git grep "arama_terimi" -- "*.json"
git grep "arama_terimi" -- "*.csv"
```

#### Repository geçmişini tarama

```bash
# Tüm commit'leri tek satırda görüntüle
git log --oneline

# Dosya değişiklikleriyle commit'leri görüntüle
git log --stat

# Belirli bir zaman aralığındaki commit'leri görüntüle
git log --since="2024-01-01" --until="2024-12-31"
```

#### Bu repository için örnekler

```bash
# Belirli mağdurları içeren commit'leri bul
git log --oneline -S "Fidan Çakır" --source --all
git log --oneline -S "İkbal Uzuner" --source --all

# Belirli konumları ara
git log --oneline -S "İzmir" --source --all
git log --oneline -S "İstanbul" --source --all

# Tarih desenlerini bul
git log --oneline -S "16/10/2024" --source --all

# Sadece veri dosyalarında ara
git log --oneline -- data.json data.csv
```

## Veri Kaynağı

- [Anıt Sayaç](https://anitsayac.com/)

## Lisans

[MIT](LISENCE)

## Yazar

[Ramazan Sancar](https://www.github.com/ramazansancar)
