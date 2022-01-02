import openpyxl

from models import Product

PATH = './insert.xlsx'

def insert():
    workbook = openpyxl.open(PATH)
    sheet = workbook['insert']
    sheet.delete_rows(1)

    for row in sheet.values:
        date, name, asin, jan, sku, fnsku, danger, sell, cost = row
        Product.create(date, name, asin, jan, sku, fnsku, danger, sell, cost)

    workbook.close()
