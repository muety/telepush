#!/usr/bin/env python3
# Author: https://github.com/peet1993

import requests
import os
import argparse

hook_url = '<URL_HERE>'
hook_recipient_ids = {
    'main': '<TOKEN_HERE>',
}

def send_notification(name, text, recipient):
    json_data = {
        'recipient_token': hook_recipient_ids[recipient],
        'origin': name,
        'text': text
    }
    resp = requests.post(hook_url, json=json_data)
    return resp.status_code, resp.text

def main():
    parser = argparse.ArgumentParser(description="Send a notification via Webhook2Telegram.")
    parser.add_argument("name", help="Name of the Bot shown in the chat.")
    parser.add_argument("message", help="Notification message to send.")
    parser.add_argument("recipients", metavar="recipients", nargs="*", default="main", choices=list(hook_recipient_ids.keys()), help="Name(s) of the recipient(s) of the notification. Possible recipients are " + ", ".join(list(hook_recipient_ids.keys())) + ". Default recipient is 'main'.")
    args = parser.parse_args()

    # Fix problem with single default in multi-choice
    if type(args.recipients) != list:
        args.recipients = [ args.recipients ]

    for current in args.recipients:
        status_code, body = send_notification(args.name, args.message, current)
        if status_code != 200:
            print("Error sending data to telegram bot: " + str(status_code) + " - " + body)

if __name__ == "__main__":
    main()