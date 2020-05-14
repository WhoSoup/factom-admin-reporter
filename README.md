# Factom Admin Chain Reporter

A small utility that queries a Factomd endpoint and then reports most Admin Block events to a discord webhook. Events not included are: MinuteNumber (deprecated), and DBSignature. For the coinbase every 25 blocks, it will print out the FCT sum but not the entire transaction. 

To configure, copy `config.ini.EXAMPLE` and rename it to `config.ini`. Place it in the current working directory when running.

* `factomd`: This should point to a factomd API endpoint. Don't include trailing /v2/. 
* `webhook`: The url to a discord webhook
* `name`: The display name the bot should use (otherwise it's left up to the webhook setting)
* `avatar`: The url to an avatar the bot should use (otherwise it's left up to the webhook setting)



