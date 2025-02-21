# Connection Request Server

### Spec
“Connection Request Server” --- a system to track the list of currently connected users, and enforce that a single user can only be connected from one computer at a time.

We operate a VPN service. Users install our application, configure it to use their VPN certificate, then connect to one of our many VPN servers. Users are free to install our app on as many computers as they like. However, they should only connect from one computer at a time. Currently, we have no way to enforce this policy. This document describes requirements for a new application called “Connection Request Server” to run on a new server dedicated to keeping track of which users are currently connected to our service, and deciding whether a given user should be allowed to connect now. The server exposes the following functions, to be used by our client app:

1. request_permission_to_connect
2. disconnect
3. heart_beat

Before connecting to the VPN, our client app will call “request_permission_to_connect”. The function will “allow” the request if and only if it believes that this user isn't currently connected from another computer. Otherwise, it will “reject” the request.

When clients disconnect, they will call the “disconnect” function. However, in some cases they may not be able to do so, such as if they suddenly lost network connectivity. Therefore, we also implement a heart-beat system. Clients will call this function repeatedly at a predetermined time interval while connected. If the server hasn't received a heart-beat from a given user in the expected time period, it will consider that user disconnected.

User Scenarios

For each of the following cases: Assume that the user installed our client app on computers A and B, and activated each with the same code X.

Trying to connect to the VPN while already connected on another computer

1. Connect on computer A → ALLOW
2. A few seconds later, while A is still connected, try to connect on computer B → REJECT
3. Some time later, click our app's “disconnect” button on computer A
4. 1 second later, try to connect on computer B. → ALLOW 
5. 1 second later, try to connect on computer A. → REJECT

Lost network connectivity or crash, then connect from another computer

1. Connect on computer A → ALLOW
2. Computer A: Unplug the network cable (to disconnect the app without giving it a chance to send us a disconnect message)
3. 1 second later: Computer B: Try to connect. In theory this should be allowed, but because we didn't get the disconnect message, we still think computer A is connected and will REJECT this request.
4. 5 minutes later: Computer B: try to connect: ALLOW. Because we didn't receive a heartbeat in more than x minutes, we consider computer A disconnected and therefore allow this request.

Lost network connectivity, then reconnect from the same computer

1. Connect on computer A → ALLOW
2. Computer A: unplug the network cable, wait for the app to notice that it's disconnected
3. Try to connect on computer B (less than 5 minutes after step #2) → REJECT because we believe A is still connected.
4. Plug the network cable back into A
5. Try to connect on computer A (still less than 5 minutes after step #2). → ALLOW. We still believe that A is connected, but because this request comes from A itself (same device_id), we allow it. This is why each client needs to identify itself with a device_id so we can uniquely identify each computer.

Connection requests at exactly the same second

1. Hit “connect” on both A and B within exactly the same second. → We'll allow both connections. However, if one of the requests came more than a second after the first request, we would have rejected it.

Connection-request API

`/request_permission_to_connect`

POST params:

1. activation_code   (this uniquely identifies a user account). Required.
2. device_id  (this uniquely identifies a computer). Required.
3. client_version  (for logging purposes only). Optional.
4. os_version  (for logging purposes only). Optional.

Return values:

```
<connection_request_response>
	<code>1</code>  // 1 means connection allowed. Code >= 400 means rejected.
	<message>Message string that will be shown to the user if code >= 400</message>
</connection_request_response>
```

Returns code==1 if and only if the given activation_code is not currently connected with another device_id. This requires keeping a data-store of “currently_connected_users” of (activation_code, device_id) pairs.

Possible values for `<code>` and `<message>`:

1. 1  Approved
2. 400 Sorry, your account is currently connected from another computer. You can use our service from multiple computers, but each account can only be connected to our network from one computer at a time. To connect from this computer now, please buy an additional account.
3. 401 Missing parameters. Sorry, we've made a note to fix this. Please try again and contact support if you continue to see this error.
4. 500 Sorry, unknown error. Please try again and contact support if you continue to see this error.

Disconnect API

`/disconnect`

POST params:  (same purpose and requirements as for /request_permission_to_connect)

1. activation_code
2. device_id
3. client_version
4. os_version

Marks the given activation_code as disconnected, regardless of device_id.

Return values: simply the string “ok” always, but this doesn't matter. The client doesn't check the return value.

Heartbeat API

`/heartbeat`

POST params:  (same purpose and requirements as for /request_permission_to_connect)

1. activation_code
2. device_id
3. client_version
4. os_version

Return values: simply the string “ok” always, but this doesn't matter. The client doesn't check the return value.

Performance Requirements

Sampling randomly from a pool of 50,000 activation_codes and 40,000 device_ids:

1. “request_permission_to_connect” and “disconnect” functions: maximum processing time per request with a load of 10 requests per second: 50ms.
2. “heart_beat”: maximum processing time per request with a load of 150 requests per second: 10ms.  (this is a lot of load … server-side needs to be specifically designed for high loads)

Requirements

1. 100% automated test coverage. Measure coverage and prove that it's 100%. We prefer rspec tests, though you're welcome to use another tool if you like, just let us know.
2. Automated tests proving that the app meets its performance requirements. For example, a tool used to generate requests against a staging environment, we run it on a few clients to generate 160 requests per second in total, then measure the server's response times.
3. Host on heroku. We'll use separate apps for test, staging, and production environments. The developer will have access to test and staging. No access to production.
4. Send an email to admins in case of an exception
5. Logging: write to a database: log the full params and return values of each call to /request_permission_to_connect and /disconnect. Do NOT log heart-beats.
6. Log data retention: automatically delete log data older than 14 days
7. No reporting requirements. This system will have no UI. The admin will manually query the DB when needed. 
8. Environment-specific configuration values:
9. heart_beat_period_minutes [int], heart_beat_grace_period_seconds [int]. If the server hasn't received a heart-beat from a given activation_code in heart_beat_period_minutes + heart_beat_grace_period_seconds, consider that account disconnected.
10. List of admin email addresses for exception notifications
11. Data-persistence: if the server crashes, it's OK to lose the state table of currently connected users and therefore consider all users disconnected. In other words, we do NOT need to keep the state table in persistent store. In memory is good enough. However, log data DOES need to persist.
12. All source-code in github, daily commits during work-days, heroku git configured as a 'remote'.

Questions

Which technology stack do you recommend?

1. Node.js + redis, rack + redis, rails3 + redis, or anything else? How many dynos do you estimate we'll need to meet our performance requirements? (it's ok if you don't know, we're just wondering if you have experience with apps of such high loads)
2. Which database?
3. Where do you store the state table of connected users? (for performance reasons, we suggest keeping this in memory only) How do you make sure that it's shared across multiple dynos?