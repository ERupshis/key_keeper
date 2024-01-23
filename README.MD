# Password manager GophKeeper
## General requirements
GophKeeper is a client-server system that allows users to safely and securely store logins, passwords, binary data and other private information.
image
The server should implement the following business logic:
- user registration, authentication and authorization;
- storage of private data;
- data synchronization between several authorized clients of the same owner;
- transfer of private data to the owner upon request.

The client shall implement the following business logic:
- authentication and authorization of users on the remote server;
- access to private data on request.

Functions whose implementation is left to the discretion of the implementer:
- server-side or client-side data creation, editing, and deletion;
- new user registration format;
- selection of storage and data storage format;
- ensuring the security of data transmission and storage;
- client-server interaction protocol;
- mechanisms for user authentication and authorization of access to information.

Additional requirements:
- the client must be distributed as a CLI application with the ability to run on Windows, Linux, and Mac OS platforms;
- the client must allow the user to get information about the version and build date of the client binary file.

Types of stored information
- login/password pairs;
- arbitrary text data;
- arbitrary binary data;
- bank card data.

For any data, it should be possible to store arbitrary textual metainformation (data belonging to a website, person or bank, lists of one-time activation codes, and so on).

## Abstract scheme of interaction with the system
The basic scenarios of user interaction with the system are described below. They are not exhaustive - solving individual scenarios (e.g., resolving data conflicts on the server) is left to the discretion of the implementer.

For a new user:
- The user receives a client for the platform he/she needs.
- The user goes through the initial registration procedure.
- The user adds new data to the client.
- The client synchronizes the data with the server.

For an existing user:
- The user receives the client under the platform he/she needs.
- The user goes through the authentication procedure.
- The client synchronizes the data with the server.
- The user requests data.
- The client displays the data for the user.