Reverse Phone Lookup OSINT with python3 and playwright for web scraping
```
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣄⠀⠀⠀⠀⠀⠀⣠⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠚⠻⠿⡇⠀⠀⠀⠀⢸⠿⠟⠓⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣠⣴⣾⣿⣶⣦⡀⢀⣤⣤⡀⢀⣴⣶⣿⣷⣦⣄⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣴⣿⣿⣿⣿⣿⣿⡇⢸⣿⣿⡇⢸⣿⣿⣿⣿⣿⣿⣦⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠘⠋⣉⡉⠙⠛⢿⣿⡇⢸⣿⣿⡇⢸⣿⡿⠛⠋⢉⣉⠙⠃⠀⠀⠀⠀
⠀⠀⢀⣤⣾⡛⠛⠛⠻⢷⣤⡙⠃⢸⣿⣿⡇⠘⢋⣤⣾⡟⠛⠛⠛⠷⣤⡀⠀⠀
⠀⢀⣾⣿⣿⡇⠀⠀⠀⠀⠙⣷⠀⠘⠛⠛⠃⠀⣾⣿⣿⣿⠀⠀⠀⠀⠈⢷⡀⠀
⠀⢸⡇⠈⠉⠀⠀⠀⠀⠀⠀⢸⡆⢰⣿⣿⡆⢰⡇⠀⠉⠁⠀⠀⠀⠀⠀⢸⡇⠀
⠀⠸⣧⠀⠀⠀⠀⠀⠀⠀⢀⡾⠀⠀⠉⠉⠀⠀⢷⡀⠀⠀⠀⠀⠀⠀⠀⣼⠇⠀
⠀⠀⠙⢷⣄⣀⠀⠀⣀⣤⡾⠁⠀⠀⠀⠀⠀⠀⠈⢷⣤⣀⠀⠀⣀⣠⡾⠋⠀⠀
⠀⠀⠀⠀⠉⠛⠛⠛⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠛⠛⠛⠉⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
```
# Install
```
pip3 install -r requirements.txt
playwright install
```

# Optional (venv)
```
python3 -m venv <venv_name>
source <venv_name>/bin/activate
pip3 install -r requirements.txt
playwright install
```

# Usage
```
usage: python3 lookup.py -n 'John Doe' -s california -c 'Santa Monica'
       python3 lookup.py -p 777-111-2222
       python3 lookup.py -n 'John Woe' -s Hawaii
```

# Content
- Owner
- Current Address
- Past Addresses
- Relatives & Associates
- Carrier

# Playwright
https://playwright.dev/python/docs/intro

# Note
Line 49: p.chromium.launch()
         - Can be changed to webkit or firefox instead of chromium.
