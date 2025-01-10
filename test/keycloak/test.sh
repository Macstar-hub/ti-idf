#curl 'http://192.168.1.100:8080/realms/admin-area/protocol/openid-connect/token'  \
#-H 'Content-Type: application/x-www-form-urlencoded' \
#--data-urlencode 'client_id=admin-booking' \
#--data-urlencode 'grant_type=password' \
#--data-urlencode 'username=macstarboy' \
#--data-urlencode 'password=testtest'

export HOST=192.168.1.100:8080
export USERNAME=macstarboy
export PASSWORD=hello
export CLIENTID=admin-booking
export REALM=admin-area
export CLIENTSECRET=LOmb3mUS4YiRWlT6s9be7Uk3NHdSGBYU

curl -X POST \
    http://$HOST/realms/$REALM/protocol/openid-connect/token \
    -H 'Content-Type: application/x-www-form-urlencoded' \
    -d username=$USERNAME \
    -d password=$PASSWORD \
    -d grant_type=password \
    -d client_id=$CLIENTID \
    -d client_secret=$CLIENTSECRET
