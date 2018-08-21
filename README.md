# database
In this project, i create a database that allow users to save data into a leveldb or a file.
I use mmap to write and read data in the file easily. In next step, i want to apply this database in blockchain.
Recent block states should be saved into a cache (leveldb) before persisting them into a file.