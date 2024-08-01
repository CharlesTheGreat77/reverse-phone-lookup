# reverse-phone-lookup
Lightweight reverse phone/name lookup tool written in Go! ☺️ This is only for US individuals only.
It goes about such by scraping usphonebook.com

![osint](https://github.com/CharlesTheGreat77/reverse-phone-lookup/assets/27988707/03facb8c-4e1d-480c-92cf-d3297408d03d)

# Prerequisite 🚀
| Prerequisite | Version |
|--------------|---------|
| Go           |  <=1.22 |
```
apt install golang-go || brew install go
```

# Install 💻
```
git clone https://github.com/CharlesTheGreat77/reverse-phone-lookup
cd reverse-phone-lookup
go mod init reverse-phone-lookup
go mod tidy
go build main.go
```

# Usage 🎯
```
./main -h
Usage of ./reverse-phone-lookup:
  -city string
        specify the city the target resides [Los Angelos]
  -fullname string
        specify the targets full name [John Doe]
  -h    show usage
  -phone string
        specify a phone number [777-999-0000]
  -state string
        specify the state the target resides [California]
```

# Example
```
./main -state California -city "Los Angelos" -fullname "John Doe"
```

# Discord Bot 🔨
The reverse-phone-lookup is now available as a discord bot. To setup the discord bot follow these steps:
1. Enter discord bot token in the **config.json** file.
2. Download the discordgo package:
```
go get github.com/bwmarrin/discordgo
```
3. Build the discord bot executable:
```
go build usphonebook_bot.go
```
4. Run the executable!
```
./usphonebook_bot.go
```
## Discord bot commands
```
Usage: <@bot_id> lookup phone=<number> fullname=<name> state=<state> city=<city>
```
example:
```
@bot lookup state=road island city=providence fullname=john doe
```
This example searches usphonebook for **john doe** in the city and state specified.
```
@bot lookup phone=777-111-2222
```


# Note
Im just getting use to Go so bear with me but once I get comfortable, I'll attempt to make it scalable for other sources
