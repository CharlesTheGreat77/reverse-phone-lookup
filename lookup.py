from playwright_stealth import stealth_sync
from playwright.sync_api import sync_playwright
from bs4 import BeautifulSoup
import re, sys

def scrape_data(html):
    pattern = r'gResults:\s*\'(.*?)\''
    match = re.search(pattern, html)
    data = match.group(1)
    if 'gResults: "null"' in data:
        data = False
    else:
        # tons of cleanup
        data = re.sub(r'(\w+):', r'"\1":', data)
        data = data.replace('&quot;', '')
        data = re.sub(r'(\w+)\s*:\s*true', r'"\1":"true"', data)
        data = re.sub(r'(\w+)\s*:\s*false', r'"\1":"false"', data)
        data = data.replace('}', ']').replace("'", '"').replace('{', '[').replace(':', '":"')
        data = re.sub(r'(\w+)\s*:\s*null', r'"\1":"null"', data)
        data = '{' + data + '}'
    return data

def get_info(data):
    # full name
    try:
        pattern = r'fullName":(.*?),'
        full_name = re.search(pattern, data)
        full_name = full_name.group(1).replace('"', '')
    except Exception as msg:
        print("[-] No data found..")
        exit(0)

    # address and address history
    pattern = r'fullAddress":(.*?):' 
    addresses = re.findall(pattern, data)
    addy_list = []
    for addy in addresses:
        addy_list.append(addy.replace(',fullAddressDisplay"', '').replace('"', ''))

    # relatives
    pattern = r'name":(.*?),'
    relatives = re.findall(pattern, data)
    family_names = []
    for name in relatives:
        if name in relatives:
            name = name.replace('"', '')
            family_names.append(name)
    
    return full_name, addy_list, family_names
    
     # find link for firstname_id_xxx
    # full_name_list = full_name.split()
    # split_name = full_name.split()
    # encoded_name = '-'.join(split_name).lower()
    # encoded_url = re.search(encoded_name + r'_id_(.*?),', json_str)
    # encoded_url = encoded_url.group(1)
    # target_url = f'https://usphonebook.com/{encoded_name}_id_{encoded_url}'

def main(phone_number):
    with sync_playwright() as p:
        browser = p.webkit.launch()
        context = browser.new_context()
        page = context.new_page()
        stealth_sync(page)
        page.goto(f'https://usphonebook.com/{phone_number}')
        page.wait_for_timeout(5000)
        html_content = page.content()
        browser.close()

    soup = BeautifulSoup(html_content, 'html.parser')
    html = str(soup.prettify)
    data = scrape_data(html)
    if data != False:
        full_name, addy_list, family_names = get_info(data)
        print(f'Phone Number: {phone_number}\nOwner: {full_name}\nAddress: {addy_list[0]}\nPrior Addresses: {", ".join(addy_list)}\nRelatives & Potentially Past Owners: {", ".join(family_names)}\nAdditional Info: https://usphonebook.com/{phone_number}')
    else:
        print(f'[-] Returned None as resuslts.. try this\nURL: https://usphonebook.com/{phone_number}')

if __name__=='__main__':
    phone_number = sys.argv[1]
    if len(phone_number) != 12:
        print("usage: python3 main.py 987-654-3221")
    else:
        main(phone_number)
