import subprocess
import time
start = time.time()
res = subprocess.call(['./GoCurl','-I',"http://example.com/"])
elapsed_time = time.time() - start
print("elapsed_time:{0}".format(elapsed_time) + "[sec]\n")
time.sleep(0.2)
start = time.time()
res = subprocess.call(['curl','-I',"http://example.com/"])
elapsed_time = time.time() - start
print("elapsed_time:{0}".format(elapsed_time) + "[sec]\n")
time.sleep(0.2)
start = time.time()
res = subprocess.call(['./GoCurl','-H',"Accept-Encoding:gzip",'-I',"http://example.com/"])
elapsed_time = time.time() - start
print("elapsed_time:{0}".format(elapsed_time) + "[sec]\n")
time.sleep(0.2)
start = time.time()
res = subprocess.call(['curl','-H',"Accept-Encoding:gzip",'-I',"http://example.com/"])
elapsed_time = time.time() - start
print("elapsed_time:{0}".format(elapsed_time) + "[sec]\n")
time.sleep(0.2)
