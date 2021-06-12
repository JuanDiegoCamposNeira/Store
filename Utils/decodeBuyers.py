import json, sys

if len(sys.argv) != 3: 
    print('Not enought arguments: \npython3 decodeBuyers.py <input_file_path> <output_json_path_file>')
    exit()

print('Buyers are already given in JSON format :)')