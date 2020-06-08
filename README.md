# Factom Admin Chain Reporter

A small utility that queries a Factomd endpoint and then reports most Admin Block events to a discord webhook. Events not included are: MinuteNumber (deprecated), and DBSignature. For the coinbase every 25 blocks, it will print out the FCT sum but not the entire transaction. 

To configure, copy `config.ini.EXAMPLE` and rename it to `config.ini`. Place it in the current working directory when running.

* `factomd`: This should point to a factomd API endpoint. Don't include trailing /v2/. 
* `webhook`: The url to a discord webhook
* `name`: The display name the bot should use. If unspecified, the default webhook setting is used
* `avatar`: The url to an avatar the bot should use. If unspecified, the default webhook setting is used



