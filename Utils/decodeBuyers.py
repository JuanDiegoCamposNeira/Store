import json, sys

filePath = '/Users/juancamposneira/Desktop/Reto-Aprendizaje/Data/Encoded/Buyers.txt'
jsonFilePath = '/Users/juancamposneira/Desktop/Reto-Aprendizaje/Data/Decoded/buyers.json'

buyers = []
with open(filePath, 'r', encoding='utf-8') as file: 
    json_string = file.read()
    buyers_obj = json.loads(json_string) 
    for buyer in buyers_obj: 
        buyer_dict = { 'type': 'Person', 'uid': '_:' + buyer['id'], 'name': buyer['name'], 'age': buyer['age'] }
        buyers.append(buyer_dict) 

# Parse data into JSON format
with open(jsonFilePath, 'w', encoding='utf-8') as jsonFile: 
    json.dump(buyers, jsonFile, indent = 4, ensure_ascii = False)