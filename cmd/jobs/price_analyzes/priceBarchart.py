import matplotlib.pyplot as plt
import numpy as np
import mysql.connector

# data from https://allisonhorst.github.io/palmerpenguins/

mydb = mysql.connector.connect(
host="127.0.0.1",
user="root",
password="test@test",
database="words"
)

tableName = "house_price"

def historyBarPlot ():
    prices, z_score = housePriceQuery(tableName)
    print (prices, z_score)

    fig, ax = plt.subplots()
    bar_container = ax.bar(z_score, prices)
    ax.set(xlabel='Z score Percent', ylabel='Per Squar Price', title='House Price Distribution')
    ax.bar_label(bar_container, fmt='{:,.0f}')

    plt.savefig('./test.png')
    # plt.show() 



def housePriceQuery(tableName): 
    test = [1, 4, 10]
    price = []
    z_score_percent = []
    mycursor = mydb.cursor()

    mycursor.execute("SELECT per_squar, z_score FROM " + tableName)

    myresult = mycursor.fetchall()

    for x in myresult:
        price.append(x[0])
        z_score_percent.append(x[1])
        

    return price, z_score_percent

historyBarPlot()