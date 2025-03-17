# This Todo list for solving problem 
* Next Steps: 
    - [] Develop auto labeling based on TF-IDF
        - [x] develop TF function.
            - [x] make pre procceing function for cleaning stopwords and make lowercase.
                - [x] make pre check stopwords baaed on regexp lib.
        - [x] develop idf function.  
            - [x] Make new lines from all strings with sepecifc delimator.
            - [x] Fix IDF word count per document.  
                - [x] Fix document count has word and count them to make idf. 
        - [x] Make Multiple TF And IDF.
            - [x] Make database connection with mysql sample.
                - [x] Make Insert And Duplication values.
                
        - [] Make use wevaite or qudarant Vector DB instead of mysql.
        - [] Make Topic Modeling.
    
    - [x] Make links manager page:
        - [x] Make home page.
        - [x] Make asset calcuation.
        - [x] Make asset DB Qurty function. 
        - [x] Make udpate function asset price from telegram.
        - [x] Make asset resault page.
        - [] Make daily note section.
        - [] Add advanced search mode.

    - [] Make API docs: 
        - [x] Make landing page.
        - [x] Make post links page.
        - [x] Make links list page.
        - [] Make search function.
            - [] Make switch for label and name column.
            - [] Make regex find in search function.
            - [x] Make MVP search function in mysql internal. 
        - [] Make authentication with keycloak.
        - [] Make Edit function in ui for recorde ui. 
        - [] Make delete function.
        - [] Make post all docs throw kafka.
            - [] Procces TF-IDF Module throw kafka.
            - [] Make back-pressure throw kafka.
            - [] Make serilizaton and deserilize with schema registery proto buff.

    - Make Price Crawling.
        - [x] Make web crawling with go.
        - [x] Make crawling with tjgu.
        - [x] Make house price from divar.
        - [x] Make outliner detection and cleaner to make correct price.
        - [] Make function to detect specefic house area.
        - [x] Make function to produce avg house price. 
        - [] Make more descriptive error handling while expercing error in divar house price crawling.
        - [] Make export price house, dollar, gold and coin to ui.
        - [] Fix duplicated data when house_price.
        - [] Make telegram bot for notify custome reports.
            - [x] Make send telegram best price house price.
            - [] Make section based house price.


    - [] Make cache layer for topic modeling and topic vectorization instead of sql DB.
        - [] Make redis SET and GET command modules.

    - [] Security:
        - [] Make read all token and security config from database.
        - [] Make all security cred with salt algo.
---
# Issues : 
- [] Find duplication values inserted in database.
    - [] Find and fix TF function duplication words.
        - [] Find reason of un succesfull bag of word clean up in archive.dir.
- [x] Fix too many connection to databases error.
- [] How to reuse old tcp connction to db in new operation.