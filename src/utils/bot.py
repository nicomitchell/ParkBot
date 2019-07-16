import tweepy
import os
from random import randint
from time import sleep
consumer_key = os.environ["TWITTER_API_KEY"]
consumer_secret = os.environ["TWITTER_API_SECRET"]
access_token = os.environ["TWITTER_ACCESS_TOKEN"]
access_token_secret = os.environ["TWITTER_ACCESS_TOKEN_SECRET"]

auth = tweepy.OAuthHandler(consumer_key, consumer_secret)
auth.set_access_token(access_token, access_token_secret)

api = tweepy.API(auth)

while(True):
    geolocation = "38.215311,-85.758189,300mi"
    tweet_candidates = api.search(
        q="filter:safe \"can't find parking\" OR \"trying to find parking\" -filter:retweets -filter:links -filter:images",
        geolocation=geolocation,
        result_type="recent",
        count=100,
        language="eng",
        tweet_mode="extended")

    s = api.create_favorite(tweet_candidates[randint(0,len(tweet_candidates)-1)].id)
    if str(input("continue? ")).lower()=="n":
        break
    sleep(randint(3800,10000))