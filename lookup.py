from playwright.sync_api import sync_playwright
from bs4 import BeautifulSoup
import re, random
from argparse import ArgumentParser

def get_url_of_user(html):
    soup = BeautifulSoup(html, 'html.parser')
    try:
        link = soup.find_all('input', {'name':'link'})
        owner = link[-1]['value']
    except IndexError:
        a_link = soup.find('a', class_='ls_contacts-btn', href=True)
        if a_link:
            owner = a_link['href']
    url = f'https://usphonebook.com{owner}'
    return url

def get_target_information(html):
    soup = BeautifulSoup(html, 'html.parser')
    name_span = soup.find('span', class_='header-name')
    full_name = name_span.get_text(strip=True)
    address_element = soup.find('p', class_='ls_contacts__text')
    current_address = address_element.get_text(strip=True)
    title_tags = soup.find_all('p', class_='ls_contacts__title')
    associates = []
    emails = []
    for title_tag in title_tags:
        if 'Current Phone Number' in title_tag.get_text(strip=True):
            person_div = title_tag.find_next('div', itemscope=True)
            if person_div:
                phone_number_span = person_div.find('span', itemprop='telephone')
                # get current phone number
                current_phone_number = phone_number_span.get_text(strip=True) if phone_number_span else None
            else:
                current_phone_number = None
        elif 'Previous Addresses:' in title_tag.get_text(strip=True):
            addresses_list = title_tag.find_next('ul')
            if addresses_list:
                addresses_list_items = addresses_list.find_all('li')
                # Extract previous addesses
                previous_addresses = [item.get_text(strip=True).rsplit(',', 1)[0] for item in addresses_list_items] if addresses_list_items else None
            else:
                previous_addresses = None
        elif 'Previous Phone Numbers:' in title_tag.get_text(strip=True):
            phone_numbers_list = title_tag.find_next('ul')
            if phone_numbers_list:
                phone_numbers_list_items = phone_numbers_list.find_all('li')
                # Extract individual phone numbers
                previous_phone_numbers = [item.find('a').get_text(strip=True) if item and item.find('a') else None for item in phone_numbers_list_items]
            else:
                previous_phone_numbers = None
        elif 'Relatives' in title_tag.get_text(strip=True):
            relatives_section = title_tag.find_next('div', class_='section-relative')
            if relatives_section:
                    relatives_items = relatives_section.find_all('span', itemprop='name')
                    # Extract individual relative names
                    relatives = [item.get_text(strip=True) for item in relatives_items] if relatives_items else None
            else:
                relatives = None
        elif 'Associates' in title_tag.get_text(strip=True):
            section_relative_div = title_tag.find_next('div', class_='section-relative')
            if section_relative_div:
                associate_links = section_relative_div.find_all('a')
                associates.extend(span.get_text(strip=True) for span in associate_links)
            else:
                associates = None
        elif 'Email' in title_tag.get_text(strip=True):
            ul_tag = title_tag.find_next('ul')
            if ul_tag:
                email_links = ul_tag.find_all('a', href=lambda x: x and 'mailto' in x)
                # Extract and append the href attribute of each <a> tag to the emails list
                emails = [link.get_text(strip=True) for link in email_links] if email_links else None
            else:
                emails = None

    associates_list = []
    if associates != None:
        for associate in associates:
            associates_list.append(associate.replace('\n',' '))

    return full_name, current_phone_number, current_address, previous_addresses, previous_phone_numbers, relatives, associates_list, emails

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
        devices = p.devices
        rand_device = random.choice(list(devices.keys()))
        device = p.devices[rand_device]
        browser = p.webkit.launch()
        context = browser.new_context(**device,)
        context = browser.new_context()
        page = context.new_page()
        page.goto(url)
        page.wait_for_timeout(random.randint(4000, 7500))
        html_content = page.content()
        url = get_url_of_user(html_content)
        context = browser.new_context()
        page = context.new_page()
        page.goto(url)
        page.wait_for_timeout(random.randint(1500, 2500))
        html = page.content()
        browser.close()
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
            full_name, current_phone_number, current_address, previous_addresses, previous_phone_numbers, relatives, associates, emails = get_target_information(html)

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
            full_name, phone_number, current_address, previous_addresses, previous_phone_numbers, relatives, associates, emails = get_target_information(html)

    if 'Current & Past Contact Info' in previous_addresses:
        previous_addresses = None

    print("[*] Results from usphonebook.com")
    print(f'Phone Number: {phone_number if phone_number is not None else "N/A"}\n'
        f'Owner: {full_name if full_name is not None else "N/A"}\n'
        f'Current Address: {current_address if current_address is not None else "N/A"}\n'
        f'Prior Addresses: {", ".join(filter(None, previous_addresses)) if previous_addresses else "N/A"}\n'
        f'Prior Phone Numbers: {", ".join(filter(None, previous_phone_numbers)) if previous_phone_numbers else "N/A"}\n'
        f'Relatives: {", ".join(filter(None, relatives)) if relatives else "N/A"}\n'
        f'Associates: {", ".join(filter(None, associates)) if associates else "N/A"}\n'
        f'Email: {", ".join(filter(None, emails)) if emails else "N/A"}\n\n')

if __name__=='__main__':
    main()
