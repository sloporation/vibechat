# VibeChat

VibeChat is an attempt to provide a self hosted alternative to Discord that is 
flexible enough to be implemented for more.

Our goal is primarily to define the spec, with our client, server and bouncer 
being an example implementation of it.

## Key Concepts

It's best to think about VibeChat as being a Postfix rip off:

1. Sender deposits message for delivery.
2. Server spools it until recipient is online.
3. Recipient accepts delivery of message.
4. Message is purged from server.

There are two key deviations from Postfix in our workflow:

1. The server hands off directly to the client, not to a mailbox.
2. The spec requires messages to be sent in json, not limited to messages.

The final piece of this puzzle is our bouncer. Because we dont have a mailbox 
server and we're delivering direct to client, users need a service to act on 
their behalf and retain message history long term.

- The bouncer acts on behalf of a client, acting like a proxy for messages
- Signalling still comes from the client, just routes via the proxy
- Clients are aware that they're connecting to a proxy, and behave appropriately
- Puts the onus of message retention on the user
- Allows the user to keep multiple devices in sync

We'll go into more detail later, but the concept in my mind is that:

- The server acts purely as a relay of JSON payloads.
- The server can optionally maintain history, at the providers discretion.
- The server can define multiple "spaces". 
- The server defines the structure of each space, similar to Discord.
    - IE static channel list or role-based display
    - server defines channel layout
    - etc.

### Rules for Development

This is kept in key concepts so that you understand that this is, infact, a key 
concept.

For servers:

- Binaries must be simple to run.
- Binaries must be easily configuration through config.yml or env vars.
- Env vars override yaml.
- The administrator has full control, but can't always read messages.

For bouncers:

- Must be dead simple to run.
- Must be dead simple to maintain.
- Must reliably store messages.
- Must be easy to prune old data upon user request.
- Must allow for one bouncer to support multiple clients/servers.

For clients:

- User should be able to operate entirely within the client.
- Must allow multiple users to be created, and multiple servers to be added.
- Must be aware of bouncers and servers, and act appropriately.
- User must be aware of what's happening at all times. No surprises.

## Ideas for Implementation

Recording this so that whoever reads this or chooses to contribute understands 
the direction that we're looking to go.

- One-to-one messaging
    - One of two primary focuses of the app.
    - Basic user@example.com sends messages to user@elpmaxe.moc.
- One-to-many messaging
    - The second primary focus of the app.
    - Basic user@example.com sends message to room #someroom@space.server.com.
- Voice
    - To operate similar to Discords.
    - One-to-one calls alert.
    - Group chats are like rooms you walk into, and message people to join you.
- Bots
    - Specifically registered to a user on a server.
    - Acts much like a typical client.
- Shopping
    - A server can list items for sale; physical or digital.
    - Users can make typical purchases and track via the app.
    - Client basically acts as the store front, instead of ecom. Server as API.
- Marketplace
    - Servers or spaces can define marketplaces for users to make listings.
    - All available listings could be made available from a single view.
    - Sketchy but could work like Facebook Marketplace.
- Donations
    - Signalling defined to allow users to donate towards server running cost.
    - Server should provide some kind of signal for bots to provision benefits.
- Events
    - Servers/spaces could accept organizing of calendar based events.
    - Calendars can be synced with caldav.

The idea is to create an environment that cultivates a community similarly to 
what we did back in the days of forums.

Provides the ability of people to connect, then provide frameworks that allow 
them to interact freely.

## Lawful Intercept

This is becoming a bit of a hot point so I'm going to bring this up early.

The *only* way to protect against lawful intercepts is to exchange encryption 
keys offline. You would have to give your public key to someone outside of the 
client.

In all message proxying specs, you can encrypt messages by signaling a key 
exchange and storing them on device. When clients lose their key, you can 
signal for a refresh.

When this is the case, the server can intercept the unencrypted key exchange 
and store them for later decryption. This means messages will appear to be end 
to end encrypted, but the server could be caching messages and decrypting out 
of band.

This is a security issue as well as a compliance requirement, and we have to 
be open that this is a possible scenario with our spec.

This is good because as a server provider, you can have an option for complying 
with lawful intercept requests and avoid fines or jail. In another sense, it's 
a negative because the same mechanics can be used by service providers to spy 
on users.

One huge benefit is that as an open project, we can only really provide the 
code but we can't really dictate how people use it. The onus is on the server 
providers, which means you can avoid being wrapped up in compromises or widenet 
data harvesting by running your own server.

And we aim to make that as simple as possible.
