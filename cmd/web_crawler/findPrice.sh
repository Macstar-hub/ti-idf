#!/bin/bash

# export text='value":"۶٬۹۵۰٬۰۰۰٬۰۰۰ تومان"'

makeCleanPrice () { 
   echo $1 
   echo $1 |  tr ',' '\n' | grep -i "value.*" | grep -i 'تومان' | tail -n 2 | tr -d \" | tr -d value | tr -d \: | tr -d 'تومان'
}

makeCleanPrice