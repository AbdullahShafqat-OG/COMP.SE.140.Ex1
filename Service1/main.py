from flask import Flask, request
import requests
import shutil
import time
import os
from datetime import datetime, timezone

app = Flask(__name__)
start_time = time.time()

SERVICE2_URL = "http://service2:3000"
STORAGE_URL = "http://storage:8080"
VSTORAGE_PATH = "/app/storage/vstorage"

def get_status_record():
    uptime_hours = (time.time() - start_time) / 3600
    
    total, used, free = shutil.disk_usage('/')
    free_mb = free // (1024 * 1024)
    
    timestamp = datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace('+00:00', 'Z')
    
    return f"{timestamp}: uptime {uptime_hours:.2f} hours, free disk in root: {free_mb} MBytes"

def write_to_vstorage(record):
    os.makedirs(os.path.dirname(VSTORAGE_PATH), exist_ok=True)
    with open(VSTORAGE_PATH, 'a') as f:
        f.write(record + '\n')

def send_to_storage(record):
    try:
        requests.post(f"{STORAGE_URL}/log", data=record, headers={'Content-Type': 'text/plain'})
    except Exception as e:
        print(f"Error sending to storage: {e}")

@app.route('/status')
def status():
    record1 = get_status_record()
    
    send_to_storage(record1)
    write_to_vstorage(record1)
    
    try:
        response = requests.get(f"{SERVICE2_URL}/status")
        record2 = response.text.strip()
        
        combined = f"{record1}\n{record2}"
        return combined, 200, {'Content-Type': 'text/plain'}
    except Exception as e:
        return f"Error contacting Service2: {e}", 500

@app.route('/log')
def log():
    try:
        response = requests.get(f"{STORAGE_URL}/log")
        return response.text, 200, {'Content-Type': 'text/plain'}
    except Exception as e:
        return f"Error contacting Storage: {e}", 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8199)