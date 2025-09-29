const express = require('express');
const axios = require('axios');
const fs = require('fs');
const path = require('path');

const app = express();
const startTime = Date.now();

const STORAGE_URL = 'http://storage:8080';
const VSTORAGE_PATH = '/app/storage/vstorage';

function getDiskUsage() {
    try {
        const { execSync } = require('child_process');
        const output = execSync('df / | tail -1', { encoding: 'utf8' });
        const parts = output.split(/\s+/);
        const availableKB = parseInt(parts[3]);
        return Math.floor(availableKB / 1024);
    } catch (error) {
        console.error('Error getting disk usage:', error);
        return 0;
    }
}

function getStatusRecord() {
    const uptimeHours = (Date.now() - startTime) / (1000 * 60 * 60);
    const freeMB = getDiskUsage();
    const timestamp = new Date().toISOString().split('.')[0] + 'Z';
    
    return `${timestamp}: uptime ${uptimeHours.toFixed(2)} hours, free disk in root: ${freeMB} Mbytes`;
}

function writeToVStorage(record) {
    try {
        const dir = path.dirname(VSTORAGE_PATH);
        if (!fs.existsSync(dir)) {
            fs.mkdirSync(dir, { recursive: true });
        }
        fs.appendFileSync(VSTORAGE_PATH, record + '\n');
    } catch (error) {
        console.error('Error writing to vstorage:', error);
    }
}

async function sendToStorage(record) {
    try {
        await axios.post(`${STORAGE_URL}/log`, record, {
            headers: { 'Content-Type': 'text/plain' }
        });
    } catch (error) {
        console.error('Error sending to storage:', error);
    }
}

app.get('/status', async (req, res) => {
    const record = getStatusRecord();
    
    await sendToStorage(record);
    writeToVStorage(record);
    
    res.set('Content-Type', 'text/plain');
    res.send(record);
});

const PORT = 3000;
app.listen(PORT, '0.0.0.0', () => {
    console.log(`Service2 running on port ${PORT}`);
});