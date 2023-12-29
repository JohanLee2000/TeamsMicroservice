# TeamsMicroservice

Developed to send messages to a Microsoft Teams channel. Requires an Incoming Webhook configuration, specifically the URL
to send messages to that channel

Usage: go run main.go <inputfile.json>

JSON file format:
{
    "WebhookURL": "<WebhookURL>",
	"Title": "<Title of Message>",
	"Text": "<Text of Message>"
}

How to use:
1) Click on the 3 dots next to the desired channel 
2) Click Connectors
3) Click Add/configure for Incoming Webhook
4) Enter an appropriate name and click Create
5) Copy the URL
6) Paste the URL in the <WebhookURL> field in json file 
7) Configure your message title/text along with any other fields in json file
8) Run 'go run main.go <inputfile.json>'


If you set up Webhook before but forgot URL:
Do steps 1-2
3) On the left side (under Manage), click on Configured
4) Under Incoming Webhook, click on # Configured
5) Click on Manage (of the name of your webhook) <- You can also edit the name of your webhook this way and set an image
6) The URL should be at the bottom to Copy


For more documentation:
- Message Cards
    https://learn.microsoft.com/en-us/outlook/actionable-messages/message-card-reference
- Adaptive Cards (actionable messages only via email) for future implementation?
    https://learn.microsoft.com/en-us/outlook/actionable-messages/adaptive-card

References
Microsoft webhook vs connectors exercise - https://learn.microsoft.com/en-us/training/modules/msteams-webhooks-connectors/
Microsoft webhook vs connectors video - https://www.youtube.com/watch?v=EtMOuBi82LI
Send a message (messageCard/AdaptiveCard) - https://github.com/atc0005/go-teams-notify
Send a message/notification - https://github.com/atc0005/send2teams
https://github.com/dasrick/go-teams-notify
