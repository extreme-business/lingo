# Lingo
Lingo is a chat application that allows users to chat with each other in real time. 

Having access to an email address also means you can have access to messages sent to you by other users. 
You will be granted a JWT by login in with a token send to your email. 
This JWT will be used to authenticate you when you want to send or receive messages.

An address of a user is the sha256 hash of the email address of the user. Email address hashes are used to protect the privacy of users, although 100% privacy is not guaranteed this way and although stored as a hash, a plain text email address is still needed to send a login token to.

protip: use a + in your email address to create a unique email address: its-a-me+askdasd@gmail.com

## domain
```mermaid
erDiagram
    MESSAGE {
        UUID id
        string sender_user_id
        string recipient_user_id
        string content
    }
    USER {
        string email_hash
    }
    ACCOUNT ||--o{ MESSAGE : sent
    ACCOUNT ||--o{ MESSAGE : received
```

Having an account model gives the following benefits:
- users can delete their account
- users can only send messages to other users that have an account

>>>>>>> dd21906 (init)
## login
```mermaid
sequenceDiagram
    User->>Lingo: request login token
    Lingo->>User: send login token to email
    User->>Lingo: login with token token from email
    Lingo->>User: send jwt token
```

## send message
```mermaid
sequenceDiagram
    User->>Lingo: send message
    Recipient->>Lingo: request messages
```