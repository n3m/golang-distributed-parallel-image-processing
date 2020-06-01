# Architecture Document

This server arquitecture consists in a Single Mother Server that handles jobs to it's workers.
Petitions are made via REST API endpoints that require previous authentication.

# Work Balance

This server is capable of balancing the loads of the workers via scheduling jobs to the lower % of usage worker.
This architecture is capable of handling as many workers as it is required in order to be scalable and efficient. The Controller will make it possible.

![Server Architecture](architecture.png)

# Functionality:

- Endpoint: [/login] Accessing this endpoint with the correct credentials (admin:password) will give you an access token.
- Endpoint: [/status] Acessing this endpoint will return information about the system status with all the current online workers information
- Endpoint: [/status/<workername>] Acessing this endpoint will return all the information about a single worker.
- Endpoint: [/workloads/test] Acessing this endpoint will trigger an end to end test to any available worker
- Endpoint: [/upload] Acessing this endpoint will give you the ability to Upload an image to the system
- Endpoint: [/logout] Acessing this endpoint will revoke the current token.

### For a user guide, please visit the file <user-guide.md> in this repository
