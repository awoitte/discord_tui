# discord_tui

A very simple, one channel at a time, text-based user interface for Discord DMs.

### Usage:
```
discord_tui -c <path to config> -u <user to chat with>
```

### Config:
```JSON
{
        "username": "<username>",
        "password": "<password>",
        "<user to chat with>": "<user's id>"
}
```

The key you provide for the user in your config is the name you give with the -u flag. It does not have to match the users actual Discord name.

A simple way to get a user's ID is to enable dev mode in the Discord client and right click on a user's avatar.
