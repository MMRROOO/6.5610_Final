To run our current code navigate to the PIR directory and run:

go build *.go   // this builds the executable
./main          // run the executable

This will create 2 peer instances and a tracker instance and start the network localy.

Future work will involve implementing the more advanced algorithms inthe paper and implementing functionality across different computers. This will involve more research into current networking systems.

After this we will need to implement more advanced algorithms from bittorent in order to make the system more efficient with large number of peers.

The current code is not secure and the LWE parameters are chosen for the sake of creating a proof of concept. There are also other shortcuts that were taken for the sake of minimizing the complexity of our program while mainting the main structure we mentioned in our paper that degrade security.

