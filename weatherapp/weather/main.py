from flask import Flask, jsonify
from flask_cors import CORS
import requests
import os

app = Flask(__name__)
CORS(app)

@app.route("/")
def health():
    return "The service is running", 200

@app.route('/<city>')
def hello(city):
    url = "https://weatherapi-com.p.rapidapi.com/current.json"
    querystring = {"q":city}
    headers = {
        'x-rapidapi-host': "weatherapi-com.p.rapidapi.com",
        'x-rapidapi-key': os.getenv("APIKEY")
    }
    response = requests.request("GET", url, headers=headers, params=querystring)
    return jsonify(response.text)


if __name__ == '__main__':
    app.run(host="0.0.0.0")
