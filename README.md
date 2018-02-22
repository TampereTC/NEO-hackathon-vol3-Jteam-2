# NEO-hackathon-vol3-Jteam-2
Some rudimentary performance tests with Shopify.sarama

Based on sarama example at: https://github.com/tcnksm-sample/sarama/blob/master/sync-producer/main.go

Kafka installed to kubernetes as in hackatlon instructions, without modifications.

Windows port forwarding setup.

Run producer. 

SyncProducer sends 5-byte varying ("00001","00002"...) messages in loop, modify code to try different setups. 

Not finished:
Number of topics should be configurable, now it is hardcoded to 10 with manual round-robin (1,2,3,4,5,6,7,8,9,0,1,2,3...) sending.
Number of bytes is hardcoded to 5, should be configurable. 

"someResults.xls" excel has some test results written down. Everything (kafka inside kubernetes, provides as standalone program in windows host) running in single laptop. 
