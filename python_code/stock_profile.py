import httpx

url = 'https://stockanalysis.com/stocks/smci/company/__data.json'
response = httpx.get(url)
raw_data = response.json()

nodes = raw_data['nodes']
data = nodes[2]['data']
data_map = data[0]

description_index = data_map['description']
stock_profile_index = data_map['profile']
stock_profile_data = data[stock_profile_index]
stock_industry_data_index = stock_profile_data['industry']
stock_sector_data_index = stock_profile_data['sector']

stock_name_index = stock_profile_data['name']
stock_country_index = stock_profile_data['country']
stock_founded_index = stock_profile_data['founded']
stock_ipo_date_index = stock_profile_data['ipoDate']
stock_ceo_index = stock_profile_data['ceo']
stock_industry_index = data[stock_industry_data_index]['value']
stock_sector_index = data[stock_sector_data_index]['value']

stock_profile_doc = {
    'name': data[stock_name_index],
    'description': data[description_index],
    'country': data[stock_country_index],
    'founded': data[stock_founded_index],
    'ipoDate': data[stock_ipo_date_index],
    'industry': data[stock_industry_index],
    'sector': data[stock_sector_index],
    'ceo': data[stock_ceo_index],
}

print(stock_profile_doc)