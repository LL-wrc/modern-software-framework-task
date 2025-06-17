from flask import Flask, request, jsonify, send_from_directory
import requests
import os

app = Flask(__name__, static_folder='../client')

# Go backend service URL
GO_BACKEND_URL = "http://localhost:8000/calculate"   

@app.route('/')
def serve_index():
    return send_from_directory(app.static_folder, 'index.html')

@app.route('/<path:path>')
def serve_static_files(path):
    return send_from_directory(app.static_folder, path)

@app.route('/api/calculate', methods=['POST'])
def handle_calculate():
    data = request.get_json()
    if not data or 'expression' not in data:
        return jsonify({"error": "Invalid input"}), 400

    expression = data['expression']

    try:

        response = requests.post(GO_BACKEND_URL, json={"expression": expression}, timeout=5)
        response.raise_for_status() # Raise an exception for HTTP errors (4xx or 5xx)
        backend_response = response.json()
        return jsonify(backend_response)

    except requests.exceptions.RequestException as e:
        print(f"Error connecting to Go backend: {e}")
        return jsonify({"error": "Calculation service unavailable"}), 503
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        return jsonify({"error": "An internal error occurred"}), 500

if __name__ == '__main__':

    client_dir = os.path.join(os.path.dirname(__file__), '..', 'client')
    if not os.path.exists(client_dir):
        print(f"Error: Client directory not found at {client_dir}")
        print("Please ensure the client files (index.html, style.css, script.js) are in the correct location.")
    else:
        print(f"Serving client files from: {os.path.abspath(client_dir)}")
        app.run(host='0.0.0.0', port=8080, debug=True)