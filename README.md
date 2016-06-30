Travel Audience
===============

Bootstrap Flow:
- Main file loads configuration using Service locator
- Parsing configuration then passing it to the dispatcher
- Dispatcher initializes workers pool, jobs queue, and numbers server
- Numbers server initializes listens to a port, and initializes Aggregator.
- Aggregator creates various channels to receive data/flags on

Request flow:
- The server will handle the request through the path "/numbers"
- The server will parse the query string, and create a number of jobs according to the passed urls.
- The server will block on an output from *AggregationQueue* channel.
- The workers will start processing jobs, fetch the data, and pass it to the response payload processor.
- The response payload processor will parse the string response and pass it back to the worker as normal array.
- The worker will push the data to "MergeQueue" channel.
- The Aggregator is monitoring *MergeQueue* channel, and appends newly received data to its own data storage.
- The Aggregator will time out after pre-configured period of time, and will disregard all items are coming to "MergeQueue" chanel.
- The Aggregator will raise a flag that the merging is done, and pass the data to the Operator
- The Operator will remove duplicates, and sort the data then serve it back to the Aggregator.
- The Aggregator will push the data to *AggregationQueue*.
- The server will respond to the client.