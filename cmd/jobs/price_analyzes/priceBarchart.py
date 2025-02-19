import matplotlib.pyplot as plt
import numpy as np
import mysql.connector

# data from https://allisonhorst.github.io/palmerpenguins/

def historyBarPlot ():
    fruit_names = ['Coffee', 'Salted Caramel', 'Pistachio']
    fruit_counts = [4000, 2000, 7000]

    fig, ax = plt.subplots()
    bar_container = ax.bar(fruit_names, fruit_counts)
    ax.set(ylabel='pints sold', title='Gelato sales by flavor', ylim=(0, 8000))
    ax.bar_label(bar_container, fmt='{:,.0f}')

    plt.show() 



def housePriceQuery(): 
    mydb = mysql.connector.connect(
    host="127.0.0.1",
    user="root",
    password="test@test",
    database="words"
    )

    mycursor = mydb.cursor()

    mycursor.execute("SELECT * FROM house_price")

    myresult = mycursor.fetchall()

    for x in myresult:
        print(x)



housePriceQuery()
historyBarPlot()