from spapi.spapi import SPAPI



def main():
    client = SPAPI()
    response = client.get_competitive_pricing(asin_list)
    print(response.json())