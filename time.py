import subprocess
import time
start = time.time()
res = subprocess.call(['./GoCurl','-b','foo=6; bar=28; baz=496','-I',"http://example.com/"])
elapsed_time = time.time() - start
print("elapsed_time:{0}".format(elapsed_time) + "[sec]\n")
time.sleep(0.2)
start = time.time()
res = subprocess.call(['curl','-b','foo=6; bar=28; baz=496','-I',"http://example.com/"])
elapsed_time = time.time() - start
print("elapsed_time:{0}".format(elapsed_time) + "[sec]\n")
time.sleep(0.2)
start = time.time()
res = subprocess.call(['./GoCurl','-b','foo=6; bar=28; baz=496','-H',"Accept-Encoding:gzip",'-I',"http://example.com/"])
elapsed_time = time.time() - start
print("elapsed_time:{0}".format(elapsed_time) + "[sec]\n")
time.sleep(0.2)
start = time.time()
res = subprocess.call(['curl','-b','foo=6; bar=28; baz=496','-H',"Accept-Encoding:gzip",'-I',"http://example.com/"])
elapsed_time = time.time() - start
print("elapsed_time:{0}".format(elapsed_time) + "[sec]\n")
time.sleep(0.2)
