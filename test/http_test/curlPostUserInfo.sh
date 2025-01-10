postUserInfo () { 
    while true 
    do 
    	curl -XPOST  -d "link=mamad" -d "name=mammadi" -d "label1=mac@mac" -d "label2=test" -d "label3=test" http://localhost/api/v1/postLinks && sleep 3 && clear 
    done
}

postUserInfo
