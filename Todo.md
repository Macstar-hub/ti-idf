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
    
    - [x] Make asset manager page:
        - [x] Make home page.
        - [x] Make asset calcuation.
        - [x] Make asset DB Qurty function. 
        - [x] Make udpate function asset price from telegram.
        - [x] Make asset resault page.
        - [] Make asset price bot with telegram source.
    - [] Make API docs: 
        - [x] Make landing page.
        - [x] Make post links page.
        - [x] Make links list page.
        - [] Make authentication with keycloak.
        - [] Make web crawling with go colly.
        
    - [] Make post all docs throw kafka.
        - [] Procces TF-IDF Module throw kafka.
        - [] Make back-pressure throw kafka.
        - [] Make serilizaton and deserilize with schema registery proto buff.

    - [] Make cache layer for topic modeling and topic vectorization instead of sql DB.
        - [] Make redis SET and GET command modules.
---
# Issues : 
- [] Find duplication values inserted in database.
    - [] Find and fix TF function duplication words.
        - [] Find reason of un succesfull bag of word clean up in archive.dir.
- [x] Fix too many connection to databases error.
- [] How to reuse old tcp connction to db in new operation.