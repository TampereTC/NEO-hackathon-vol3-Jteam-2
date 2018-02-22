# -*- coding: utf-8 -*-
"""
add threads
add partitioning 
etc...
"""
import time 

from kafka import KafkaConsumer, KafkaProducer

totalMessages = 10000
producer = KafkaProducer(bootstrap_servers='kafka:9092')


startTime = time.time()

for x in range(0, totalMessages):
    msgText = str(int((x/10000)%10)) +  str(int((x/1000)%10)) +  str(int((x/100)%10)) +  str(int((x/10)%10)) + str(int(x%10))
    
    future = producer.send('kkTopicPython', bytes(msgText, encoding='utf-8'))
    result = future.get()
    if ( int(x%1000) == 0): 
        print ("message ", x , "result: ", result)
    producer.flush()
    
elapsedTime = time.time() - startTime
msgPerSec = totalMessages / elapsedTime
print ("Time elapsed ", elapsedTime, " sec")
print ( msgPerSec, " msg per second.")
producer.close()