#!/usr/bin/env python
"""
Very simple HTTP server in python.
Usage::
    ./dummy-web-server.py [<port>]
Send a GET request::
    curl http://localhost
Send a HEAD request::
    curl -I http://localhost
Send a POST request::
    curl -d "foo=bar&bin=baz" http://localhost
"""
from http.server import BaseHTTPRequestHandler, HTTPServer
import socketserver
import pickle
import urllib.request
import json
from pprint import pprint
from pandas.io.json import json_normalize
import pandas as pd
from sklearn import preprocessing
from sklearn.preprocessing import PolynomialFeatures
from sklearn import datasets, linear_model
from sklearn.metrics import mean_squared_error, r2_score
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import scale
from sklearn.preprocessing import PolynomialFeatures
import numpy as np
from sklearn.preprocessing import StandardScaler
from sklearn.pipeline import Pipeline
from sklearn.metrics import mean_squared_error
from math import sqrt
import os
import errno
from pymongo import MongoClient
import urllib.parse as urlparse

class S(BaseHTTPRequestHandler):

    def do_training(self, contents_json, csp):
        finaldfInstances_timestamp = pd.DataFrame()
        for x in contents_json:
            for y in x["experimentloop"]:
                for z in y["experiments"]:
                    df_temp = json_normalize(z["instances"])
                    df_temp["CoreCount"] = x["corecount"]
                    df_temp["InstanceType"] = x["instancetype"]
                    df_temp["NumInstances"] = z["numinstances"]
                    if(x["instancetype"] == 't2.nano'):
                        df_temp["Mem_gib"] = 0.5
                        df_temp["cpu_credits_per_hour"] = 3

                    elif(x["instancetype"] == 't2.micro'):
                        df_temp["Mem_gib"] = 1
                        df_temp["cpu_credits_per_hour"] = 6

                    elif (x["instancetype"] == 't2.small'):
                        df_temp["Mem_gib"] = 1
                        df_temp["cpu_credits_per_hour"] = 12

                    elif (x["instancetype"] == 't2.medium'):
                        df_temp["Mem_gib"] = 4
                        df_temp["cpu_credits_per_hour"] = 24


                    elif(x["instancetype"] == 't2.large'):
                        df_temp["Mem_gib"] = 8
                        df_temp["cpu_credits_per_hour"] = 36
                    df_temp["Region"] = x["region"]
                    finaldfInstances_timestamp = finaldfInstances_timestamp.append(df_temp, ignore_index=True)
        finaldfInstances_timestamp["BootTime"] = finaldfInstances_timestamp["running"] - finaldfInstances_timestamp["pending"]
        finaldfInstances_timestamp["ShuttingDownTime"] = finaldfInstances_timestamp["terminated"] - finaldfInstances_timestamp["shuttingdown"]
        finaldfInstances_timestamp = finaldfInstances_timestamp.drop(['sshlogin'], axis=1)
        finaldfInstances_timestamp = finaldfInstances_timestamp.drop(['other'], axis=1)
        finaldfInstances_timestamp = finaldfInstances_timestamp.drop(['stopped'], axis=1)
        finaldfInstances_timestamp = finaldfInstances_timestamp.drop(['pending'], axis=1)
        finaldfInstances_timestamp = finaldfInstances_timestamp.drop(['running'], axis=1)
        finaldfInstances_timestamp = finaldfInstances_timestamp.drop(['shuttingdown'], axis=1)
        finaldfInstances_timestamp = finaldfInstances_timestamp.drop(['terminated'], axis=1)
        finaldfInstances_timestamp = finaldfInstances_timestamp.drop(['Region'], axis=1)
        finaldfInstances_timestamp = finaldfInstances_timestamp.drop(['cpu_credits_per_hour'], axis=1)
        finaldfInstances_timestamp = finaldfInstances_timestamp.sort_values(by=['CoreCount', 'Mem_gib'])
        finaldfInstances_timestamp = finaldfInstances_timestamp.reset_index(drop=True)

        t2nano = finaldfInstances_timestamp.loc[finaldfInstances_timestamp['InstanceType'] == "t2.nano"]
        t2micro = finaldfInstances_timestamp.loc[finaldfInstances_timestamp['InstanceType'] == "t2.micro"]
        t2small = finaldfInstances_timestamp.loc[finaldfInstances_timestamp['InstanceType'] == "t2.small"]
        t2medium = finaldfInstances_timestamp.loc[finaldfInstances_timestamp['InstanceType'] == "t2.medium"]
        t2large = finaldfInstances_timestamp.loc[finaldfInstances_timestamp['InstanceType'] == "t2.large"]
        t2nano = t2nano.reset_index(drop=True)
        t2micro = t2micro.reset_index(drop=True)
        t2small = t2small.reset_index(drop=True)
        t2medium = t2medium.reset_index(drop=True)
        t2large = t2large.reset_index(drop=True)

        self.do_training_save("t2.nano", csp,t2nano)
        self.do_training_save("t2.micro", csp,t2micro)
        self.do_training_save("t2.small", csp,t2small)
        self.do_training_save("t2.medium", csp,t2medium)
        self.do_training_save("t2.large", csp,t2large)

    def do_training_save(self, instancetype, csp,  instancedata):
        instancedata = instancedata.drop(['InstanceType'], axis=1)
        df_boot = instancedata.copy()
        df_boot = df_boot.drop(['ShuttingDownTime'], axis=1)
        df_boot = df_boot.sort_values(by=['NumInstances'])
        df_boot_X = df_boot.drop(['BootTime'], axis=1)
        df_boot_X = df_boot_X.values
        df_boot_y = df_boot['BootTime'].values

        df_shut = instancedata.copy()
        df_shut = df_shut.drop(['BootTime'], axis=1)
        df_shut = df_shut.sort_values(by=['NumInstances'])
        df_shut_X = df_shut.drop(['ShuttingDownTime'], axis=1)
        df_shut_X = df_shut_X.values
        df_shut_y = df_shut['ShuttingDownTime'].values


        for style, width, degree in (("g-", 2, 8), ("b--", 2, 2), ("r-", 2, 1)):
            polybig_features = PolynomialFeatures(degree=degree, include_bias=False)
            std_scaler = StandardScaler()
            lin_reg = linear_model.LinearRegression()
            polynomial_regression = Pipeline([
                ("poly_features", polybig_features),
                ("std_scaler", std_scaler),
                ("lin_reg", lin_reg),
            ])
            polynomial_regression.fit(df_boot_X, df_boot_y)
            y_newbig = polynomial_regression.predict(df_boot_X)
            rms = sqrt(mean_squared_error(df_boot_y, y_newbig))
            print("Degree= ", degree," RMS = ", rms)
            # save the model to disk
            filename = "boot/"+csp+"/"+instancetype+".sav"
            if not os.path.exists(os.path.dirname(filename)):
                try:
                    os.makedirs(os.path.dirname(filename))
                except OSError as exc: # Guard against race condition
                    if exc.errno != errno.EEXIST:
                        raise
            pickle.dump(polynomial_regression, open(filename, 'wb'))

        for style, width, degree in (("g-", 2, 8), ("b--", 2, 2), ("r-", 2, 1)):
            polybig_features = PolynomialFeatures(degree=degree, include_bias=False)
            std_scaler = StandardScaler()
            lin_reg = linear_model.LinearRegression()
            polynomial_regression = Pipeline([
                ("poly_features", polybig_features),
                ("std_scaler", std_scaler),
                ("lin_reg", lin_reg),
            ])
            polynomial_regression.fit(df_shut_X, df_shut_y)
            y_newbig = polynomial_regression.predict(df_shut_X)
            rms = sqrt(mean_squared_error(df_shut_y, y_newbig))
            print("Degree= ", degree," RMS = ", rms)
            # save the model to disk
            filename = "shut/"+csp+"/"+instancetype+".sav"
            if not os.path.exists(os.path.dirname(filename)):
                try:
                    os.makedirs(os.path.dirname(filename))
                except OSError as exc: # Guard against race condition
                    if exc.errno != errno.EEXIST:
                        raise
            pickle.dump(polynomial_regression, open(filename, 'wb'))

    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

    def do_GET(self):
        self._set_headers()
        output=''
        if self.path == '/train':
            client = MongoClient(mongo_host, 27017, username=mongo_username, password=mongo_password)
            db = client['VM_BOOT_SHUTDOWN_RATE_DB']
            collection = db['data']
            contents = collection.find()
            #contents_json = contents.json()
            self.wfile.write("started".encode())
            self.do_training(contents, "aws")
        elif '/getPrediction' in self.path:
            parsed = urlparse.urlparse(self.path)

            cloudProvider = urlparse.parse_qs(parsed.query)['csp'][0]
            instanceType = urlparse.parse_qs(parsed.query)['instanceType'][0]
            numInstances = urlparse.parse_qs(parsed.query)['numInstances'][0]
            ratetype = urlparse.parse_qs(parsed.query)['type'][0]

            print(urlparse.parse_qs(parsed.query)['numInstances'])
            print(cloudProvider,instanceType,numInstances)
            #print(type(numInstances))
            mem = ''
            loaded_model = ''
            error = 0
            coreCount = 0
            if cloudProvider == "aws":
                filename = ratetype+"/"+cloudProvider+"/"+instanceType+".sav"
                print(filename)
                if(instanceType == 't2.nano'):
                    mem = 0.5
                    coreCount = 1
                    loaded_model = pickle.load(open(filename, 'rb'))
                elif(instanceType == 't2.micro'):
                    mem = 1
                    coreCount = 1
                    loaded_model = pickle.load(open(filename, 'rb'))
                elif (instanceType == 't2.small'):
                    mem = 1
                    coreCount = 1
                    loaded_model = pickle.load(open(filename, 'rb'))
                elif (instanceType == 't2.medium'):
                    mem = 2
                    loaded_model = pickle.load(open(filename, 'rb'))
                elif(instanceType == 't2.large'):
                    mem = 2
                    loaded_model = pickle.load(open(filename, 'rb'))
                else:
                    error =1
                if error == 0:
                    val = [[coreCount,float(numInstances),mem]]
                    predict = loaded_model.predict(val)
                    output=""+str(predict[0])
                    print (predict)
                    print (output)
                    self.wfile.write(output.encode())
                else:
                    self.wfile.write("error instanceType not matched".encode())
            else:
                self.wfile.write("error csp not matched".encode())
        else:
            self.wfile.write("Path not found".encode())
    def do_HEAD(self):
        self._set_headers()

mongo_host = os.environ['MONGODB_HOST']
mongo_port = os.environ['MONGODB_PORT']
mongo_username= os.environ['MONGODB_USER']
mongo_password = os.environ['MONGODB_PASS']

def run(server_class=HTTPServer, handler_class=S, port=9001):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print ('Starting httpd...')
    httpd.serve_forever()

if __name__ == "__main__":
    from sys import argv

    if len(argv) == 2:
        run(port=int(argv[1]))
    else:
        run()