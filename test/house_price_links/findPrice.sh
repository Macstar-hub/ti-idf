Divar search example: 
1-  curl 'https://divar.ir/s/tehran/buy-residential/west-tehran-pars'

2- for search and find all links: 
curl "https://divar.ir/s/tehran/buy-residential/west-tehran-pars?size=65-80" | grep -i 'https://divar.ir/v/' | tr "," "\n"  | grep -i "url" | tr -d "\"" | cut -d ":" -f2-


3- for line convert: 
#!/bin/bash

input_file="priceList.txt"  # Replace with the actual filename containing the URLs

# Read the file and format the output
formatted=$(awk '{printf "\"%s\", ", $0}' "$input_file")

# Remove the trailing comma and space
formatted=${formatted%, }

echo "$formatted"