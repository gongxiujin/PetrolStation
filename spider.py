import os

import pandas
import json
import datetime


def run():
    page = pandas.read_html("https://oil.usd-cny.com/", header=0, index_col=0)
    petrol = page[1]
    petrol.rename(columns=dict(zip(petrol.columns.to_list(), [i.split("号")[0] for i in petrol.columns.to_list()])),
                  inplace=True)
    petrol.replace("-", "0", inplace=True)
    petrol.replace("", "0", inplace=True)
    petrol.astype(float)
    r = list()
    for k, v in petrol.loc['山东'].to_dict().items():
        r.append({"version": k, "price": v, "day": datetime.datetime.now().strftime("%Y-%m-%d")})
    os.system("redis-cli -n 1 del dailyPetrol")
    os.system("redis-cli -n 1 set dailyPetrol '%s'" % json.dumps(r))
    os.system("redis-cli -n 1 set DailyInsert 1")


if __name__ == '__main__':
    run()