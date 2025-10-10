# SPDX-FileCopyrightText: 2021-2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT

curl -X POST https://id.twitch.tv/oauth2/authorize?client_id=${TWITCH_CLIENT_ID}&redirect_uri=https://www.bitsmasher.net/twitch/callback&response_type=code&scope=user%3Aread%3Aemail+channel%3Aread%3Asubscriptions&state=somerandomstringXhere

exit 0

curl -X POST https://id.twitch.tv/oauth2/token \
     -H 'Content-Type: application/x-www-form-urlencoded' \
     -d "client_id=${TWITCH_CLIENT_ID}" \
     -d "client_secret=TWICH_CLIENT_SECRET" \
     -d 'code=YOUR_AUTHORIZATION_CODE_FROM_STEP_2' \
     -d 'grant_type=authorization_code' \
     -d 'redirect_uri=YOUR_REGISTERED_REDIRECT_URI'
