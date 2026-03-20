# VibeChat Spec

## Constants

- All messages are sent via the payload spec
    - Payload spec defines the action, sender, recipient and payload
- Payload changes depending on the action sent
- All data is base64 encoded for serialization

## Payload Spec

All payloads are identical at the top level. It's to provide a consistent 
recipient/sender combo:


| Key | Purpose                                                               |
|-----|-----------------------------------------------------------------------|
| a   | The action being performed                                            |
| s   | The sender of the payload                                             |
| r   | The recipient of the payload                                          |
| p   | The payload itself                                                    |


For example, a message payload for a key exchange going from user-to-user:

```
{
    a: 0,
    s: "alice@funclub.com",
    r: "bob@unfunclub.com",
    p: {
        ...
    }
}
```

To send the message itself, you'd just change the action and provide the 
appropriate payload.

## Key Exchange

When a user wants to message another user using encryption, you first need to 
do a key exchange.

While this could theoretically be done offline, our first payload blob is 
dedicated to exchanging keys when two users connect.


| Key | Purpose                                                               |
|-----|-----------------------------------------------------------------------|
| p   | The public key for that particular user                               |


There isn't much to it. The public key is just provided. The client stores it. 
The key is removed from the server.

```
{
    a: 0,
    s: "alice@funclub.com",
    r: "bob@unfunclub.com",
    p: {
        p: "c2VyaWFsaXplZCBwdWIga2V5Cg==",
    }
}
```

## Direct Message

A user messaging a user is also relatively simple.


| Key | Purpose                                                               |
|-----|-----------------------------------------------------------------------|
| m   | The message being sent to the recipient                               |
| i   | List of images associated with this message                           |
| r   | The message ID of the message being replied to                        |
| t   | The message ID that this is being sent as a thread to                 |


```
{
    a: 1,
    s: "alice@funclub.com",
    r: "bob@unfunclub.com",
    p: {
        m: "c29tZSBtZXNzYWdlCg==",
        i: ["iVBOR..."],
        r: "ZmRhODUyYjAtZTk5Yy00MjhmLWE4NDMtMTUxNWI3YzlmODAzCg=="
    }
}
```


