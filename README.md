# reverse-phone-lookup
Lightweight reverse phone/name lookup tool written in Go! ☺️ This is only for US individuals only.
It goes about such by scraping usphonebook.com

# Prerequisite
| Prerequisite | Version |
|--------------|---------|
| Go           |  <=1.22 |
```
apt install golang-go || brew install go
```

# Install
```
git clone https://github.com/CharlesTheGreat77/reverse-phone-lookup
cd reverse-phone-lookup
go mod init main
go build reverse-phone-lookup.go
```

# Usage
```
./reverse-phone-lookup -h
Usage of ./reverse-phone-lookup:
  -city string
        specify the city the target resides [Los Angelos]
  -fullname string
        specify the first name to the target [bob]
  -h    show usage
  -phone string
        specify a phone number [777-999-0000]
  -state string
        specify the state the target resides [California]
```
# Note
Im just getting use to Go so bear with me but once I get comfortable, I'll attempt to make it scalable for other sources