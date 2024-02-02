package commands

const (
	helpMsg = `available commands:
	- 'add [type]' - to add record with type = [text, creds, card, bin]
	- 'update' - to update record
	- 'delete' - to delete record
	- 'get [type]' - to show stored records with type = [any, text, creds, card, bin
	- 'extract [type]' - to decode and save binary file from local storage with type [bin]

	- 'server [type]' - for manipulation with server with type = [login, register, push, pull]

	- 'exit' - to close application`
)

func (c *Commands) Help() {
	c.iactr.Printf(helpMsg)
}
