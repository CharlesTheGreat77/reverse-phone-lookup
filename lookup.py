from playwright_stealth import stealth_sync
from playwright.sync_api import sync_playwright
from bs4 import BeautifulSoup
import re, sys
from argparse import ArgumentParser

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


def encode_args(unencoded_args):
    split_list = unencoded_args.split(' ')
    encoded_data = ''
    if (len(split_list) > 1):
        encoded_data+= split_list[0]
        for i in range(1, len(split_list)):
            encoded_data  += f'-{split_list[i]}'
    else:
        encoded_data += unencoded_args # only needs to be encoded if there's spaces
    return encoded_data.lower()

def get_content_from_usphonebook(url):
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)
        context = browser.new_context()
        page = context.new_page()
        stealth_sync(page)
        page.goto(url)
        page.wait_for_timeout(1500)
        html_content = page.content()
        browser.close()
        html = soup_html(html_content)
        return html

def soup_html(html):
    soup = BeautifulSoup(html, 'html.parser')
    html = str(soup.prettify)
    return html

def main():
    parser = ArgumentParser(description="Reverse Phone Lookup with Playwright")
    parser.add_argument('-p', '--phone_number', help='phone number formatted like 999-111-2222', type=str, required=False)
    parser.add_argument('-n', '--name', help='name of the individual [first last]', type=str, required=False)
    parser.add_argument('-c', '--city', help='enter the city the individual resides [only used with -n]', type=str, required=False)
    parser.add_argument('-s', '--state', help='enter the state the individual resides [only used with -n]', type=str, required=False)
    args = parser.parse_args()

    phone_number = args.phone_number
    name = args.name
    city = args.city
    state = args.state

    if phone_number:
        if len(phone_number) != 12:
            print(parser.usage)
            exit(0)
        else:
            url = f'https://usphonebook.com/{phone_number}'
            html = get_content_from_usphonebook(url)
            target_content = scrape_data(html)
            full_name, addy_list, family_names = get_info(target_content)

            print("[*] Results from usphonebook:")
            print(f'Phone Number: {phone_number}\nOwner: {full_name}\nAddress: {addy_list[0]}\nPrior Addresses: {", ".join(addy_list)}\nRelatives & Potentially Past Owners: {", ".join(family_names)}\nAdditional Info: https://usphonebook.com/{phone_number}\n\n')

    if name:
        if name and state:
            encoded_name = encode_args(name)
            encoded_state = encode_args(state)
            if city:
                encoded_city = encode_args(city)
                url = f'https://usphonebook.com/{encoded_name}/{encoded_state}/{encoded_city}'
            else:
                url = f'https://usphonebook.com/{encoded_name}/{encoded_state}'
            html = get_content_from_usphonebook(url)
            target_content = scrape_data(html)
            full_name, addy_list, family_names = get_info(target_content)
            print("[*] Results from usphonebook")
            print(f'Phone Number: {phone_number}\nOwner: {full_name}\nAddress: {addy_list[0]}\nPrior Addresses: {", ".join(addy_list)}\nRelatives & Potentially Past Owners: {", ".join(family_names)}\nAdditional Info: {url}\n')

if __name__=='__main__':
    main()
