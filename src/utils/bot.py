import tweepy
import os

consumer_key = os.environ["TWITTER_API_KEY"]
consumer_secret = os.environ["TWITTER_API_SECRET"]
access_token = os.environ["TWITTER_ACCESS_TOKEN"]
access_token_secret = os.environ["TWITTER_ACCESS_TOKEN_SECRET"]

auth = tweepy.OAuthHandler(consumer_key, consumer_secret)
auth.set_access_token(access_token, access_token_secret)

api = tweepy.API(auth)
geolocation = "38.215311,-85.758189,300mi"
lang = "en"

tweet_candidates = api.search(
    q="filter:safe \"can't find parking\" OR \"trying to find parking\" -filter:retweets -filter:links -filter:images",
    geolocation=geolocation,
    result_type="recent",
    count=100,
    language="eng",
    tweet_mode="extended")
for tweet in tweet_candidates: 
    print(tweet.full_text)
